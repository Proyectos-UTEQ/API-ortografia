package utils

import "Proyectos-UTEQ/api-ortografia/internal/interfaces"

func SendNotification(notifier interfaces.Notifier, message string) error {
	return notifier.SendNotification(message)
}

func ResetPassword(notifier interfaces.Notifier, message string, url string) error {
	return notifier.ResetPassword(message, url)
}
