// parser contains useful functions to parse a scripts file.
package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

/* ParseFile goes through the file named filename and performs all of the
   actions listed in that file.

The file follows the following format:
  Every command is a single character that takes up a line
  Any command that requires arguments must have those arguments in the 2nd line.
  The commands are as follows:
    circle: add a circle to the edge matrix -
	    takes 4 arguments (cx, cy, cz, r)
  	hermite: add a hermite curve to the edge matrix -
      takes 8 arguments (x0, y0, x1, y1, rx0, ry0, rx1, ry1)
	  bezier: add a bezier curve to the edge matrix -
	    takes 8 arguments (x0, y0, x1, y1, x2, y2, x3, y3)
    line: add a line to the edge matrix -
	    takes 6 arguemnts (x0, y0, z0, x1, y1, z1)
	  ident: set the transform matrix to the identity matrix -
	  scale: create a scale matrix, then multiply the transform matrix by the
      scale matrix -
	    takes 3 arguments (sx, sy, sz)
    translate: create a translation matrix, then multiply the transform matrix
      by the translation matrix -
	    takes 3 arguments (tx, ty, tz)
    rotate: create an rotation matrix, then  multiply the transform matrix by
      the rotation matrix -
	    takes 2 arguments (axis, theta) axis should be x y or z
    apply: apply the current transformation matrix to the edge  matrix
	  display: draw the lines of the edge matrix to the screen display  the screen
	  save: draw the lines of the edge matrix to the screen save the screen to a
       file -
	    takes 1 argument (file name)
	  quit: end parsing
*/
func ParseFile(filename string,
	transform [][]float64,
	edges [][]float64,
	screen [][][]int) {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Immediate operations (no arguments)
		if line == "ident" {
			MakeIdentity(transform)
			continue
		} else if line == "display" {
			ClearScreen(screen)
			DrawLines(edges, screen)
			DisplayScreen(screen)
			continue
		} else if line == "clear" {
			edges = make([][]float64, 4)
			continue
		} else if line == "apply" {
			MultiplyMatrices(&transform, &edges)
			continue
		} else if line == "quit" {
			return
		} else if line == "draw" {
			DrawLines(edges, screen)
			continue
		} else if line == "show" {
			DisplayScreen(screen)
			continue
		} else if strings.Contains(line, "color") {
			SetColor(strings.Fields(line)[1])
			continue
		}

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		scanner.Scan()

		// Non-immediate operations (has arguments)
		params := scanner.Text()

		if line == "save" {
			WriteScreenToExtension(screen, params)
		} else if line == "line" {
			AddEdge(edges, FloatParams(params)...)
		} else if line == "circle" {
			AddCircle(edges, FloatParams(params)...)
		} else if line == "sphere" {
			AddSphere(edges, FloatParams(params)...)
		} else if line == "box" {
			AddBox(edges, FloatParams(params)...)
		} else if line == "torus" {
			AddTorus(edges, FloatParams(params)...)
		} else if line == "hermite" || line == "bezier" {
			p := FloatParams(params)
			AddCurve(edges, p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7], 0.001, line)
		} else {
			var stepTransform [][]float64

			if line == "move" {
				stepTransform = MakeTranslationMatrix(FloatParams(params)...)
			} else if line == "scale" {
				stepTransform = MakeDilationMatrix(FloatParams(params)...)
			} else if line == "rotate" {
				args := strings.Fields(params)
				numDegrees, err := strconv.ParseFloat(args[1], 64)
				if err != nil {
					panic(err)
				}

				switch args[0] {
				case "x":
					stepTransform = MakeRotX(numDegrees)
				case "y":
					stepTransform = MakeRotY(numDegrees)
				case "z":
					stepTransform = MakeRotZ(numDegrees)
				}
			}

			if len(transform) == 0 {
				transform = stepTransform
			} else {
				MultiplyMatrices(&stepTransform, &transform)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func FloatParams(text string) (args []float64) {
	args = []float64{}
	for _, v := range strings.Fields(text) {
		floated, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		args = append(args, floated)
	}
	return
}
