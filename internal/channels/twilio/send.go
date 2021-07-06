package twilio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/message"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func (tw *TwilioVoice) Send(mes *message.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), tw.timeout)
	defer cancel()

	u := tw.apiPrefix + "/Accounts/" + tw.sid + "/Calls.json"

	twiml := tw.twiML
	twiml = strings.Replace(twiml, "{TEXT}", mes.Text, -1)
	if twiml == "" {
		twiml = mes.Text
	}

	buf := bytes.NewBuffer(nil)

	w := multipart.NewWriter(buf)
	if err := w.WriteField("From", tw.from); err != nil {
		return err
	}
	if err := w.WriteField("To", tw.to); err != nil {
		return err
	}
	if err := w.WriteField("Twiml", twiml); err != nil {
		return err
	}
	err := w.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.SetBasicAuth(tw.sid, tw.token)

	resp, err := tw.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read response body, %w", err)
	}

	tw.logger.Debug("twilio response", zap.ByteString("response", respBody))

	if resp.StatusCode != http.StatusCreated {
		tw.logger.Error("unexpected status code from twilio request",
			zap.Int("status", resp.StatusCode),
			zap.ByteString("body", respBody),
		)
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}
