package notify

import (
	notifier "github.com/deckarep/gosx-notifier"
)

// SendNotification sends a platform specific desktop notification.
func SendNotification(title string, message string) {

	notification := notifier.Notification{
		Group:   "com.wpds.cli",
		Title:   title,
		Message: message,
		Sound:   notifier.Glass,
	}

	return notification.Push()

}
