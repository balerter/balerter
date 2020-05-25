package dismock

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/json/option"

	"github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// Messages mocks a Messages request.
//
// This method will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func (m *Mocker) Messages(channelID discord.Snowflake, limit uint, messages []discord.Message) {
	if messages == nil {
		messages = []discord.Message{}
	}

	if len(messages) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent messages (%d vs. %d)", len(messages), limit))
	}

	const hardLimit uint = 100

	var after discord.Snowflake

	for i := 0; i <= len(messages)/int(hardLimit); i++ {
		var (
			from = uint(i) * hardLimit
			to   = uint(math.Min(float64(from+hardLimit), float64(len(messages))))

			fetch = to - from // we expect this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// expect either limit or hardlimit, depending on which is smaller, instead.
			if fetch < hardLimit {
				fetch = uint(math.Min(float64(limit), float64(hardLimit)))
			}

			limit -= fetch
		} else { // this means there is no limit, hence we should expect hardlimit
			fetch = hardLimit
		}

		m.messagesRange(channelID, 0, after, 0, fmt.Sprintf("Messages #%d", i+1), fetch, messages[from:to])

		if to-from < hardLimit {
			break
		}

		after = messages[to-1].ID
	}
}

// MessagesAround mocks a MessagesAround request.
//
// This method will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func (m *Mocker) MessagesAround(channelID, around discord.Snowflake, limit uint, messages []discord.Message) {
	switch {
	case limit == 0:
		limit = 50
	case limit > 100:
		limit = 100
	}

	if messages == nil {
		messages = []discord.Message{}
	}

	if len(messages) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent messages (%d vs. %d)", len(messages), limit))
	}

	m.messagesRange(channelID, 0, 0, around, "MessagesAround", limit, messages)
}

// MessagesBefore mocks a MessagesBefore request.
//
// This method will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func (m *Mocker) MessagesBefore(channelID, before discord.Snowflake, limit uint, messages []discord.Message) {
	if messages == nil {
		messages = []discord.Message{}
	}

	if len(messages) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent messages (%d vs. %d)", len(messages), limit))
	}

	const hardLimit = 100

	req := len(messages)/hardLimit + 1

	from := uint(math.Min(float64(uint(req)*hardLimit), float64(len(messages))))

	for i := req; i > 0; i-- {
		no := req - i + 1

		to := from
		from = uint(math.Max(float64(0), float64(int(to-hardLimit))))

		fetch := to - from // we expect this as the sent limit

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// expect either limit or hardlimit, depending on which is smaller, instead.
			if fetch < hardLimit {
				fetch = uint(math.Min(float64(limit), float64(hardLimit)))
			}

			limit -= fetch
		} else { // this means there is no limit, hence we should expect hardlimit
			fetch = hardLimit
		}

		m.messagesRange(channelID, before, 0, 0, fmt.Sprintf("MessagesBefore #%d", no), fetch, messages[from:to])

		if to-from < hardLimit {
			break
		}

		before = messages[from].ID
	}
}

// MessagesAfter mocks a MessagesAfter request.
//
// This method will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func (m *Mocker) MessagesAfter(channelID, after discord.Snowflake, limit uint, messages []discord.Message) {
	if messages == nil {
		messages = []discord.Message{}
	}

	if len(messages) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent messages (%d vs. %d)", len(messages), limit))
	}

	const hardLimit uint = 100

	for i := 0; i <= len(messages)/int(hardLimit); i++ {
		var (
			from = uint(i) * hardLimit
			to   = uint(math.Min(float64(from+hardLimit), float64(len(messages))))

			fetch = to - from // we expect this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// expect either limit or hardlimit, depending on which is smaller, instead.
			if fetch < hardLimit {
				fetch = uint(math.Min(float64(limit), float64(hardLimit)))
			}

			limit -= fetch
		} else { // this means there is no limit, hence we should expect hardlimit
			fetch = hardLimit
		}

		m.messagesRange(channelID, 0, after, 0, fmt.Sprintf("MessagesAfter #%d", i+1), fetch, messages[from:to])

		if to-from < hardLimit {
			break
		}

		after = messages[to-1].ID
	}
}

// messagesRange mocks a single request to the GET /messages endpoint.
//
// This method will sanitize Message.ID, Message.ChannelID and
// Message.Author.ID.
func (m *Mocker) messagesRange(
	channelID, before, after, around discord.Snowflake, name string, limit uint, messages []discord.Message,
) {
	for i, m := range messages {
		messages[i] = sanitize.Message(m, 1, channelID, 1)
	}

	m.MockAPI(name, http.MethodGet, "/channels/"+channelID.String()+"/messages",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.FormatUint(uint64(limit), 10)},
			}

			if after != 0 {
				expect["after"] = []string{after.String()}
			}

			if before != 0 {
				expect["before"] = []string{before.String()}
			}

			if around != 0 {
				expect["around"] = []string{around.String()}
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, messages)
		})
}

// Message mocks a Message request.
//
// The ID field and the ChannelID field of the passed Message must be set.
//
// This method will sanitize Message.Author.ID.
func (m *Mocker) Message(msg discord.Message) {
	msg = sanitize.Message(msg, 1, 1, 1)

	m.MockAPI("Message", http.MethodGet, "/channels/"+msg.ChannelID.String()+"/messages/"+msg.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, msg)
		})
}

