package notify

import (
	"log"

	toast "gopkg.in/toast.v1"
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
	note.Audio = toast.Default
	if err := note.Push(); err != nil {
		log.Println("ERROR:", err)
	}
}

func notification(appName string, title string, text string, iconPath string) toast.Notification {
	return toast.Notification{
		AppID:   appName,
		Title:   title,
		Message: text,
		Icon:    iconPath,
	}
}
