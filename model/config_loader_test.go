package model

import "testing"

func BenchmarkOptimiseIncompleteParamas(b *testing.B) {
	l, err := NewLoader(&ConfigLoaderOption{
		ExecPath:            "../",
		ConfigPath:          "../_example/test_auto_layout.yml",
		AppVersion:          "0.9.5",
		ConfigSchemaVersion: "1.1.0",
	})
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		l.optimiseIncompleteParams()
	}
}

func TestOptimiseIncompleteParamas(t *testing.T) {
	l, err := NewLoader(&ConfigLoaderOption{
		ExecPath:            "../",
		ConfigPath:          "../_example/test_auto_layout.yml",
		AppVersion:          "0.9.5",
		ConfigSchemaVersion: "1.1.0",
	})
	if err != nil {
		t.Fatalf("error:%v", err)
	}

	l.optimiseIncompleteParams()

	config := l.GetConfig()
	if len(config.Layout) == 0 {
		t.Fatalf("error:%v", config)
	}
	rH := config.Layout[0].Rows[0].Height
	if rH != "33%" {
		t.Fatalf("error:%v", rH)
	}
	rH = config.Layout[0].Rows[1].Height
	if rH != "33%" {
		t.Fatalf("error:%v", rH)
	}
	rH = config.Layout[0].Rows[2].Height
	if rH != "34%" {
		t.Fatalf("error:%v", rH)
	}
	rH = config.Layout[1].Rows[0].Height
	if rH != "100%" {
		t.Fatalf("error:%v", rH)
	}
	cW := config.Layout[0].Rows[0].Cols[0].Width
	if cW != 4 {
		t.Fatalf("error:%v", cW)
	}
	cW = config.Layout[0].Rows[0].Cols[1].Width
	if cW != 4 {
		t.Fatalf("error:%v", cW)
	}
	cW = config.Layout[0].Rows[0].Cols[2].Width
	if cW != 4 {
		t.Fatalf("error:%v", cW)
	}
	cW = config.Layout[0].Rows[1].Cols[0].Width
	if cW != 5 {
		t.Fatalf("error:%v", cW)
	}
	cW = config.Layout[0].Rows[1].Cols[1].Width
	if cW != 3 {
		t.Fatalf("error:%v", cW)
	}
	cW = config.Layout[0].Rows[1].Cols[2].Width
	if cW != 4 {
		t.Fatalf("error:%v", cW)
	}
	sH := config.Layout[0].Rows[0].Cols[0].Stacks[0].Height
	if sH != "33%" {
		t.Fatalf("error:%v", sH)
	}
	sH = config.Layout[0].Rows[0].Cols[0].Stacks[1].Height
	if sH != "33%" {
		t.Fatalf("error:%v", sH)
	}
	sH = config.Layout[0].Rows[0].Cols[0].Stacks[2].Height
	if sH != "34%" {
		t.Fatalf("error:%v", sH)
	}
	sH = config.Layout[0].Rows[2].Cols[0].Stacks[0].Height
	if sH != "100%" {
		t.Fatalf("error:%v", sH)
	}
}
