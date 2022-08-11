package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type LorenzAttractor struct {
	scale float64

	state0 []float64
	oldPos rl.Vector3

	history     []rl.Vector3
	max_history int
	t_step      float64
	Color       rl.Color

	rho   float64
	sigma float64
	beta  float64
}

func (l *LorenzAttractor) Init(scaleTo float64, timeStep float64) {
	l.state0 = []float64{1.0, 1.0, 1.0}
	l.scale = scaleTo
	l.t_step = timeStep
	l.max_history = 4000
	l.oldPos = rl.Vector3Zero()
	l.history = []rl.Vector3{l.oldPos}

	// default rho, sigma and beta variables
	l.rho = 28.0
	l.sigma = 10.0
	l.beta = 8.0 / 3.0

}

// set different rho, sigma and beta values
// defaults are rho = 28, sigma = 10, beta = 8/3
func (l *LorenzAttractor) SetVariables(rho float64, sigma float64, beta float64) {
	l.rho = rho
	l.sigma = sigma
	l.beta = beta
}

func (l *LorenzAttractor) UpdateAndDraw(camera *rl.Camera3D) {
	x, y, z := l.f(l.state0, 0.0)
	l.state0[0] += x * l.t_step
	l.state0[1] += y * l.t_step
	l.state0[2] += z * l.t_step

	newPos := rl.Vector3{X: float32(x * l.scale), Y: float32(y * l.scale), Z: float32(z * l.scale)}
	his_len := len(l.history)

	// pop first element
	if his_len < l.max_history {
		l.history = append(l.history, newPos)
		his_len++
	} else {
		l.history = l.history[1:]
		his_len--
	}

	rl.BeginMode3D(*camera)

	if his_len > 10 {
		// draw the attractor
		for i := 5; i < his_len; i += 4 {
			points := SmoothLine(l.history[i-5:i], 3)
			for j := 1; j < len(points); j++ {
				rl.DrawLine3D(points[j-1], points[j], l.Color)
			}

		}
	}

	rl.DrawLine3D(l.oldPos, newPos, rl.Red)
	l.oldPos = newPos

	rl.EndMode3D()
}

func (l *LorenzAttractor) f(state []float64, t float64) (float64, float64, float64) {
	x := state[0]
	y := state[1]
	z := state[2]
	return l.sigma * (y - x), x*(l.rho-z) - y, x*y - l.beta*z
}

func (l *LorenzAttractor) Reset() {
	l.Init(l.scale, l.t_step)
}

// create a smooth bezier curve between the points
func SmoothLine(pointsToCurve []rl.Vector3, smoothness float32) []rl.Vector3 {
	// https://answers.unity.com/questions/392606/line-drawing-how-can-i-interpolate-between-points.html
	if smoothness < 1.0 {
		smoothness = 1.0
	}
	pointsLength := len(pointsToCurve)
	curvedLength := (pointsLength * int(smoothness)) - 1
	curvedPoints := []rl.Vector3{}

	t := 0.0
	var points []rl.Vector3
	for pointInTime := 0; pointInTime < curvedLength+1; pointInTime++ {
		t = InvLerp(0, curvedLength, pointInTime)

		points = make([]rl.Vector3, len(pointsToCurve))
		copy(points, pointsToCurve)

		for j := pointsLength - 1; j > 0; j-- {
			for i := 0; i < j; i++ {
				points[i] = rl.Vector3Add(rl.Vector3Multiply(points[i], float32(1.0-t)), rl.Vector3Multiply(points[i+1], float32(t)))
			}
		}
		curvedPoints = append(curvedPoints, points[0])
		// curvedPoints[pointInTime] = points[0]
	}

	return curvedPoints
}

// https://www.gamedev.net/articles/programming/general-and-gameplay-programming/inverse-lerp-a-super-useful-yet-often-overlooked-function-r5230/
func InvLerp(a int, b int, v int) float64 {
	return float64(v-a) / float64(b-a)
}
