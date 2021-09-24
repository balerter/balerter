package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SendPhotoMessage send PhotoMessage to the Telegram API
func (api *API) SendPhotoMessage(mes *PhotoMessage) error {
	body, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("error marshaling a message, %w", err)
	}

	return api.sendMessage(body, methodSendPhoto)
}

// SendTextMessage send TextMessage to the Telegram API
func (api *API) SendTextMessage(mes *TextMessage) error {
	body, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("error marshaling a message, %w", err)
	}

	return api.sendMessage(body, methodSendMessage)
}

type tgResponse struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

func (api *API) sendMessage(body []byte, method string) error {
	req, err := http.NewRequest(http.MethodPost, api.endpoint+method, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error generate request to telegram, %w", err)
	}
	req.Header.Add("Content-type", "application/json")

	res, err := api.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error send request, %w", err)
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error read response body, %w", err)
	}

	tgResp := &tgResponse{}
	err = json.Unmarshal(body, tgResp)
	if err != nil {
		return fmt.Errorf("error unmarshal response body, %w", err)
	}

	if !tgResp.OK {
		return fmt.Errorf("error send message %d: %s", tgResp.ErrorCode, tgResp.Description)
	}

	return nil
}
