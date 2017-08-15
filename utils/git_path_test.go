package utils

import (
	"github.com/qmu/mcc/utils"
	"testing"
)

func TestGetDotGitPath(t *testing.T) {
	result, err := utils.GetDotGitPath("/Users/tamurayoshiya/Sites/go/src/github.com/qmu/mcc/utils")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	t.Log(result)
}
