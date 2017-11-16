package vector

// RectangleCollection is
type RectangleCollection struct {
	rects []*Rectangle
}

// NewRectangleCollection constructs a ConfigManager
func NewRectangleCollection() (r *RectangleCollection, err error) {
	r = new(RectangleCollection)
	return
}

// Register is
func (r *RectangleCollection) Register(opt *RectangleOptions) {
	rect := NewRectangle(opt)
	r.rects = append(r.rects, rect)
}

// GetRectangle is
func (r *RectangleCollection) GetRectangle(idx int) *Rectangle {
	return r.rects[idx]
}

// CalcDistances is
func (r *RectangleCollection) CalcDistances() {
	type distance struct {
		direction string
		from      *Rectangle
		to        *Rectangle
		toX       int
		toY       int
		toTab     int
		fromIndex int
		toIndex   int
		value     float64
	}
	directions := []string{
		"top",
		"right",
		"bottom",
		"left",
	}
	distances := []*distance{}
	for _, d := range directions {
		for i1, w1 := range r.rects {
			for i2, w2 := range r.rects {
				if i1 == i2 || w1.Tab != w2.Tab {
					continue
				}
				distances = append(distances, &distance{
					direction: d,
					from:      w1,
					to:        w2,
					toX:       w2.Center.x,
					toY:       w2.Center.y,
					toTab:     w2.Tab,
					fromIndex: w1.index,
					toIndex:   w2.index,
					value:     w1.getDistance(w2, d),
				})
			}
		}
	}
	for i, f := range r.rects {
		var nearestTop *distance
		var nearestBottom *distance
		var nearestLeft *distance
		var nearestRight *distance
		for _, d := range distances {
			if f.index == d.toIndex || f.index != d.fromIndex || d.value < 0 {
				continue
			}
			// if the distance is almost same
			// first widget should be primary
			if d.direction == "top" && !f.edge.top && (!f.firstStack || f.rowIndex-1 == d.to.rowIndex) {
				if nearestTop == nil || nearestTop.value > d.value {
					if f.firstStack && f.colIndex == 0 && d.to.colIndex == 0 && d.to.lastStack {
						d.value = d.value / 2
					}
					nearestTop = d
				}
			} else if d.direction == "right" && !f.edge.right && f.rowIndex == d.to.rowIndex && f.colIndex != d.to.colIndex {
				if nearestRight == nil || nearestRight.value > d.value {
					if f.firstStack && d.to.firstStack && f.colIndex+1 == d.to.colIndex {
						d.value = d.value / 2
					}
					nearestRight = d
				}
			} else if d.direction == "bottom" && !f.edge.bottom && (!f.lastStack || f.rowIndex+1 == d.to.rowIndex) {
				if nearestBottom == nil || nearestBottom.value > d.value {
					if f.lastStack && f.colIndex == 0 && d.to.colIndex == 0 && d.to.firstStack {
						d.value = d.value / 2
					}
					nearestBottom = d
				}
			} else if d.direction == "left" && !f.edge.left && f.rowIndex == d.to.rowIndex && f.colIndex != d.to.colIndex {
				if nearestLeft == nil || nearestLeft.value > d.value {
					if f.firstStack && d.to.firstStack && f.colIndex-1 == d.to.colIndex {
						d.value = d.value / 2
					}
					nearestLeft = d
				}
			}
		}
		if nearestTop != nil {
			r.rects[i].TopWidgetIndex = nearestTop.to.index
		} else {
			r.rects[i].TopWidgetIndex = -1
		}
		if nearestBottom != nil {
			r.rects[i].BottomWidgetIndex = nearestBottom.to.index
		} else {
			r.rects[i].BottomWidgetIndex = -1
		}
		if nearestLeft != nil {
			r.rects[i].LeftWidgetIndex = nearestLeft.to.index
		} else {
			r.rects[i].LeftWidgetIndex = -1
		}
		if nearestRight != nil {
			r.rects[i].RightWidgetIndex = nearestRight.to.index
		} else {
			r.rects[i].RightWidgetIndex = -1
		}
	}
	// prevent moving on incorrect widget if the cursor is on the edge widget.
	// this works if the layout is not filled with widgets.
	for i, r1 := range r.rects {
		for _, r2 := range r.rects {
			if r1.index != r2.index && r1.BottomWidgetIndex == r2.index && r1.area.lb.y == r2.area.lb.y {
				r.rects[i].BottomWidgetIndex = -1
			}
			if r1.index != r2.index && r1.RightWidgetIndex == r2.index && r1.area.rb.x == r2.area.rb.x {
				r.rects[i].RightWidgetIndex = -1
			}
		}
	}
}
