package discordgo

import (
	"github.com/jonas747/gojay"
	"github.com/pkg/errors"
)

// This file contains all the possible structs that can be
// handled by AddHandler/EventHandler.
// DO NOT ADD ANYTHING BUT EVENT HANDLER STRUCTS TO THIS FILE.
//go:generate go run tools/cmd/eventhandlers/main.go

// Connect is the data for a Connect event.
// This is a sythetic event and is not dispatched by Discord.
type Connect struct{}

// Disconnect is the data for a Disconnect event.
// This is a sythetic event and is not dispatched by Discord.
type Disconnect struct{}

// RateLimit is the data for a RateLimit event.
// This is a sythetic event and is not dispatched by Discord.
type RateLimit struct {
	*TooManyRequests
	URL string
}

// Event provides a basic initial struct for all websocket events.
type Event struct {
	Operation GatewayOP          `json:"op"`
	Sequence  int64              `json:"s"`
	Type      string             `json:"t"`
	RawData   gojay.EmbeddedJSON `json:"d"`
	// Struct contains one of the other types in this file.
	Struct interface{} `json:"-"`
}

// implement gojay.UnmarshalerJSONObject
func (evt *Event) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "op":
		return dec.Int((*int)(&evt.Operation))
	case "s":
		return dec.Int64(&evt.Sequence)
	case "t":
		return dec.String(&evt.Type)
	case "d":
		if cap(evt.RawData) > 1000000 && len(evt.RawData) < 1000000 {
			evt.RawData = nil
		} else if evt.RawData != nil {
			evt.RawData = evt.RawData[:0]
		}

		return dec.AddEmbeddedJSON(&evt.RawData)
	}

	return nil
}

func (evt *Event) NKeys() int {
	return 0
}

type GuildEvent interface {
	GetGuildID() int64
}

type ChannelEvent interface {
	GetChannelID() int64
}

// A Ready stores all data for the websocket READY event.
type Ready struct {
	Version         int          `json:"v"`
	SessionID       string       `json:"session_id"`
	User            *SelfUser    `json:"user"`
	ReadState       []*ReadState `json:"read_state"`
	PrivateChannels []*Channel   `json:"private_channels"`
	Guilds          []*Guild     `json:"guilds"`

	// Undocumented fields
	Settings          *Settings            `json:"user_settings"`
	UserGuildSettings []*UserGuildSettings `json:"user_guild_settings"`
	Relationships     []*Relationship      `json:"relationships"`
	Presences         []*Presence          `json:"presences"`
	Notes             map[string]string    `json:"notes"`
}

// ChannelCreate is the data for a ChannelCreate event.
type ChannelCreate struct {
	*Channel
}

// ChannelUpdate is the data for a ChannelUpdate event.
type ChannelUpdate struct {
	*Channel
}

// ChannelDelete is the data for a ChannelDelete event.
type ChannelDelete struct {
	*Channel
}

// ChannelPinsUpdate stores data for a ChannelPinsUpdate event.
type ChannelPinsUpdate struct {
	LastPinTimestamp string `json:"last_pin_timestamp"`
	ChannelID        int64  `json:"channel_id,string"`
	GuildID          int64  `json:"guild_id,string,omitempty"`
}

func (cp *ChannelPinsUpdate) GetGuildID() int64 {
	return cp.GuildID
}

func (cp *ChannelPinsUpdate) GetChannelID() int64 {
	return cp.ChannelID
}

// GuildCreate is the data for a GuildCreate event.
type GuildCreate struct {
	*Guild
}

// GuildUpdate is the data for a GuildUpdate event.
type GuildUpdate struct {
	*Guild
}

// GuildDelete is the data for a GuildDelete event.
type GuildDelete struct {
	*Guild
}

// GuildBanAdd is the data for a GuildBanAdd event.
type GuildBanAdd struct {
	User    *User `json:"user"`
	GuildID int64 `json:"guild_id,string"`
}

