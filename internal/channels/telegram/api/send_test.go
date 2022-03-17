package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPI_SendTextMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		errParseForm := req.ParseMultipartForm(1024)
		require.NoError(t, errParseForm)

		assert.Equal(t, "1", req.Form.Get("chat_id"))
		assert.Equal(t, "a", req.Form.Get("text"))
		assert.Equal(t, "MarkdownV2", req.Form.Get("parse_mode"))

		resp := `{"ok":true}`
		rw.Write([]byte(resp))
	}))
	defer srv.Close()

	a := &API{
		endpoint:   srv.URL + "/",
		httpClient: &http.Client{},
	}

	err := a.SendTextMessage(&TextMessage{
		ChatID: 1,
		Text:   "a",
	})

	require.NoError(t, err)
}

func TestAPI_SendPhotoMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		errParseForm := req.ParseMultipartForm(1024)
		require.NoError(t, errParseForm)

		assert.Equal(t, "1", req.Form.Get("chat_id"))
		assert.Equal(t, "b", req.Form.Get("caption"))
		assert.Equal(t, "MarkdownV2", req.Form.Get("parse_mode"))
		f, ok := req.MultipartForm.File["photo"]
		require.True(t, ok)
		require.Equal(t, 1, len(f))
		assert.Equal(t, "image.png", f[0].Filename)

		resp := `{"ok":true}`
		rw.Write([]byte(resp))
	}))
	defer srv.Close()

	a := &API{
		endpoint:   srv.URL + "/",
		httpClient: &http.Client{},
	}

	err := a.SendPhotoMessage(&PhotoMessage{
		ChatID:  1,
		Photo:   "a",
		Caption: "b",
	})

	require.NoError(t, err)
}
