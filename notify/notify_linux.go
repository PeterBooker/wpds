package notify

import (
	"os/exec"
	"time"
)

const (
	notificationBreakpoint = 5 * time.Minute
)

// SendNotification sends a platform specific desktop notification.
func SendNotification(title string, message string, elapsed time.Duration) {

	if elapsed < notificationBreakpoint {
		return
	}

	exec.Command("notify-send", title, message)

	return

}
