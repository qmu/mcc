package model

import (
	"strconv"

	m2s "github.com/mitchellh/mapstructure"
	"github.com/qmu/mcc/utils"
	"github.com/qmu/mcc/widget"
)

const (
	// widgets section
	vErrLackOfWidgetID                   = "'widgets[].id' should have value"
	vErrLackOfWidgetType                 = "'widgets[].type' should have value"
	vErrLackOfWidgetTitle                = "'widgets[].title' should have value"
	vErrLackOfNoteContent                = "'widgets[].type=note' should have content"
	vErrLackOfTextFilePath               = "'widgets[].type=text_file' should have path"
	vErrLackOfDockerStatusContent        = "'widgets[].type=docker_status' should have content"
	vErrLackOfDockerStatusName           = "'widgets[].type=docker_status' should have value of content[].name"
	vErrLackOfDockerStatusContainer      = "'widgets[].type=docker_status' should have value of content[].container"
	vErrLackOfDockerStatusMetrics        = "'widgets[].type=docker_status' should have value of content[].metrics"
	vErrLackOfDockerStatusInvalidMetrics = "'widgets[].type=docker_status' metrics should be 'cpu' or 'memory'"
	vErrLackOfMenuContent                = "'widgets[].type=menu' should have content"
	vErrLackOfMenuName                   = "'widgets[].type=menu' should have value of content[].name"
	vErrLackOfMenuCategory               = "'widgets[].type=menu' should have value of content[].category"
	vErrLackOfMenuDescription            = "'widgets[].type=menu' should have value of content[].description"
	vErrLackOfMenuCommand                = "'widgets[].type=menu' should have value of content[].command"
	vErrLackOfGithubIssueRegex           = "'widgets[].type=github_issue' should have issue_regex"
	vErrLackOfTailFilePath               = "'widgets[].type=tail_file' should have path"
	// layout section
	vErrLackOfTabs              = "'layout should have array of tab"
	vErrLackOfTabName           = "'layout[].name' should have value"
	vErrLackOfTabRows           = "'layout[].rows' should have array of row"
	vErrRowHeightInvalid        = "'layout[].rows[].height' should be '1%' ~ '100%'"
	vErrLackOfRowCols           = "'layout[].rows[].cols' should have array of col"
	vErrColWidthInvalid         = "'layout[].rows[].cols[].width' should be 1~12"
	vErrLackOfColStacks         = "'layout[].rows[].cols[].stacks' should have array of stack"
	vErrLackOfStackID           = "'layout[].rows[].cols[].stacks[].id' should have value"
	vErrWidgetDoesNotExit       = "'layout[].rows[].cols[].stacks[].id' should be defined in widgets[].id"
	vErrStackHeightInvalid      = "'layout[].rows[].cols[].stacks[].height' should be '1%' ~ '100%'"
	vErrInvalidTotalRowHeight   = "total of the layout[].rows[].height in a tab should be <= 100%"
	vErrInvalidTotalColWidth    = "total of the layout[].rows[].cols[].width in a row should be <= 12"
	vErrInvalidTotalStackHeight = "total of the layout[].rows[].cols[].stacks[].height in a col should be <= 100"
)

// ConfigValidator is
type ConfigValidator struct {
}

// validationError is
type validationError struct {
	condition bool
	message   string
	position  string
}

// NewConfigValidator constructs a ConfigLoader
func NewConfigValidator() (c *ConfigValidator, err error) {
	c = new(ConfigValidator)
	return
}

func (c *ConfigValidator) validate(config *ConfRoot) (vErr []*validationError, err error) {
	verr1, err := c.validateWidgets(config)
	if err != nil {
		return
	}
	verr2, err := c.validateLayout(config)
	if err != nil {
		return
	}
	vErr = append(verr1, verr2...)
	return
}

