package dashboard

// WidgetOptions defines options for each type of widgets
type WidgetOptions struct {
	extendedWidget *ExtendedWidget
	envs           []map[string]string
	execPath       string
	timezone       string
}

// GetContent is
func (w *WidgetOptions) GetContent() interface{} {
	return w.extendedWidget.widget.Content
}

// GetHeight is
func (w *WidgetOptions) GetHeight() int {
	return w.extendedWidget.height
}

// GetWidth is
func (w *WidgetOptions) GetWidth() int {
	return w.extendedWidget.width
}

// GetTitle is
func (w *WidgetOptions) GetTitle() string {
	return w.extendedWidget.title
}
