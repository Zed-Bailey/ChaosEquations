package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//Translated from c++ to go
// https://github.com/HackerPoet/Chaos-Equations/blob/master/Main.cpp

type Vertex struct {
	pos rl.Vector2
	col rl.Color
}

const (
	window_w = 800
	window_h = 450

	speedMult       = 2.5
	iters           = 800
	steps_per_frame = 500
	delta_per_step  = 1e-5
	delta_minimum   = 1e-7
)

var (
	t_start = 0.0
	t_end   = 3.0

	rolling_delta = delta_per_step
	plot_scale    = 1.75
	plot_x        = 0.0
	plot_y        = 0.0
)

func main() {

	rl.InitWindow(window_w, window_h, "Chaos Equations")

	camera := rl.Camera2D{}
	camera.Target = rl.NewVector2(0, 0)
	camera.Offset = rl.NewVector2(0, 0)
	camera.Rotation = 0.0
	camera.Zoom = 1.0

	rl.SetTargetFPS(60)
	history := make([]rl.Vector2, iters)

	numVertexs := iters * steps_per_frame
	vertex_array := make([]Vertex, numVertexs)

	// generate random colours
	for i := 0; i < numVertexs; i++ {
		vertex_array[i].col = rl.NewColor(uint8(rl.GetRandomValue(0, 255)), uint8(rl.GetRandomValue(0, 255)), uint8(rl.GetRandomValue(0, 255)), uint8(rl.GetRandomValue(0, 255)))
	}

	time := t_start
	fadeSpeed := 10

	for !rl.WindowShouldClose() {
		if time > t_end {
			time = t_start
		}

		const steps = steps_per_frame
		const delta = delta_per_step * speedMult
		rolling_delta = rolling_delta*0.99 + delta*0.01
		// rolling_delta = 0.00005*0.99 + delta*0.01

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginBlendMode(rl.BlendSubtractColors)

		rl.DrawRectangle(0, 0, window_w, window_h, rl.NewColor(uint8(fadeSpeed), uint8(fadeSpeed), uint8(fadeSpeed), 0))
		for step := 0; step < steps; step++ {
			isOffScreen := true
			x := time
			y := time

			for iter := 0; iter < iters; iter++ {
				nx, ny := newXY(x, y, time)
				x, y = nx, ny
				screenpt := ToScreen(rl.Vector2{X: float32(x), Y: float32(y)})

				vertex_array[step*iters+iter].pos = screenpt
				// Check if dynamic delta should be adjusted
				if screenpt.X > 0.0 && screenpt.Y > 0.0 && screenpt.X < window_w && screenpt.Y < window_h {

					dx := float64(history[iter].X) - x
					dy := float64(history[iter].Y) - y
					dist := 500.0 * math.Sqrt(dx*dx+dy*dy)

					rolling_delta = math.Min(rolling_delta, math.Max(delta/(dist+1e-5), delta_minimum*speedMult))
					isOffScreen = false
				}
				history[iter].X = float32(x)
				history[iter].Y = float32(y)
			}

			if isOffScreen {
				time += 0.01
			} else {
				time += rolling_delta
			}
		}

		rl.BeginMode2D(camera)

		// draw vertex's
		for i := 1; i < numVertexs; i++ {
			rl.DrawPixelV(vertex_array[i].pos, vertex_array[i].col)
		}

		rl.EndMode2D()
		rl.EndBlendMode()

		rl.DrawText(fmt.Sprintf("time t = %f\nrolling delta = %f", time, rolling_delta), 20, 20, 15, rl.White)
		rl.EndDrawing()
		// fmt.Printf("t=%f, rd=%f\n", time, rolling_delta)
	}

	rl.CloseWindow()
}

func ToScreen(pos rl.Vector2) rl.Vector2 {
	s := plot_scale * float64(window_h/2)
	nx := float64(window_w)*0.5 + (float64(pos.X)-plot_x)*s
	ny := float64(window_h)*0.5 + (float64(pos.Y)-plot_y)*s

	return rl.Vector2{X: float32(nx), Y: float32(ny)}
}

func newXY(x, y, t float64) (float64, float64) {
	// newX := -math.Pow(y, 2) - math.Pow(t, 2) + (t * x)
	// newY := (y * t) + (x * y)
	newX := -math.Pow(x, 2) + (x * t) + y
	newY := math.Pow(x, 2) - math.Pow(y, 2) - math.Pow(t, 2) - (x * y) + (y * t) - x + y
	return newX, newY
}

func CenterPlot(history []rl.Vector2) {
	min_x := -math.MaxFloat32
	min_y := -math.MaxFloat32
	max_x := math.MaxFloat32
	max_y := math.MaxFloat32

	for i := 0; i < len(history); i++ {
		min_x = math.Min(min_x, float64(history[i].X))
		max_x = math.Max(max_x, float64(history[i].X))
		min_y = math.Min(min_y, float64(history[i].Y))
		max_y = math.Max(max_y, float64(history[i].Y))
	}

	max_x = math.Min(max_x, 4.0)
	max_y = math.Min(max_y, 4.0)
	min_x = math.Max(min_x, -4.0)
	min_y = math.Max(min_y, -4.0)
	plot_x = (max_x + min_x) * 0.5
	plot_y = (max_y + min_y) * 0.5
	plot_scale = 1.0 / math.Max(math.Max(max_x-min_x, max_y-min_y)*0.6, 0.1)

	fmt.Printf("Centered Plot\nx: %f\ny: %f\nscale: %f\n", plot_x, plot_y, plot_scale)
}