func (c *ConfigValidator) validateWidgets(config *ConfRoot) (vErr []*validationError, err error) {
	vErr = nil
	for i1, w := range config.Widgets {
		// all widgetNode should have id, type, title
		if w.ID == "" {
			vErr = append(vErr, &validationError{
				message:  vErrLackOfWidgetID,
				position: "widgets[" + strconv.Itoa(i1) + "]",
			})
		}
		if w.Type == "" {
			vErr = append(vErr, &validationError{
				message:  vErrLackOfWidgetType,
				position: "widgets[" + strconv.Itoa(i1) + "]",
			})
		}
		if w.Title == "" {
			vErr = append(vErr, &validationError{
				message:  vErrLackOfWidgetTitle,
				position: "widgets[" + strconv.Itoa(i1) + "]",
			})
		}
		// type=note widget, should have "content"
		if w.Type == "note" && w.Content == nil {
			vErr = append(vErr, &validationError{
				message:  vErrLackOfNoteContent,
				position: "widgets[" + strconv.Itoa(i1) + "]",
			})
		}
		// type=text_file widget, should have "path"
		if w.Type == "text_file" && w.Path == "" {
			vErr = append(vErr, &validationError{
				message:  vErrLackOfTextFilePath,
				position: "widgets[" + strconv.Itoa(i1) + "]",
			})
		}
		if w.Type == "docker_status" {
			// type=docker_status widget, should have "content"
			if w.Content == nil {
				vErr = append(vErr, &validationError{
					message:  vErrLackOfDockerStatusContent,
					position: "widgets[" + strconv.Itoa(i1) + "]",
				})
			} else {
				containers := &[]widget.Container{}
				if err = m2s.Decode(w.Content, containers); err != nil {
					return
				}
				// type=docker_status widget, should have "metrics","name", "container" in the "content"
				for _, ct := range *containers {
					if ct.Name == "" {
						vErr = append(vErr, &validationError{
							message:  vErrLackOfDockerStatusName,
							position: "widgets[" + strconv.Itoa(i1) + "]",
						})
					}
					if ct.Container == "" {
						vErr = append(vErr, &validationError{
							message:  vErrLackOfDockerStatusContainer,
							position: "widgets[" + strconv.Itoa(i1) + "]",
						})
					}
					if ct.Metrics == "" {
						vErr = append(vErr, &validationError{
							message:  vErrLackOfDockerStatusContainer,
							position: "widgets[" + strconv.Itoa(i1) + "]",
						})
					}
					// type=docker_status widget, "metrics" should be "cpu" or "memory"
					if ct.Metrics != "cpu" && ct.Metrics != "memory" {
						vErr = append(vErr, &validationError{
							message:  vErrLackOfDockerStatusInvalidMetrics,
							position: "widgets[" + strconv.Itoa(i1) + "]",
						})
					}
				}
			}
		}
		if w.Type == "menu" {
			// type=menu widget, should have "content"
			if w.Content == nil {
				vErr = append(vErr, &validationError{
					message:  vErrLackOfMenuContent,
					position: "widgets[" + strconv.Itoa(i1) + "]",
				})
			} else {
				// type=menu widget, "content" should have "category", "name", "description", "command"
				menus := &[]widget.Menu{}
				if err = m2s.Decode(w.Content, menus); err != nil {
					return
				}
				for _, m := range *menus {
					if m.Name == "" {
						vErr = append(vErr, &validationError{
							message:  vErrLackOfMenuName,
							position: "widgets[" + strconv.Itoa(i1) + "]",
						})
					}
					if m.Category == "" {
						vErr = append(vErr, &validationError{
							message:  vErrLackOfMenuCategory,
							position: "widgets[" + strconv.Itoa(i1) + "]",
						})
					}
					if m.Description == "" {
						vErr = append(vErr, &validationError{
							message:  vErrLackOfMenuDescription,
							position: "widgets[" + strconv.Itoa(i1) + "]",
						})
					}
					if m.Command == "" {
						vErr = append(vErr, &validationError{
							message:  vErrLackOfMenuCommand,
							position: "widgets[" + strconv.Itoa(i1) + "]",
						})
					}
				}
			}
		}
		if w.Type == "github_issue" {
			// type=github_issue widget, should have "issue_regex"
			if w.IssueRegex == "" {
				vErr = append(vErr, &validationError{
					message:  vErrLackOfGithubIssueRegex,
					position: "widgets[" + strconv.Itoa(i1) + "]",
				})
			}
		}
		if w.Type == "tail_file" {
			// type=tail_file widget, should have "path"
			if w.Path == "" {
				vErr = append(vErr, &validationError{
					message:  vErrLackOfTailFilePath,
					position: "widgets[" + strconv.Itoa(i1) + "]",
				})
			}
		}
	}
	return
}

