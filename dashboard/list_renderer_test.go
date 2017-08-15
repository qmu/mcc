package dashboard

import (
	"testing"

	"github.com/qmu/mcc/dashboard"
)

func TestRender(t *testing.T) {
	opt := &dashboard.ListRendererOption{
		Header:        buildHeader(),
		Body:          buildBody(),
		MaxH:          8,
		LineHighLight: true,
	}
	listRenderer := dashboard.NewListRenderer(opt)
	items := listRenderer.RenderActually()
	if items[len(items)-1] != "row4" {
		t.Error("invalid list display")
	}
	items = listRenderer.Move("bottom")
	if items[len(items)-1] != "[row9](fg-black,bg-green)" {
		t.Error("invalid list display")
	}
	if items[len(items)-1] != "[row9](fg-black,bg-green)" {
		t.Error("invalid list display")
	}
	items = listRenderer.Move("up")
	if items[len(items)-1] != "row9" {
		t.Error("invalid list display")
	}
	if items[len(items)-2] != "[row8](fg-black,bg-green)" {
		t.Error("invalid list display")
	}
	items = listRenderer.Move("up")
	items = listRenderer.Move("up")
	items = listRenderer.Move("up")
	if items[len(items)-1] != "row8" {
		t.Error("invalid list display")
	}
	items = listRenderer.Move("top")
	if items[len(items)-1] != "row4" {
		t.Error("invalid list display")
	}
	listRenderer.Move("bottom")
	items = listRenderer.ResetRender()
	if items[len(items)-1] != "row4" {
		t.Error("invalid list display")
	}
}

func stringify(list []string) (result string) {
	result += "\n"
	for _, l := range list {
		result += l + "\n"
	}
	return
}

func buildHeader() (result []string) {
	result = []string{
		"header1",
		"-----------------------------------------------",
	}
	return
}

func buildBody() (result []string) {
	result = []string{
		"row1",
		"row2",
		"row3",
		"row4",
		"row5",
		"row6",
		"row7",
		"row8",
		"row9",
	}
	return
}
