package notify

import (
	"log"

	"github.com/deckarep/gosx-notifier"
)

// Notify displays a desktop notification
func Notify(appName string, title string, text string, iconPath string) {
	note := notification(appName, title, text, iconPath)
	if err := note.Push(); err != nil {
		log.Println("ERROR:", err)
	}
}

// Alert displays a desktop notification and plays a default system sound
func Alert(appName string, title string, text string, iconPath string) {
	note := notification(appName, title, text, iconPath)
	note.Sound = gosxnotifier.Default
	if err := note.Push(); err != nil {
		log.Println("ERROR:", err)
	}
}

func notification(appName string, title string, text string, iconPath string) *gosxnotifier.Notification {
	head := ""
	if text == "" {
		head = title
		title = ""
	} else {
		head = text
	}
	note := gosxnotifier.NewNotification(head)
	note.Title = appName
	note.Subtitle = title
	note.AppIcon = iconPath // (10.9+ ONLY)
	return note
}
