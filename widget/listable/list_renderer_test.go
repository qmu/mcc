package listable

import "testing"

func BenchmarkRender(b *testing.B) {
	opt := &ListRendererOption{
		Header:        buildHeader(),
		Body:          buildBody(),
		MaxH:          8,
		LineHighLight: true,
	}
	listRenderer := NewListRenderer(opt)
	for i := 0; i < b.N; i++ {
		listRenderer.RenderActually()
		listRenderer.MoveCursorWithFocus("bottom")
		listRenderer.MoveCursorWithFocus("up")
		listRenderer.Deactivate()
	}
}

func TestRender(t *testing.T) {
	opt := &ListRendererOption{
		Header:        buildHeader(),
		Body:          buildBody(),
		MaxH:          8,
		LineHighLight: true,
	}
	listRenderer := NewListRenderer(opt)
	items := listRenderer.RenderActually()
	if items[len(items)-1] != "row4" {
		t.Fatalf("invalid list display %v", items[len(items)-1])
	}
	items = listRenderer.MoveCursorWithFocus("bottom")
	if items[len(items)-1] != "[row9](fg-black,bg-green)" {
		t.Fatalf("invalid list display %v", items[len(items)-1])
	}
	if items[len(items)-1] != "[row9](fg-black,bg-green)" {
		t.Fatalf("invalid list display %v", items[len(items)-1])
	}
	items = listRenderer.MoveCursorWithFocus("up")
	if items[len(items)-1] != "row9" {
		t.Fatalf("invalid list display %v", items[len(items)-1])
	}
	if items[len(items)-2] != "[row8](fg-black,bg-green)" {
		t.Fatalf("invalid list display %v", items[len(items)-1])
	}
	items = listRenderer.MoveCursorWithFocus("up")
	items = listRenderer.MoveCursorWithFocus("up")
	items = listRenderer.MoveCursorWithFocus("up")
	if items[len(items)-1] != "row8" {
		t.Fatalf("invalid list display %v", items[len(items)-1])
	}
	items = listRenderer.MoveCursorWithFocus("top")
	if items[len(items)-1] != "row4" {
		t.Fatalf("invalid list display %v", items[len(items)-1])
	}
	listRenderer.MoveCursorWithFocus("bottom")
	items = listRenderer.Deactivate()
	if items[len(items)-1] != "row9" {
		t.Fatalf("invalid list display %v", items[len(items)-1])
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
