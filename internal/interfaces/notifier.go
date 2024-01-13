package interfaces

type Notifier interface {
	SendNotification(message string) error
	ResetPassword(message string, url string) error
}
