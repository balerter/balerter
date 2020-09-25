package webhook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookName(t *testing.T) {
	w := &Webhook{name: "foo"}
	assert.Equal(t, "foo", w.Name())
}
