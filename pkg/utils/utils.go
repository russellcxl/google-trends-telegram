package utils

import (
	"encoding/json"
	"os"
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

func ReadJSONFile(filename string, data interface{}) error {
    // Read the file contents
    content, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    // Unmarshal the JSON data
    err = json.Unmarshal(content, data)
    if err != nil {
        return err
    }
    return nil
}


func WriteJSONFile(filename string, data interface{}) error {
    // Open the file for writing
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    // Marshal the data into JSON
    jsonData, err := json.MarshalIndent(data, " ", " ")
    if err != nil {
        return err
    }
    // Write the JSON data to the file
    _, err = file.Write(jsonData)
    if err != nil {
        return err
    }
    return nil
}





