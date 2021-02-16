package main

import (
	_ "embed"
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"github.com/fzipp/geom"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	resW = 800
	resH = 600
)

type Rect struct {
	geom.Vec2
	geom.Size
}

var (
	//go:embed shader/vertex.glsl
	vertexShader string
	//go:embed shader/fragment.glsl
	fragmentShader string
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	win, err := initGraphics(resW, resH)
	defer glfw.Terminate()

	check(err)

	vertexShader, fragmentShader, err := loadShaders(vertexShader, fragmentShader)
	check(err)

	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	var vertices quadArray
	quadCoords(Rect{geom.V2(100, 100), geom.Size{W: 600, H: 400}}, &vertices)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices), gl.STATIC_DRAW)

	program, err := loadProgram(vertexShader, fragmentShader)
	check(err)
	gl.UseProgram(program)

	// Init and load PMV matrices
	var pmatrix, mvmatrix geom.Mat4
	pmatrix.Ortho(0, resW, resH, 0, -1.0, 1.0)
	loadMatrix(&pmatrix, program, "pmatrix")
	mvmatrix.ID()
	loadMatrix(&mvmatrix, program, "mvmatrix")

	attr := uint32(gl.GetAttribLocation(program, gl.Str("vertex\x00")))

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.EnableVertexAttribArray(attr)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.VertexAttribPointer(attr, dimensions, gl.FLOAT, false, 0, nil)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, quadVertexCount)
		gl.DisableVertexAttribArray(attr)

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGraphics(width, height int) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("could not initialize GLFW: %w", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(width, height, "OpenGL Demo", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create windoe: %w", err)
	}

	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("could not initialize OpenGL bindings: %w", err)
	}

	return win, nil
}

const (
	quadVertexCount = 4
	dimensions      = 2
)

type quadArray [quadVertexCount * dimensions]float32

func quadCoords(r Rect, coords *quadArray) {
	coords[0] = r.X
	coords[1] = r.Y
	coords[2] = r.X + r.W
	coords[3] = r.Y
	coords[4] = r.X
	coords[5] = r.Y + r.H
	coords[6] = coords[2]
	coords[7] = coords[5]
}

func check(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