// SendText mocks a SendText request.
//
// The ChannelID field and the Content field of the passed Message must be set.
//
// This method will sanitize Message.ID, Message.Author.ID, Message.Embeds.Type
// and Message.Embeds.Color.
func (m *Mocker) SendText(msg discord.Message) {
	m.sendMessageComplex("SendText", api.SendMessageData{
		Content: msg.Content,
	}, msg)
}

// SendEmbed mocks a SendEmbed request.
//
// The ChannelID field and the Embed field of the passed Message must be set.
//
// This method will sanitize Message.ID, Message.Author.ID, Message.Embeds.Type
// and Message.Embeds.Color.
func (m *Mocker) SendEmbed(msg discord.Message) {
	m.sendMessageComplex("SendEmbed", api.SendMessageData{
		Embed: &msg.Embeds[0],
	}, msg)
}

// SendMessage mocks a SendMessage request.
//
// The ChannelID field and the Content field of the passed Message must be set.
//
// This method will sanitize Message.ID, Message.Author.ID, Message.Embeds.Type
// and Message.Embeds.Color.
func (m *Mocker) SendMessage(embed *discord.Embed, msg discord.Message) {
	d := api.SendMessageData{
		Content: msg.Content,
	}

	if embed != nil {
		d.Embed = embed

		if len(msg.Embeds) == 0 {
			msg.Embeds = append(msg.Embeds, *d.Embed)
		}
	}

	m.sendMessageComplex("SendMessage", d, msg)
}

// EditText mocks a EditText request.
//
// The ID field, the ChannelID field and the Content field of the passed
// Message must be set.
//
// This method will sanitize Message.Author.ID.
func (m *Mocker) EditText(msg discord.Message) {
	m.editMessageComplex("EditText", api.EditMessageData{
		Content: option.NewNullableString(msg.Content),
	}, msg)
}

// EditEmbed mocks a EditEmbed request.
//
// The ID field, the ChannelID field and the Embed[0] field of the passed
// Message must be set.
//
// This method will sanitize Message.Author.ID, Message.Embed.Type and
// Message.Embed.Color.
func (m *Mocker) EditEmbed(msg discord.Message) {
	m.editMessageComplex("EditEmbed", api.EditMessageData{
		Embed: &msg.Embeds[0],
	}, msg)
}

// EditMessage mocks a EditMessage request.
//
// The ID field, the ChannelID field, the Content field of the passed Message
// must be set.
//
// This method will sanitize Message.Author.ID, Message.Embed.Type and
// Message.Embed.Color.
func (m *Mocker) EditMessage(embed *discord.Embed, msg discord.Message, suppressEmbeds bool) {
	d := api.EditMessageData{
		Content: option.NewNullableString(msg.Content),
		Embed:   embed,
	}

	if suppressEmbeds {
		d.Flags = &discord.SuppressEmbeds
	}

	m.editMessageComplex("EditMessage", d, msg)
}

// EditMessageComplex mocks a EditMessageComplex request.
//
// The ID field and the ChannelID field of the passed Message must be set.
//
// This method will sanitize Message.Author.ID, Message.Embed.Type and
// Message.Embed.Color.
func (m *Mocker) EditMessageComplex(d api.EditMessageData, msg discord.Message) {
	m.editMessageComplex("EditMessageComplex", d, msg)
}

// editMessageComplex mocks a EditMessageComplex request.
//
// The ID field and the ChannelID field of the passed Message must be set.
//
// This method will sanitize Message.Author.ID, Message.Embed.Type and
// Message.Embed.Color.
func (m *Mocker) editMessageComplex(name string, d api.EditMessageData, msg discord.Message) {
	msg = sanitize.Message(msg, 1, 1, 1)

	if d.Embed != nil {
		if d.Embed.Type == "" {
			d.Embed.Type = discord.NormalEmbed
		}

		if d.Embed.Color == 0 {
			d.Embed.Color = discord.DefaultEmbedColor
		}
	}

	m.MockAPI(name, http.MethodPatch, "/channels/"+msg.ChannelID.String()+"/messages/"+msg.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.EditMessageData), &d)
			mockutil.WriteJSON(t, w, msg)
		})
}

// DeleteMessage mocks a DeleteMessage request.
func (m *Mocker) DeleteMessage(channelID, messageID discord.Snowflake) {
	m.MockAPI("DeleteMessage", http.MethodDelete, "/channels/"+channelID.String()+"/messages/"+messageID.String(), nil)
}

type deleteMessagesPayload struct {
	Messages []discord.Snowflake `json:"messages"`
}

// DeleteMessages mocks a DeleteMessages request.
func (m *Mocker) DeleteMessages(channelID discord.Snowflake, messageIDs []discord.Snowflake) {
	m.MockAPI("DeleteMessages", http.MethodPost, "/channels/"+channelID.String()+"/messages/bulk-delete",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := deleteMessagesPayload{
				Messages: messageIDs,
			}

			mockutil.CheckJSON(t, r.Body, new(deleteMessagesPayload), &expect)
			w.WriteHeader(http.StatusNoContent)
		})
}
