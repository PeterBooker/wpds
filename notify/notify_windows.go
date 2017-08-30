package notify

import (
	toast "github.com/jacobmarshall/go-toast"
)

// SendNotification sends a platform specific desktop notification.
func SendNotification(title string, message string) error {

	notification := toast.Notification{
		AppID:   "WPDS",
		Title:   title,
		Message: message,
	}

	return notification.Push()

}
