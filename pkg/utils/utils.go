package utils

import (
	"time"
)

// IsToday checks if a string in the format "Monday, 2 January 2006" is today
func IsToday(dateStr string) bool {
	
    // Parse the date string into a time.Time object
    date, err := time.Parse("Monday, 2 January 2006", dateStr)
    if err != nil {
        return false
    }

    // Get the current time and truncate it to midnight
    now := time.Now().Truncate(24 * time.Hour)

    // Truncate the parsed date to midnight as well
    date = date.Truncate(24 * time.Hour)

    // Compare the two dates
    return now.Equal(date)
}