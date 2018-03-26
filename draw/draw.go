// Package draw contains useful functions to manipulate the edge matrix and draw
// it onto a screen.
package draw

import (
	"github.com/jkao1/yet-another-3d-thing/display"
	"github.com/jkao1/yet-another-3d-thing/matrix"

	"math"
)

var DefaultDrawColor []int = []int{0, 0, 0}

// DrawLines draws an edge matrix onto a screen.
func DrawLines(edges [][]float64, screen [][][]int) {
	for i := 0; i < len(edges[0])-1; i += 2 {
		point := matrix.ExtractColumn(edges, i)
		nextPoint := matrix.ExtractColumn(edges, i+1)
		x0, y0 := point[0], point[1]
		x1, y1 := nextPoint[0], nextPoint[1]
		DrawLine(screen, x0, y0, x1, y1)
	}
}

// AddPoint adds a point to an edge matrix.
func AddPoint(m [][]float64, x, y, z float64) {
	m[0] = append(m[0], x)
	m[1] = append(m[1], y)
	m[2] = append(m[2], z)
	m[3] = append(m[3], 1)
}

// AddEdge adds an edge (two points) to an edge matrix.
func AddEdge(m [][]float64, params ...float64) {
	x0, y0, z0 := params[0], params[1], params[2]
	x1, y1, z1 := params[3], params[4], params[5]
	AddPoint(m, x0, y0, z0)
	AddPoint(m, x1, y1, z1)
}

// AddCircle adds a circle of center (cx, cy, cz) and radius r to an edge
// matrix.
func AddCircle(m [][]float64, params ...float64) {
	cx, cy, _, r := params[0], params[1], params[2], params[3]
	for t := 0.0; t <= 1.0; t += 0.001 {
		x := r*math.Cos(2*math.Pi*t) + cx
		y := r*math.Sin(2*math.Pi*t) + cy
		AddPoint(m, x, y, 0)
	}
}

// AddCurve adds the curve bounded by the 4 points passed as parameters
// to an edge matrix.
func AddCurve(m [][]float64, x0, y0, x1, y1, x2, y2, x3, y3, step float64, curveType string) {
	xCoefs := generateCurveCoefs(x0, x1, x2, x3, curveType)
	yCoefs := generateCurveCoefs(y0, y1, y2, y3, curveType)

	for t := 0.0; t <= 1.0; t += step {
		x := CubicEval(t, xCoefs)
		y := CubicEval(t, yCoefs)

		AddPoint(m, x, y, 0)
	}
}

func generateCurveCoefs(p0, p1, p2, p3 float64, curveType string) [][]float64 {
	m := make([][]float64, 4)
	var coefGenerator [][]float64
	if curveType == "hermite" {
		coefGenerator = matrix.MakeHermite()
	} else if curveType == "bezier" {
		coefGenerator = matrix.MakeBezier()
	}
	m[0] = []float64{p0}
	m[1] = []float64{p1}
	m[2] = []float64{p2}
	m[3] = []float64{p3}
	matrix.MultiplyMatrices(&coefGenerator, &m)
	return m
}

// AddBox adds the points for a rectagular prism whose upper-left corner is
// (x, y, z) with width, height and depth dimensions.
func AddBox(m [][]float64, a ...float64) {
	x, y, z, width, height, depth := a[0], a[1], a[2], a[3], a[4], a[5]
	AddEdge(m, x, y, z, x+width, y, z)
	AddEdge(m, x, y, z, x, y-height, z)
	AddEdge(m, x, y, z, x, y, z-depth)

	AddEdge(m, x, y-height, z, x, y-height, z-depth)
	AddEdge(m, x, y-height, z, x+width, y-height, z)
	AddEdge(m, x, y-height, z-depth, x+width, y-height, z-depth)
	AddEdge(m, x+width, y-height, z, x+width, y-height, z-depth)

	AddEdge(m, x, y, z-depth, x+width, y, z-depth)
	AddEdge(m, x, y, z-depth, x, y-height, z-depth)

	AddEdge(m, x+width, y, z, x+width, y-height, z)
	AddEdge(m, x+width, y, z, x+width, y, z-depth)
	AddEdge(m, x+width, y, z-depth, x+width, y-height, z-depth)
}

// AddSphere adds all the points for a sphere with center (cx, cy, cz) and
// radius r.
func AddSphere(m [][]float64, a ...float64) {
	cx, cy, cz, r := a[0], a[1], a[2], a[3]
	for _, p := range GenerateSphere(cx, cy, cz, r) {
		AddEdge(m, p[0], p[1], p[2], p[0]+1, p[1]+1, p[2]+1)
	}
}

