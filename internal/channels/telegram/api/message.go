package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type TextMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type PhotoMessage struct {
	ChatID  int64  `json:"chat_id"`
	Photo   string `json:"photo"`
	Caption string `json:"caption,omitempty"`
}

func (message *PhotoMessage) getPhotoFile() (fileName string, err error) {
	output, err := ioutil.TempFile(os.TempDir(), "balerter.*")
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	fileName = output.Name()
	defer output.Close()

	response, err := http.Get(message.Photo)
	if err != nil {
		fmt.Println("Error while downloading", message.Photo, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", message.Photo, "-", err)
		return
	}

	return
}

func NewPhotoMessage(chatID int64, photo, text string) *PhotoMessage {
	return &PhotoMessage{
		ChatID:  chatID,
		Photo:   photo,
		Caption: text,
	}
}

func NewTextMessage(chatID int64, text string) *TextMessage {
	return &TextMessage{
		ChatID: chatID,
		Text:   text,
	}
}
