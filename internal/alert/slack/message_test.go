package slack

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSlack_createSlackMessage(t *testing.T) {

	// reference
	_ = slackMessage{
		Text: "alert-text",
		Blocks: []block{
			{
				Type: "section",
				Section: &blockText{
					Type:  "mrkdwn",
					Emoji: false,
					Text:  "alert-text",
				},
				Context: nil,
			},
			{
				Type:    "context",
				Section: nil,
				Context: &[]blockText{
					{
						Type:  "mrkdwn",
						Emoji: false,
						Text:  "Alert name: alertName",
					},
				},
			},
		},
	}

	mes := createSlackMessage("alert-name", "alert-text")
	assert.Equal(t, "alert-text", mes.Text)
	require.Equal(t, 2, len(mes.Blocks))
	b := mes.Blocks[0]
	require.NotNil(t, b.Section)
	assert.Equal(t, "mrkdwn", b.Section.Type)
	assert.Equal(t, false, b.Section.Emoji)
	assert.Equal(t, "alert-text", b.Section.Text)
	assert.Nil(t, b.Context)

	b = mes.Blocks[1]
	assert.Nil(t, b.Section)
	require.NotNil(t, b.Context)
	assert.Equal(t, "context", b.Type)

	require.Equal(t, 1, len(*b.Context))
	cb := (*(b.Context))[0]
	assert.Equal(t, "mrkdwn", cb.Type)
	assert.Equal(t, false, cb.Emoji)
	assert.Equal(t, "Alert name: alert-name", cb.Text)
}
