package model

import "testing"

func BenchmarkValidateLayout(b *testing.B) {
	v, err := NewConfigValidator()
	if err != nil {
		return
	}

	// vErrInvalidTotalStackHeight
	conf := ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID: "hoge",
			},
		},
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "50%",
						Cols: []colNode{
							colNode{
								Width: 7,
								Stacks: []stackNode{
									stackNode{
										ID:     "hoge",
										Height: "80%",
									},
									stackNode{
										ID:     "hoge",
										Height: "30%",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for i := 0; i < b.N; i++ {
		v.validateLayout(&conf)
	}
}

func BenchmarkValidateWidgets(b *testing.B) {
	v, err := NewConfigValidator()
	if err != nil {
		return
	}

	// vErrLackOfDockerStatusInvalidMetrics
	conf := ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "docker_status",
				Content: []interface{}{
					map[interface{}]interface{}{
						"name":      "name1",
						"metrics":   "hoge",
						"container": "hoge_container",
					},
				},
			},
		},
	}
	for i := 0; i < b.N; i++ {
		v.validateWidgets(&conf)
	}
}
func TestValidate(t *testing.T) {
	v, _ := NewConfigValidator()
	conf := ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "tail_file",
			},
			widgetNode{
				ID:    "widget2",
				Title: "widget2",
				Type:  "text_file",
				Path:  "./",
			},
		},
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "50%",
						Cols: []colNode{
							colNode{
								Width: 7,
								Stacks: []stackNode{
									stackNode{
										ID:     "widget1",
										Height: "80%",
									},
									stackNode{
										ID:     "widget2",
										Height: "30%",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	vErrs, err := v.validate(&conf)
	if err != nil {
		return
	}
	if vErrs[0].message != vErrLackOfTailFilePath {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}
	if vErrs[1].message != vErrInvalidTotalStackHeight {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[1].message, err)
	}
}

func TestValidateWidgets(t *testing.T) {
	v, _ := NewConfigValidator()

	// vErrLackOfWidgetID
	conf := ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				Type:  "note",
				Title: "hoge",
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfWidgetID {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfWidgetType
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "hoge",
				Title: "hoge",
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfWidgetType {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfWidgetTitle
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:   "hoge",
				Type: "note",
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfWidgetTitle {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfNoteContent
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "note",
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfNoteContent {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfTextFilePath
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "text_file",
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfTextFilePath {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfDockerStatusContent
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "docker_status",
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfDockerStatusContent {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfDockerStatusName
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "docker_status",
				Content: []interface{}{
					map[interface{}]interface{}{
						"metrics":   "cpu",
						"container": "hoge_container",
					},
				},
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfDockerStatusName {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfDockerStatusContainer
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "docker_status",
				Content: []interface{}{
					map[interface{}]interface{}{
						"name":    "name1",
						"metrics": "cpu",
					},
				},
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfDockerStatusContainer {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfDockerStatusContainer
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "docker_status",
				Content: []interface{}{
					map[interface{}]interface{}{
						"name":      "name1",
						"container": "hoge_container",
					},
				},
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfDockerStatusContainer {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfDockerStatusInvalidMetrics
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "docker_status",
				Content: []interface{}{
					map[interface{}]interface{}{
						"name":      "name1",
						"metrics":   "hoge",
						"container": "hoge_container",
					},
				},
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfDockerStatusInvalidMetrics {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfMenuContent
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "menu",
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfMenuContent {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfMenuName
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "menu",
				Content: []interface{}{
					map[interface{}]interface{}{
						"category":    "hoge",
						"description": "hoge",
						"command":     "ls",
					},
				},
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfMenuName {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfMenuCategory
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "menu",
				Content: []interface{}{
					map[interface{}]interface{}{
						"name":        "hoge",
						"description": "hoge",
						"command":     "ls",
					},
				},
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfMenuCategory {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfMenuDescription
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "menu",
				Content: []interface{}{
					map[interface{}]interface{}{
						"name":     "hoge",
						"category": "hoge",
						"command":  "ls",
					},
				},
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfMenuDescription {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfMenuCommand
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "menu",
				Content: []interface{}{
					map[interface{}]interface{}{
						"name":        "hoge",
						"category":    "hoge",
						"description": "hoge",
					},
				},
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfMenuCommand {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfGithubIssueRegex
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "github_issue",
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfGithubIssueRegex {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfTailFilePath
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID:    "widget1",
				Title: "widget1",
				Type:  "tail_file",
			},
		},
	}
	if vErrs, err := v.validateWidgets(&conf); vErrs[0].message != vErrLackOfTailFilePath {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}
}

func TestValidateLayout(t *testing.T) {
	v, _ := NewConfigValidator()

	// vErrLackOfTabs
	conf := ConfRoot{}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrLackOfTabs {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfTabName
	conf = ConfRoot{
		Layout: []tabNode{
			tabNode{},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrLackOfTabName {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfTabRows
	conf = ConfRoot{
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrLackOfTabRows {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrRowHeightInvalid
	conf = ConfRoot{
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "101%",
						Cols:   []colNode{},
					},
				},
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrRowHeightInvalid {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfRowCols
	conf = ConfRoot{
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "100%",
					},
				},
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrLackOfRowCols {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrColWidthInvalid
	conf = ConfRoot{
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "100%",
						Cols: []colNode{
							colNode{
								Width: 13,
							},
						},
					},
				},
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrColWidthInvalid {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfColStacks
	conf = ConfRoot{
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "100%",
						Cols: []colNode{
							colNode{
								Width: 12,
							},
						},
					},
				},
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrLackOfColStacks {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrLackOfStackID
	conf = ConfRoot{
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "100%",
						Cols: []colNode{
							colNode{
								Width: 12,
								Stacks: []stackNode{
									stackNode{
										Height: "100%",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrLackOfStackID {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrWidgetDoesNotExit
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID: "hoge",
			},
		},
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "100%",
						Cols: []colNode{
							colNode{
								Width: 12,
								Stacks: []stackNode{
									stackNode{
										ID:     "fuga",
										Height: "100%",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrWidgetDoesNotExit {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrWidgetDoesNotExit
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID: "hoge",
			},
		},
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "100%",
						Cols: []colNode{
							colNode{
								Width: 12,
								Stacks: []stackNode{
									stackNode{
										ID:     "hoge",
										Height: "101%",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrWidgetDoesNotExit {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrInvalidTotalRowHeight
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID: "hoge",
			},
		},
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "50%",
						Cols: []colNode{
							colNode{
								Width: 12,
								Stacks: []stackNode{
									stackNode{
										ID: "hoge",
									},
								},
							},
						},
					},
					rowNode{
						Height: "60%",
						Cols: []colNode{
							colNode{
								Width: 12,
								Stacks: []stackNode{
									stackNode{
										ID: "hoge",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrInvalidTotalRowHeight {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrInvalidTotalColWidth
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID: "hoge",
			},
		},
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "50%",
						Cols: []colNode{
							colNode{
								Width: 7,
								Stacks: []stackNode{
									stackNode{
										ID: "hoge",
									},
								},
							},
							colNode{
								Width: 6,
								Stacks: []stackNode{
									stackNode{
										ID: "hoge",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrInvalidTotalColWidth {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}

	// vErrInvalidTotalStackHeight
	conf = ConfRoot{
		Widgets: []widgetNode{
			widgetNode{
				ID: "hoge",
			},
		},
		Layout: []tabNode{
			tabNode{
				Name: "Tab1",
				Rows: []rowNode{
					rowNode{
						Height: "50%",
						Cols: []colNode{
							colNode{
								Width: 7,
								Stacks: []stackNode{
									stackNode{
										ID:     "hoge",
										Height: "80%",
									},
									stackNode{
										ID:     "hoge",
										Height: "30%",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if vErrs, err := v.validateLayout(&conf); vErrs[0].message != vErrInvalidTotalStackHeight {
		t.Fatalf("Get validation error: %v | error:%v", vErrs[0].message, err)
	}
}