// GenerateSphere generates all the points along the surface of a sphere with
// center (cx, cy, cz) and radius r. It returns a matrix of the points.
func GenerateSphere(cx, cy, cz, r float64) [][]float64 {
	points := make([][]float64, 0)
	for i := 0.0; i <= 1.0; i += 0.01 {
		fi := 2 * math.Pi * i
		for j := 0.0; j <= 0.5; j += 0.01 {
			theta := 2 * math.Pi * j
			x := r*math.Cos(theta) + cx
			y := r*math.Sin(theta)*math.Cos(fi) + cy
			z := r*math.Sin(theta)*math.Sin(fi) + cz
			points = append(points, []float64{x, y, z})
		}
	}
	return points
}

// AddTorus adds all the points required to make a torus with center
// (cx, cy, cz) and radii r1 and r2.
func AddTorus(m [][]float64, a ...float64) {
	cx, cy, cz, r1, r2 := a[0], a[1], a[2], a[3], a[4]
	for _, p := range GenerateTorus(cx, cy, cz, r1, r2) {
		AddEdge(m, p[0], p[1], p[2], p[0]+1, p[1]+1, p[2]+1)
	}
}

// GenerateTorus  generates all the points along the surface of a torus with
// center (cx, cy, cz) and radii r1 and r2.
func GenerateTorus(cx, cy, cz, r2, r1 float64) [][]float64 {
	points := make([][]float64, 0)
	for i := 0.0; i <= 1.0; i += 0.01 {
		fi := 2 * math.Pi * i
		for j := 0.0; j <= 1.0; j += 0.01 {
			theta := 2 * math.Pi * j
			x := math.Cos(fi)*(r2*math.Cos(theta)+r1) + cx
			y := r2*math.Sin(theta) + cy
			z := -1*math.Sin(fi)*(r2*math.Cos(theta)+r1) + cz
			points = append(points, []float64{x, y, z})
		}
	}
	return points
}

// CubicEval evaluates a cubic function with variable x and coefficients.
func CubicEval(x float64, coefs [][]float64) (y float64) {
	for i := 3.0; i >= 0.0; i-- {
		y += coefs[int64(3-i)][0] * math.Pow(x, i)
	}
	return
}

// DrawLine draws a line from (x0, y0) to (x1, y1) onto a screen.
func DrawLine(screen [][][]int, x0, y0, x1, y1 float64) {
	if x1 < x0 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	A := y1 - y0
	B := x0 - x1
	x := x0
	y := y0

	if B == 0 { // vertical line
		if y1 < y0 {
			y0, y1 = y1, y0
		}

		y = y0
		for y <= y1 {
			plot(screen, x, y)
			y++
		}

		return
	}

	slope := A / (-B)
	var d float64

	if slope >= 0 && slope <= 1 { // octant 1
		d = 2*A + B
		for x <= x1 && y <= y1 {
			plot(screen, x, y)
			if d > 0 {
				y++
				d += 2 * B
			}
			x++
			d += 2 * A
		}
	}

	if slope > 1 { // octant 2
		d = A + 2*B
		for x <= x1 && y <= y1 {
			plot(screen, x, y)
			if d < 0 {
				x++
				d += 2 * A
			}
			y++
			d += 2 * B
		}
	}

	if slope < 0 && slope >= -1 { // octant 8
		d = 2*A - B
		for x <= x1 && y >= y1 {
			plot(screen, x, y)
			if d < 0 {
				y--
				d -= 2 * B
			}
			x++
			d += 2 * A
		}
	}

	if slope < -1 { // octant 7
		d = A - 2*B
		for x <= x1 && y >= y1 {
			plot(screen, x, y)
			if d > 0 {
				x++
				d += 2 * A
			}
			y--
			d -= 2 * B
		}
	}
}

// SetColor sets the color to draw with.
func SetColor(color string) {
	if color == "yellow" {
		DefaultDrawColor = []int{255, 255, 0}
	} else if color == "white" {
		DefaultDrawColor = []int{255, 255, 255}
	} else if color == "black" {
		DefaultDrawColor = make([]int, 3)
	}
}

// plot draws a point (x, y) onto a screen with the default draw color.
func plot(screen [][][]int, x, y float64) {
	newX, newY := float64ToInt(x), display.YRES-float64ToInt(y)-1
	if newX >= 0 && newX < display.XRES && newY >= 0 && newY < display.YRES {
		screen[newY][newX] = DefaultDrawColor[:]
	}
}

// DrawLineFromParams gets arguments from a params slice.
func DrawLineFromParams(screen [][][]int, params ...float64) {
	if len(params) >= 4 {
		DrawLine(screen, params[0], params[1], params[2], params[3])
	}
}

// float64ToInt rounds a float64 without truncating it. It returns an int.
func float64ToInt(f float64) int {
	if f-float64(int(f)) < 0.5 {
		return int(f)
	}
	return int(f + 1)
}