func (gba *GuildBanAdd) GetGuildID() int64 {
	return gba.GuildID
}

// GuildBanRemove is the data for a GuildBanRemove event.
type GuildBanRemove struct {
	User    *User `json:"user"`
	GuildID int64 `json:"guild_id,string"`
}

func (e *GuildBanRemove) GetGuildID() int64 {
	return e.GuildID
}

// GuildMemberAdd is the data for a GuildMemberAdd event.
type GuildMemberAdd struct {
	*Member
}

// GuildMemberUpdate is the data for a GuildMemberUpdate event.
type GuildMemberUpdate struct {
	*Member
}

// GuildMemberRemove is the data for a GuildMemberRemove event.
type GuildMemberRemove struct {
	*Member
}

// GuildRoleCreate is the data for a GuildRoleCreate event.
type GuildRoleCreate struct {
	*GuildRole
}

// GuildRoleUpdate is the data for a GuildRoleUpdate event.
type GuildRoleUpdate struct {
	*GuildRole
}

// A GuildRoleDelete is the data for a GuildRoleDelete event.
type GuildRoleDelete struct {
	RoleID  int64 `json:"role_id,string"`
	GuildID int64 `json:"guild_id,string"`
}

func (e *GuildRoleDelete) GetGuildID() int64 {
	return e.GuildID
}

// A GuildEmojisUpdate is the data for a guild emoji update event.
type GuildEmojisUpdate struct {
	GuildID int64    `json:"guild_id,string"`
	Emojis  []*Emoji `json:"emojis"`
}

func (e *GuildEmojisUpdate) GetGuildID() int64 {
	return e.GuildID
}

// A GuildMembersChunk is the data for a GuildMembersChunk event.
type GuildMembersChunk struct {
	GuildID int64     `json:"guild_id,string"`
	Members []*Member `json:"members"`
}

func (e *GuildMembersChunk) GetGuildID() int64 {
	return e.GuildID
}

// GuildIntegrationsUpdate is the data for a GuildIntegrationsUpdate event.
type GuildIntegrationsUpdate struct {
	GuildID int64 `json:"guild_id,string"`
}

func (e *GuildIntegrationsUpdate) GetGuildID() int64 {
	return e.GuildID
}

// MessageAck is the data for a MessageAck event.
type MessageAck struct {
	MessageID int64 `json:"message_id,string"`
	ChannelID int64 `json:"channel_id,string"`
}

// MessageCreate is the data for a MessageCreate event.
type MessageCreate struct {
	*Message
}

// MessageUpdate is the data for a MessageUpdate event.
type MessageUpdate struct {
	*Message
}

// MessageDelete is the data for a MessageDelete event.
type MessageDelete struct {
	*Message
}

// MessageReactionAdd is the data for a MessageReactionAdd event.
type MessageReactionAdd struct {
	*MessageReaction
}

// MessageReactionRemove is the data for a MessageReactionRemove event.
type MessageReactionRemove struct {
	*MessageReaction
}

// MessageReactionRemoveAll is the data for a MessageReactionRemoveAll event.
type MessageReactionRemoveAll struct {
	*MessageReaction
}

// PresencesReplace is the data for a PresencesReplace event.
type PresencesReplace []*Presence

// PresenceUpdate is the data for a PresenceUpdate event.
//easyjson:json
type PresenceUpdate struct {
	Presence
	GuildID int64 `json:"guild_id,string"`
}

func (e *PresenceUpdate) GetGuildID() int64 {
	return e.GuildID
}

// implement gojay.UnmarshalerJSONObject
func (p *PresenceUpdate) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "guild_id":
		return errors.Wrap(DecodeSnowflake(&p.GuildID, dec), key)
	default:
		return p.Presence.UnmarshalJSONObject(dec, key)
	}
}

func (p *PresenceUpdate) NKeys() int {
	return 0
}

// Resumed is the data for a Resumed event.
type Resumed struct {
	Trace []string `json:"_trace"`
}

