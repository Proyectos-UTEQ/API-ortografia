package utils

import (
	"fmt"
	"time"
)

func GetDate(date time.Time) string {
	return date.Format("02/01/2006")
}

// GetFullDateOrNull Retorna la fecha en formato string en caso de ser nil se retorna el nil.
func GetFullDateOrNull(date *time.Time) *string {
	if date == nil {
		return nil
	} else {
		dateString := date.Format("02/01/2006 15:04:05")
		return &dateString
	}
}

func GetFullDate(date time.Time) string {
	return date.Format("02/01/2006 15:04:05")
}

func GetDateOrNullForString(date *string) *string {
	if date == nil {
		return nil
	}

	t, err := time.Parse(time.RFC3339, *date)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	dateString := t.Format("02/01/2006")
	return &dateString
}
