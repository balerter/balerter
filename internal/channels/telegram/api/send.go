package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

const (
	methodSendMessage = "sendMessage"
	methodSendPhoto   = "sendPhoto"
)

// SendPhotoMessage send PhotoMessage to the Telegram API
func (api *API) SendPhotoMessage(mes *PhotoMessage) error {
	fields := map[string]string{
		"chat_id":    strconv.Itoa(int(mes.ChatID)),
		"photo":      mes.Photo,
		"caption":    mes.Caption,
		"parse_mode": "MarkdownV2",
	}

	return api.sendMessage(fields, methodSendPhoto)
}

// SendTextMessage send TextMessage to the Telegram API
func (api *API) SendTextMessage(mes *TextMessage) error {
	fields := map[string]string{
		"chat_id":    strconv.Itoa(int(mes.ChatID)),
		"text":       mes.Text,
		"parse_mode": "MarkdownV2",
	}

	return api.sendMessage(fields, methodSendMessage)
}

func (api *API) buildMultipartBody(fields map[string]string) (string, *bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)

	w := multipart.NewWriter(buf)

	for k, v := range fields {
		var f io.Writer
		var e error

		switch k {
		case "photo":
			if strings.HasPrefix(v, "http") {
				f, e = w.CreateFormField(k)
			} else {
				f, e = w.CreateFormFile(k, "image.png")
			}
		default:
			f, e = w.CreateFormField(k)
		}
		if e != nil {
			return "", nil, fmt.Errorf("error create field %s, %w", k, e)
		}
		_, errCopy := io.Copy(f, bytes.NewBufferString(v))
		if errCopy != nil {
			return "", nil, fmt.Errorf("error copy value for field %s", k)
		}
	}

	errClose := w.Close()
	if errClose != nil {
		return "", nil, fmt.Errorf("error close writer, %w", errClose)
	}

	return w.FormDataContentType(), buf, nil
}

type tgResponse struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

func (api *API) sendMessage(fields map[string]string, method string) error {
	contentType, body, errBuildBody := api.buildMultipartBody(fields)
	if errBuildBody != nil {
		return fmt.Errorf("error build body, %w", errBuildBody)
	}

	req, errCreateRequest := http.NewRequest(http.MethodPost, api.endpoint+method, body)
	if errCreateRequest != nil {
		return fmt.Errorf("error generate request to telegram, %w", errCreateRequest)
	}
	req.Header.Add("Content-type", contentType)

	res, errDo := api.httpClient.Do(req)
	if errDo != nil {
		return fmt.Errorf("error send request, %w", errDo)
	}

	defer res.Body.Close()

	respBody, errReadRespBody := io.ReadAll(res.Body)
	if errReadRespBody != nil {
		return fmt.Errorf("error read response body, %w", errReadRespBody)
	}

	tgResp := &tgResponse{}
	errDecodeRespBody := json.Unmarshal(respBody, tgResp)
	if errDecodeRespBody != nil {
		return fmt.Errorf("error unmarshal response body, %w", errDecodeRespBody)
	}

	if !tgResp.OK {
		return fmt.Errorf("error send message %d: %s", tgResp.ErrorCode, tgResp.Description)
	}

	return nil
}