func (c *ConfigValidator) validateLayout(config *ConfRoot) (vErr []*validationError, err error) {
	vErr = nil
	// ConfRoot.Layout should be set
	if len(config.Layout) == 0 {
		vErr = append(vErr, &validationError{
			message:  vErrLackOfTabs,
			position: "root",
		})
	}
	for i1, t := range config.Layout {
		// tabNode.Name should be set
		if t.Name == "" {
			vErr = append(vErr, &validationError{
				message:  vErrLackOfTabName,
				position: "layout[" + strconv.Itoa(i1) + "]",
			})
		}
		// tabNode.Rows should be set
		if len(t.Rows) == 0 {
			vErr = append(vErr, &validationError{
				message:  vErrLackOfTabRows,
				position: "layout[" + strconv.Itoa(i1) + "]",
			})
		}
		rowHeight := 0
		for i2, r := range t.Rows {
			// rowNode.Height should be "[1-100]%" format
			hp := utils.Percentalize(100, r.Height)
			if hp < 0 || hp > 100 {
				vErr = append(vErr, &validationError{
					message:  vErrRowHeightInvalid,
					position: "layout[" + strconv.Itoa(i1) + "].rows[" + strconv.Itoa(i2) + "].height",
				})
			}
			rowHeight = rowHeight + hp
			// rowNode.Cols should be set
			if len(r.Cols) == 0 {
				vErr = append(vErr, &validationError{
					message:  vErrLackOfRowCols,
					position: "layout[" + strconv.Itoa(i1) + "].rows[" + strconv.Itoa(i2) + "]",
				})
			}
			colWidth := 0
			for i3, cl := range r.Cols {
				// colNode.Width should be [1-12]
				if cl.Width < 0 || cl.Width > 12 {
					vErr = append(vErr, &validationError{
						message:  vErrColWidthInvalid,
						position: "layout[" + strconv.Itoa(i1) + "].rows[" + strconv.Itoa(i2) + "].cols[" + strconv.Itoa(i3) + "].width",
					})
				}
				colWidth = colWidth + cl.Width
				// colNode.Stacks should be set
				if len(cl.Stacks) == 0 {
					vErr = append(vErr, &validationError{
						message:  vErrLackOfColStacks,
						position: "layout[" + strconv.Itoa(i1) + "].rows[" + strconv.Itoa(i2) + "].cols[" + strconv.Itoa(i3) + "]",
					})
				}
				stackHeight := 0
				for i4, s := range cl.Stacks {
					// stackNode should have "id"
					if s.ID == "" {
						vErr = append(vErr, &validationError{
							message:  vErrLackOfStackID,
							position: "layout[" + strconv.Itoa(i1) + "].rows[" + strconv.Itoa(i2) + "].cols[" + strconv.Itoa(i3) + "].stacks[" + strconv.Itoa(i4) + "]",
						})
					}
					// stackNode.ID should be defined on widgets sectoin
					defined := false
					for _, w := range config.Widgets {
						if s.ID == w.ID {
							defined = true
							break
						}
					}
					if !defined {
						vErr = append(vErr, &validationError{
							message:  vErrWidgetDoesNotExit,
							position: "widgets[" + strconv.Itoa(i4) + "].id = " + s.ID + "",
						})
					}
					// stackNode.Height should be "[1-100]%" format
					hp := utils.Percentalize(100, s.Height)
					if hp < 0 || hp > 100 {
						vErr = append(vErr, &validationError{
							message:  vErrWidgetDoesNotExit,
							position: "layout[" + strconv.Itoa(i1) + "].rows[" + strconv.Itoa(i2) + "].cols[" + strconv.Itoa(i3) + "].stacks[" + strconv.Itoa(i4) + "]",
						})
					}
					stackHeight = stackHeight + hp
				}
				// total of the stackNode.Height in a colNode.Stacks should be <= 100%
				if stackHeight > 100 {
					vErr = append(vErr, &validationError{
						message:  vErrInvalidTotalStackHeight,
						position: "layout[" + strconv.Itoa(i1) + "].rows[" + strconv.Itoa(i2) + "].cols[" + strconv.Itoa(i3) + "].stacks[]",
					})
				}
			}
			// total of the colNode.Width in a rowNode.Cols should be <= 12
			if colWidth > 12 {
				vErr = append(vErr, &validationError{
					message:  vErrInvalidTotalColWidth,
					position: "layout[" + strconv.Itoa(i1) + "].rows[" + strconv.Itoa(i2) + "].cols[]",
				})
			}
		}
		// total of the rowNode.Height in a tabNode.Rows should be <= 100%
		if rowHeight > 100 {
			vErr = append(vErr, &validationError{
				message:  vErrInvalidTotalRowHeight,
				position: "layout[" + strconv.Itoa(i1) + "].rows[]",
			})
		}
	}
	return
}
