package notify

import (
	notifier "github.com/deckarep/gosx-notifier"
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

	notification := notifier.Notification{
		Group:   "com.wpds.cli",
		Title:   title,
		Message: message,
		Sound:   notifier.Glass,
	}

	notification.Push()

	return

}
