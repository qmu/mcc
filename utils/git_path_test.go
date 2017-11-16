package utils

import (
	"testing"
)

func TestGetDotGitPath(t *testing.T) {
	result, err := GetDotGitPath("/Users/tamurayoshiya/Sites/go/src/github.com/qmu/mcc/utils")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	t.Log(result)
}
