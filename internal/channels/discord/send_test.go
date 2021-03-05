package discord

import (
	"bytes"
	"fmt"
	"github.com/balerter/balerter/internal/message"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/dismock/pkg/dismock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	successCases := []struct {
		name string
		data api.SendMessageData
	}{
		{
			name: "no files",
			data: api.SendMessageData{
				Content: "abc",
			},
		},
		{
			name: "with file",
			data: api.SendMessageData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("def"),
					},
				},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := dismock.NewSession(t)

				expect := sanitizeMessage(discord.Message{
					ChannelID: 123,
				}, 1, 1, 1)

				cp := c.data

				cp.Files = make([]api.SendMessageFile, len(c.data.Files))
				copy(cp.Files, c.data.Files) // the readers of the file will be consumed twice

				// the files are copied now, but the reader for them may be a pointer and wasn't
				// deep copied. therefore we create two readers using the data from the original
				// reader
				for i, f := range c.data.Files {
					b, err := ioutil.ReadAll(f.Reader)
					require.NoError(t, err)

					cp.Files[i].Reader = bytes.NewBuffer(b)
					c.data.Files[i].Reader = bytes.NewBuffer(b)
				}

				m.SendMessageComplex(c.data, expect)

				actual, err := s.SendMessageComplex(expect.ChannelID, cp)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)

				m.Eval()
			})
		}
	})

	failureCases := []struct {
		name  string
		data1 api.SendMessageData
		data2 api.SendMessageData
	}{
		{
			name: "different content",
			data1: api.SendMessageData{
				Content: "abc",
			},
			data2: api.SendMessageData{
				Content: "cba",
			},
		},
		{
			name: "different file",
			data1: api.SendMessageData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("def"),
					},
				},
			},
			data2: api.SendMessageData{
				Files: []api.SendMessageFile{
					{
						Name:   "abc",
						Reader: bytes.NewBufferString("fed"),
					},
				},
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				tMock := new(testing.T)

				m, s := dismock.NewSession(tMock)

				expect := sanitizeMessage(discord.Message{
					ChannelID: 123,
				}, 1, 1, 1)

				m.SendMessageComplex(c.data1, expect)

				actual, err := s.SendMessageComplex(expect.ChannelID, c.data2)
				require.NoError(t, err)

				assert.Equal(t, expect, *actual)
				assert.True(t, tMock.Failed())
			})
		}
	})
}

func sanitizeMessage(m discord.Message, id, channelID, authorID discord.Snowflake) discord.Message {
	if m.ID <= 0 {
		m.ID = id
	}

	if m.ChannelID <= 0 {
		m.ChannelID = channelID
	}

	m.Author = User(m.Author, authorID)

	return m
}

func User(u discord.User, id discord.Snowflake) discord.User {
	if u.ID <= 0 {
		u.ID = id
	}

	return u
}

type sessionMock struct {
	mock.Mock
}

func (m *sessionMock) SendMessage(channelID discord.Snowflake, content string, embed *discord.Embed) (*discord.Message, error) {
	args := m.Called(channelID, content, embed)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*discord.Message), args.Error(1)
}

func TestTest_error_send(t *testing.T) {
	m := &sessionMock{}
	m.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("err1"))
	d := &Discord{
		session: m,
	}
	mes := &message.Message{
		Text: "foo",
	}
	err := d.Send(mes)
	assert.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

func TestTest(t *testing.T) {
	m := &sessionMock{}
	m.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	d := &Discord{
		session: m,
	}
	mes := &message.Message{
		Text: "foo",
	}
	err := d.Send(mes)
	assert.NoError(t, err)
}
