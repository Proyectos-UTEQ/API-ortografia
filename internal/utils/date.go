package utils

import "time"

func GetDate(date time.Time) string {
	return date.Format("02/01/2006")
}

func GetFullDate(date time.Time) string {
	return date.Format("02/01/2006 15:04:05")
}
