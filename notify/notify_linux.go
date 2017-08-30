package notify

import (
	"os/exec"
)

// SendNotification sends a platform specific desktop notification.
func SendNotification(title string, message string) {

	exec.Command("notify-send", title, message)

}
