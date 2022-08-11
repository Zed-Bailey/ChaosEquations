package main

import (
	"fmt"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	winWidth  = 800 * 2
	winHeight = 450 * 2
)

var (
	number_of_attractors = 1
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		n, err := strconv.ParseInt(args[0], 0, 64)
		if err != nil {
			fmt.Printf("-- Failed to parse the cli argument: %s. Make sure it is a number", args[0])
			return
		}
		number_of_attractors = int(n)

	}
	rl.SetTargetFPS(60)
	rl.InitWindow(winWidth, winHeight, "Lorenz Attractor")

	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(0.0, 0.0, 15.0)
	camera.Target = rl.NewVector3(1.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	rl.SetCameraMode(camera, rl.CameraOrbital)
	// time
	t := 0.0
	// the time step
	t_step := 0.01

	attractor := []LorenzAttractor{}

	colours := []rl.Color{
		rl.Red,
		rl.Purple,
		rl.Green,
		rl.Blue,
		rl.Gold,
		rl.Lime,
	}

	fmt.Printf("-- Creating %d attractors\n", number_of_attractors)
	// scale the lorenz attractor down
	scale := 0.2

	if number_of_attractors > 1 {
		for i := 0; i < number_of_attractors; i++ {
			attractor = append(attractor, LorenzAttractor{})
			attractor[i].Init(scale, t_step)
			attractor[i].Color = colours[rl.GetRandomValue(0, int32(len(colours)-1))]
			attractor[i].SetVariables(float64(rl.GetRandomValue(25, 30)), float64(rl.GetRandomValue(8, 12)), float64(rl.GetRandomValue(6, 8))/3)
		}
		// set the first attractyor with the correct values
		attractor[0].SetVariables(28, 10, 8.0/3.0)
	} else {
		attractor = append(attractor, LorenzAttractor{})
		attractor[0].Init(scale, t_step)
		attractor[0].Color = rl.Gold
	}

	fmt.Println("-- Zoom with mouse wheel")

	for !rl.WindowShouldClose() {

		// updateCamera(&camera)
		rl.UpdateCamera(&camera)
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		for i := 0; i < number_of_attractors; i++ {
			attractor[i].UpdateAndDraw(&camera)
		}

		rl.DrawText(fmt.Sprintf("t = %f", t), 15, 15, 15, rl.White)
		rl.EndDrawing()

		t += t_step
	}

	rl.CloseWindow()
}

var (
	camDist   int32   = 100
	rotAngle  float32 = 45.0
	tiltAngle float32 = 45.0

	rotSpeed  float32 = 0.25
	moveSpeed         = 10

	cursorPosition = rl.GetMousePosition()
)

// https://gist.github.com/JeffM2501/000787070aef421a00c02ae4cf799ea1
func updateCamera(camera *rl.Camera3D) {
	// Update Camera
	if rl.IsMouseButtonDown(1) {
		newPos := rl.GetMousePosition()
		rotAngle += (newPos.X - cursorPosition.X) * rotSpeed
		tiltAngle += (newPos.Y - cursorPosition.Y) * rotSpeed

		// clamp the tilt so we don't get gymbal lock
		if tiltAngle > 89.0 {
			tiltAngle = 89
		}
		if tiltAngle < 1.0 {
			tiltAngle = 1
		}

	}

	cursorPosition = rl.GetMousePosition()

	moveVec := rl.Vector3Zero()

	if rl.IsKeyDown(rl.KeyW) {
		moveVec.Z = -float32(moveSpeed) * rl.GetFrameTime()
	}
	if rl.IsKeyDown(rl.KeyS) {
		moveVec.Z = float32(moveSpeed) * rl.GetFrameTime()
	}
	if rl.IsKeyDown(rl.KeyA) {
		moveVec.X = -float32(moveSpeed) * rl.GetFrameTime()
	}
	if rl.IsKeyDown(rl.KeyD) {
		moveVec.X = float32(moveSpeed) * rl.GetFrameTime()
	}
	// update zoom
	camDist += rl.GetMouseWheelMove()
	if camDist < 1 {
		camDist = 1
	}

	camPos := rl.Vector3{X: 0, Y: 0, Z: float32(camDist)}

	tiltMat := rl.MatrixRotateX(tiltAngle * rl.GetFrameTime()) // a matrix for the tilt rotation
	rotMat := rl.MatrixRotateY(rotAngle * rl.GetFrameTime())   // a matrix for the plane rotation
	mat := rl.MatrixMultiply(tiltMat, rotMat)                  // the combined transformation matrix for the camera position

	camPos = rl.Vector3Transform(camPos, mat)      // transform the camera position into a vector in world space
	moveVec = rl.Vector3Transform(moveVec, rotMat) // transform the movement vector into world space, but ignore the tilt so it is in plane

	camera.Target = rl.Vector3Add(camera.Target, moveVec) // move the target to the moved position

	camera.Position = rl.Vector3Add(camera.Target, camPos) // offset the camera position by the vector from the target position

}
