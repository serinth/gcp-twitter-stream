package config

import "testing"
import "os"

func TestGetConfig(t *testing.T) {
	os.Chdir("../")
	_, err := GetConfig()
	if err != nil {
		t.Error(err)
	}
}
