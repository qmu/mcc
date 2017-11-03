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
				var val float64
				if d == "top" {
					val = w1.toTop(w2)
				} else if d == "right" {
					val = w1.toRight(w2)
				} else if d == "bottom" {
					val = w1.toBottom(w2)
				} else if d == "left" {
					val = w1.toLeft(w2)
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
					value:     val,
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
			if d.direction == "top" && !f.edge.top {
				if nearestTop == nil || nearestTop.value > d.value {
					nearestTop = d
				}
			} else if d.direction == "right" && !f.edge.right {
				if nearestRight == nil || nearestRight.value > d.value {
					nearestRight = d
				}
			} else if d.direction == "bottom" && !f.edge.bottom {
				if nearestBottom == nil || nearestBottom.value > d.value {
					nearestBottom = d
				}
			} else if d.direction == "left" && !f.edge.left {
				if nearestLeft == nil || nearestLeft.value > d.value {
					nearestLeft = d
				}
			}
		}
		if nearestTop != nil {
			r.rects[i].TopWidgetIndex = nearestTop.to.index
		}
		if nearestBottom != nil {
			r.rects[i].BottomWidgetIndex = nearestBottom.to.index
		}
		if nearestLeft != nil {
			r.rects[i].LeftWidgetIndex = nearestLeft.to.index
		}
		if nearestRight != nil {
			r.rects[i].RightWidgetIndex = nearestRight.to.index
		}
	}
}
