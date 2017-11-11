package vector

import "math"

// Rectangle is
type Rectangle struct {
	index      int
	rowIndex   int
	colIndex   int
	lastCol    bool
	firstStack bool
	lastStack  bool
	edge       Edge
	area       Area

	Height            int
	Width             int
	Tab               int
	Center            Point
	BottomWidgetIndex int
	TopWidgetIndex    int
	LeftWidgetIndex   int
	RightWidgetIndex  int
}

// RectangleOptions is
type RectangleOptions struct {
	Index      int
	RowIndex   int
	ColIndex   int
	LastCol    bool
	FirstStack bool
	LastStack  bool
	WindowW    int
	WindowH    int
	TabIndex   int
	CenterX    int
	CenterY    int
	Height     int
	Width      int
}

// NewRectangle is
func NewRectangle(opt *RectangleOptions) (r *Rectangle) {
	windowW := opt.WindowW
	windowH := opt.WindowH
	r = new(Rectangle)
	center := Point{
		x: opt.CenterX,
		y: opt.CenterY,
	}
	r = &Rectangle{
		index:      opt.Index,
		rowIndex:   opt.RowIndex,
		colIndex:   opt.ColIndex,
		lastCol:    opt.LastCol,
		firstStack: opt.FirstStack,
		lastStack:  opt.LastStack,
		Tab:        opt.TabIndex,
		Height:     opt.Height,
		Width:      opt.Width,
		Center:     center,
		edge: Edge{
			top:    center.y-opt.Height/2-1 <= 0,
			right:  center.x+opt.Width/2+1 >= windowW,
			bottom: center.y+opt.Height/2+1 >= windowH,
			left:   center.x-opt.Width/2-1 <= 0,
		},
		area: Area{
			lt: Point{
				x: center.x - opt.Width/2,
				y: center.y - opt.Height/2,
			},
			rt: Point{
				x: center.x + opt.Width/2,
				y: center.y - opt.Height/2,
			},
			lb: Point{
				x: center.x - opt.Width/2,
				y: center.y + opt.Height/2,
			},
			rb: Point{
				x: center.x + opt.Width/2,
				y: center.y + opt.Height/2,
			},
		},
	}

	return
}

// Edge is
type Edge struct {
	top    bool
	right  bool
	bottom bool
	left   bool
}

// Area is
type Area struct {
	lt Point
	rt Point
	lb Point
	rb Point
}

// Point is
type Point struct {
	x int
	y int
}

func (r *Rectangle) getDistance(from *Rectangle, d string) float64 {
	var val float64
	if d == "top" {
		val = r.toTop(from)
	} else if d == "right" {
		val = r.toRight(from)
	} else if d == "bottom" {
		val = r.toBottom(from)
	} else if d == "left" {
		val = r.toLeft(from)
	}
	return val
}
func (r *Rectangle) toTop(wTop *Rectangle) float64 {
	wBottom := r
	if wBottom.Center.y < wTop.Center.y {
		return -1
	}
	v := ltTolb(wBottom.area, wTop.area)
	v = v + rtTorb(wBottom.area, wTop.area)
	v = v + vectorDistance(wBottom.Center, wTop.Center)
	return v
}
func (r *Rectangle) toBottom(wBottom *Rectangle) float64 {
	wTop := r
	if wTop.Center.y > wBottom.Center.y {
		return -1
	}
	v := lbTolt(wTop.area, wBottom.area)
	v = v + rbTort(wTop.area, wBottom.area)
	v = v + vectorDistance(wTop.Center, wBottom.Center)
	return v
}
func (r *Rectangle) toRight(wRight *Rectangle) float64 {
	wLeft := r
	if wLeft.Center.x > wRight.Center.x {
		return -1
	}
	v := rtTolt(wLeft.area, wRight.area)
	v = v + rbTolb(wLeft.area, wRight.area)
	v = v + vectorDistance(wLeft.Center, wRight.Center)
	return v
}
func (r *Rectangle) toLeft(wLeft *Rectangle) float64 {
	wRight := r
	if wRight.Center.x < wLeft.Center.x {
		return -1
	}
	v := lbTorb(wRight.area, wLeft.area)
	v = v + ltTort(wRight.area, wLeft.area)
	v = v + vectorDistance(wRight.Center, wLeft.Center)
	return v
}

func ltTolb(areaBottom Area, areaTop Area) float64 {
	return vectorDistance(areaBottom.lt, areaTop.lb)
}
func rtTorb(areaBottom Area, areaTop Area) float64 {
	return vectorDistance(areaBottom.rt, areaTop.rb)
}
func lbTolt(areaTop Area, areaBottom Area) float64 {
	return vectorDistance(areaTop.lb, areaBottom.lt)
}
func rbTort(areaTop Area, areaBottom Area) float64 {
	return vectorDistance(areaTop.rb, areaBottom.rt)
}
func rtTolt(areaLeft Area, areaRight Area) float64 {
	return vectorDistance(areaLeft.rt, areaRight.lt)
}
func rbTolb(areaLeft Area, areaRight Area) float64 {
	return vectorDistance(areaLeft.rb, areaRight.lb)
}
func lbTorb(areaRight Area, areaLeft Area) float64 {
	return vectorDistance(areaRight.lb, areaLeft.rb)
}
func ltTort(areaRight Area, areaLeft Area) float64 {
	return vectorDistance(areaRight.lt, areaLeft.rt)
}
func vectorDistance(fromPoint Point, toPoint Point) (distance float64) {
	x1 := fromPoint.x
	y1 := fromPoint.y
	x2 := toPoint.x
	y2 := toPoint.y
	distance = math.Pow(float64((x2-x1)*(x2-x1)+(y2-y1)*(y2-y1)), 0.5)
	return
}
func abs(value int) int {
	return int(math.Abs(float64(value)))
}
