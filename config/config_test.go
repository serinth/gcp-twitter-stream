package config

import (
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	cwd, err := os.Getwd()
	_, err2 := GetConfig(cwd + "/config.json")
	if err2 != nil {
		t.Error(err)
	}
}
