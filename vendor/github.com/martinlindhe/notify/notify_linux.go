package notify

import (
	"log"
	"os/exec"
)

// Notify displays a desktop notification
func Notify(appName string, title string, text string, iconPath string) {
	notification(appName, title, text, iconPath)
}

// Alert displays a desktop notification and plays a default system sound
func Alert(appName string, title string, text string, iconPath string) {
	notification(appName, title, text, iconPath)
	if err := exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga").Run(); err != nil {
		log.Println("error playing alert sound:", err)
	}
}

func notification(appName string, title string, text string, iconPath string) {
	cmd := exec.Command("notify-send", "-i", iconPath, title, text)
	cmd.Run()
}
