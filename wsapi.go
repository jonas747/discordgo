package discordgo

import (
	"encoding/json"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type wsWriter struct {
	session *Session

	conn           net.Conn
	closer         chan interface{}
	incoming       chan interface{}
	sendCloseQueue chan []byte

	writer *wsutil.Writer

	sendRatelimiter chan bool
}

func (w *wsWriter) Run() {
	w.sendRatelimiter = make(chan bool)
	go w.runSendRatelimiter()

	w.writer = wsutil.NewWriter(w.conn, ws.StateClientSide, ws.OpText)

	for {
		select {
		case <-w.closer:
			select {
			// Ensure we send the close frame
			case body := <-w.sendCloseQueue:
				w.sendClose(body)
			default:
			}
			return
		case msg := <-w.incoming:
			var err error
			switch t := msg.(type) {
			case []byte:
				err = w.writeRaw(t)
			default:
				err = w.writeJson(t)
			}

			if err != nil {
				w.session.log(LogError, "Error writing to gateway: %s", err.Error())
				return
			}
		case body := <-w.sendCloseQueue:
			w.sendClose(body)
		}
	}
}

func (w *wsWriter) runSendRatelimiter() {
	tc := time.NewTicker(time.Millisecond * 500)
	defer tc.Stop()
	defer func() {
		close(w.sendRatelimiter)
	}()

	for {
		select {
		case <-tc.C:
			select {
			case w.sendRatelimiter <- true:
			case <-w.closer:
				return
			}
		case <-w.closer:
			return
		}
	}
}

func (w *wsWriter) writeJson(data interface{}) error {
	serialized, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return w.writeRaw(serialized)
}

func (w *wsWriter) writeRaw(data []byte) error {
	<-w.sendRatelimiter

	_, err := w.writer.WriteThrough(data)
	if err != nil {
		return err
	}

	return w.writer.Flush()
}

func (w *wsWriter) sendClose(body []byte) error {
	_, err := w.writer.WriteThrough(body)
	return err
}

func (w *wsWriter) Queue(data interface{}) {
	select {
	case <-time.After(time.Second * 10):
	case <-w.closer:
	case w.incoming <- data:
	}
}

func (w *wsWriter) QueueClose(body []byte) {
	select {
	case <-time.After(time.Second * 10):
	case <-w.closer:
	case w.sendCloseQueue <- body:
	}
}

type wsHeartBeater struct {
	sync.Mutex

	writer      *wsWriter
	sequence    *int64
	receivedAck bool
	missedAcks  int
	stop        chan interface{}

	// Called when we received no Ack from last heartbeat
	onNoAck func()

	lastAck  time.Time
	lastSend time.Time
}

func (wh *wsHeartBeater) ReceivedAck() {
	wh.Lock()
	wh.receivedAck = true
	wh.missedAcks = 0
	wh.lastAck = time.Now()
	wh.Unlock()
}

func (wh *wsHeartBeater) UpdateSequence(seq int64) {
	atomic.StoreInt64(wh.sequence, seq)
}

func (wh *wsHeartBeater) Run(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:

			wh.Lock()
			hasReceivedAck := wh.receivedAck
			wh.receivedAck = false
			wh.missedAcks++
			missed := wh.missedAcks

			wh.lastSend = time.Now()
			wh.Unlock()

			if !hasReceivedAck && wh.onNoAck != nil && missed > 4 {
				wh.onNoAck()
			}

			wh.SendBeat()
		case <-wh.stop:
			return
		}
	}
}

func (wh *wsHeartBeater) SendBeat() {
	seq := atomic.LoadInt64(wh.sequence)

	wh.writer.Queue(&outgoingEvent{
		Operation: GatewayOPHeartbeat,
		Data:      seq,
	})
}

func (wh *wsHeartBeater) Times() (send time.Time, ack time.Time) {
	wh.Lock()
	send = wh.lastSend
	ack = wh.lastAck
	wh.Unlock()
	return
}