// RelationshipAdd is the data for a RelationshipAdd event.
type RelationshipAdd struct {
	*Relationship
}

// RelationshipRemove is the data for a RelationshipRemove event.
type RelationshipRemove struct {
	*Relationship
}

var _ gojay.UnmarshalerJSONObject = (*TypingStart)(nil)

// TypingStart is the data for a TypingStart event.
type TypingStart struct {
	UserID    int64 `json:"user_id,string"`
	ChannelID int64 `json:"channel_id,string"`
	Timestamp int   `json:"timestamp"`
	GuildID   int64 `json:"guild_id,string,omitempty"`
}

// implement gojay.UnmarshalerJSONObject
func (ts *TypingStart) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "user_id":
		return DecodeSnowflake(&ts.UserID, dec)
	case "channel_id":
		return DecodeSnowflake(&ts.ChannelID, dec)
	case "guild_id":
		return DecodeSnowflake(&ts.GuildID, dec)
	case "timestamp":
		return dec.Int(&ts.Timestamp)
	}

	return nil
}

func (ts *TypingStart) NKeys() int {
	return 0
}

func (e *TypingStart) GetGuildID() int64 {
	return e.GuildID
}

func (e *TypingStart) GetChannelID() int64 {
	return e.ChannelID
}

// UserUpdate is the data for a UserUpdate event.
type UserUpdate struct {
	*User
}

// UserSettingsUpdate is the data for a UserSettingsUpdate event.
type UserSettingsUpdate map[string]interface{}

// UserGuildSettingsUpdate is the data for a UserGuildSettingsUpdate event.
type UserGuildSettingsUpdate struct {
	*UserGuildSettings
}

// UserNoteUpdate is the data for a UserNoteUpdate event.
type UserNoteUpdate struct {
	ID   int64  `json:"id,string"`
	Note string `json:"note"`
}

// VoiceServerUpdate is the data for a VoiceServerUpdate event.
type VoiceServerUpdate struct {
	Token    string `json:"token"`
	GuildID  int64  `json:"guild_id,string"`
	Endpoint string `json:"endpoint"`
}

func (e *VoiceServerUpdate) GetGuildID() int64 {
	return e.GuildID
}

// VoiceStateUpdate is the data for a VoiceStateUpdate event.
type VoiceStateUpdate struct {
	*VoiceState
}

// MessageDeleteBulk is the data for a MessageDeleteBulk event
type MessageDeleteBulk struct {
	Messages  IDSlice `json:"ids,string"`
	ChannelID int64   `json:"channel_id,string"`
	GuildID   int64   `json:"guild_id,string"`
}

func (e *MessageDeleteBulk) GetGuildID() int64 {
	return e.GuildID
}

func (e *MessageDeleteBulk) GetChannelID() int64 {
	return e.ChannelID
}

// WebhooksUpdate is the data for a WebhooksUpdate event
type WebhooksUpdate struct {
	GuildID   int64 `json:"guild_id,string"`
	ChannelID int64 `json:"channel_id,string"`
}

func (e *WebhooksUpdate) GetGuildID() int64 {
	return e.GuildID
}

func (e *WebhooksUpdate) GetChannelID() int64 {
	return e.ChannelID
}

// InviteCreate is the data for the InviteCreate event
type InviteCreate struct {
	GuildID   int64 `json:"guild_id,string"`
	ChannelID int64 `json:"channel_id,string"`

	Code      string    `json:"code"`
	CreatedAt Timestamp `json:"created_at"`

	MaxAge    int  `json:"max_age"`
	MaxUses   int  `json:"max_uses"`
	Temporary bool `json:"temporary"`
	Uses      int  `json:"uses"`

	Inviter *InviteUser `json:"inviter"`
}

// InviteDelete is the data for the InviteDelete event
type InviteDelete struct {
	GuildID   int64  `json:"guild_id,string"`
	ChannelID int64  `json:"channel_id,string"`
	Code      string `json:"code"`
}
