package session

import (
	"strconv"
	"testing"
)

func TestRedis(t *testing.T) {

	r := New("localhost:6379", "", 0)

	t.Run("value", func(t *testing.T) {
		key := "test"
		expectedVal := "1"
		r.SetValue(key, expectedVal, 0)
		actualVal, err := r.GetValue(key)
		if err != nil {
			t.Error(err)
		}
		if expectedVal != actualVal {
			t.Error("Expected val does not equal actual val")
		}
	})

	t.Run("list", func(t *testing.T) {
		key := "nums"
		count := 10
		for i := 0; i < count; i++ {
			if err := r.AddToList(key, strconv.Itoa(i)); err != nil {
				t.Error(err)
			}
		}
		res, err := r.GetList(key)
		if err != nil {
			t.Error(err)
		}
		if len(res) != count {
			t.Error("list not set properly")
		}
	})
}
