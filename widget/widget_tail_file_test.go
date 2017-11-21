package widget

import (
	"path/filepath"
	"testing"
)

func BenchmarkNewTailFileWidget(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d, err := filepath.Abs("../") // current dir is ./widget
		if err != nil {
			return
		}
		_, err = NewTailFileWidget(&Option{
			ExecPath: d,
			Height:   100,
			Width:    100,
			Title:    "test",
			Type:     "tail_file",
			Path:     d + "/_example/example.log",
		})
		if err != nil {
			return
		}
	}
}

func TestNewTailFileWidget(t *testing.T) {
	d, err := filepath.Abs("../") // current dir is ./widget
	if err != nil {
		return
	}
	_, err = NewTailFileWidget(&Option{
		ExecPath: d,
		Height:   100,
		Width:    100,
		Title:    "test",
		Type:     "tail_file",
		Path:     d + "/_example/example.log",
	})
	if err != nil {
		return
	}
}
