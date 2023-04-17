package utils

import (
	"fmt"
	"testing"
)

type userIDs struct {
	UserIDs []int64 `json:"user_ids"`
}

func TestUtils(t *testing.T) {
	var ids userIDs
	ReadJSONFile("../../data/allowed_users.json", &ids)
	fmt.Printf("%+v", ids)
}