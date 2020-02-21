package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (api *API) SendPhotoMessage(mes *PhotoMessage) error {
	body, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("error marshaling a message, %w", err)
	}

	return api.sendMessage(body, methodSendPhoto)
}

func (api *API) SendTextMessage(mes *TextMessage) error {
	body, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("error marshaling a message, %w", err)
	}

	return api.sendMessage(body, methodSendMessage)
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

	res.Body.Close()

	return nil
}
