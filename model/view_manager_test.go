package model

import (
	"testing"

	ui "github.com/gizak/termui"
)

//
// func BenchmarkOptimiseIncompleteParamas(b *testing.B) {
// 	l, err := NewLoader(&ConfigLoaderOption{
// 		ExecPath:            "../",
// 		ConfigPath:          "../_example/test_auto_layout.yml",
// 		AppVersion:          "0.9.5",
// 		ConfigSchemaVersion: "1.1.0",
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	for i := 0; i < b.N; i++ {
// 		l.optimiseIncompleteParams()
// 	}
// }

func TestNewViewManager(t *testing.T) {
	// initialize termui
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	vm, err := NewViewManager(&ConfigLoaderOption{
		ExecPath:            "../",
		ConfigPath:          "../_example/example.yml",
		AppVersion:          "0.9.5",
		ConfigSchemaVersion: "1.1.0",
	})
	if err != nil {
		return
	}
	vm.SwitchTab(0)
	vm.HasWidget("github_issue")
	if vm.HasWidget("git_status") {
		vm.GetActiveWidgetsOf("git_status")
	}
	vm.NextWidget("right")
	vm.NextWidget("bottom")
	vm.GetGithubHost()
}
