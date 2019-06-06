// +build extra full

package base

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// declared in func_opengl.go
var freeFns map[string]*func() = make(map[string]*func(), 0)
var cSources map[string]**uint8 = make(map[string]**uint8, 0)
var gifs map[string]*gif.GIF = make(map[string]*gif.GIF, 0)

// gogl
func op_gl_Init() {
	gl.Init()
}

func op_gl_Strs(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	dsSource := ReadStr(fp, inp1)
	fnName := ReadStr(fp, inp2)

	csources, free := gl.Strs(dsSource + string('\000'))

	freeFns[fnName] = &free
	cSources[fnName] = csources
}

func op_gl_Free(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fnName := ReadStr(fp, inp1)

	(*freeFns[fnName])()
	delete(freeFns, fnName)
	delete(cSources, fnName)
}

func op_gl_NewTexture(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	file := ReadStr(fp, inp1)

	imgFile, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("texture %q not found on disk: %v\n", file, err))
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	outB1 := encoder.SerializeAtomic(int32(texture))
	WriteMemory(out1Offset, outB1)
}

func op_gl_NewGIF(expr *CXExpression, fp int) {
	path := ReadStr(fp, expr.Inputs[0])

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		panic(fmt.Sprintf("file not found %q, %v", path, err))
	}

	reader := bufio.NewReader(file)
	gif, err := gif.DecodeAll(reader)
	if err != nil {
		panic(fmt.Sprintf("failed to decode file %q, %v", path, err))
	}

	gifs[path] = gif

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(int32(len(gif.Image))))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromI32(int32(gif.LoopCount)))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[2]), FromI32(int32(gif.Config.Width)))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[3]), FromI32(int32(gif.Config.Height)))
}

func op_gl_FreeGIF(expr *CXExpression, fp int) {
	gifs[ReadStr(fp, expr.Inputs[0])] = nil
}

func op_gl_GIFFrameToTexture(expr *CXExpression, fp int) {
	path := ReadStr(fp, expr.Inputs[0])
	frame := ReadI32(fp, expr.Inputs[1])
	texture := ReadI32(fp, expr.Inputs[2])

	gif := gifs[path]
	img := gif.Image[frame]
	delay := int32(gif.Delay[frame])
	disposal := int32(gif.Disposal[frame])

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	gl.BindTexture(gl.TEXTURE_2D, uint32(texture))
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(delay))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromI32(disposal))
}

// gl_0_0
func op_gl_MatrixMode(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.MatrixMode(uint32(ReadI32(fp, inp1)))
}

func op_gl_Rotatef(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.Rotatef(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3), ReadF32(fp, inp4))
}

func op_gl_Translatef(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Translatef(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func op_gl_LoadIdentity() {
	gl.LoadIdentity()
}

func op_gl_PushMatrix() {
	gl.PushMatrix()
}

func op_gl_PopMatrix() {
	gl.PopMatrix()
}

func op_gl_EnableClientState(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.EnableClientState(uint32(ReadI32(fp, inp1)))
}

func op_gl_Color3f(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Color3f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func op_gl_Color4f(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.Color4f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3), ReadF32(fp, inp4))
}

func op_gl_Begin(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.Begin(uint32(ReadI32(fp, inp1)))
}

func op_gl_End() {
	gl.End()
}

func op_gl_Normal3f(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Normal3f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func op_gl_Vertex2f(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.Vertex2f(ReadF32(fp, inp1), ReadF32(fp, inp2))
}

func op_gl_Vertex3f(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Vertex3f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func op_gl_Lightfv(expr *CXExpression, fp int) {
	// pointers
	panic("gl.Lightfv")
}

func op_gl_Frustum(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]
	gl.Frustum(ReadF64(fp, inp1), ReadF64(fp, inp2), ReadF64(fp, inp3), ReadF64(fp, inp4), ReadF64(fp, inp5), ReadF64(fp, inp6))
}

func op_gl_TexEnvi(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.TexEnvi(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)), ReadI32(fp, inp3))
}

func op_gl_Ortho(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]
	gl.Ortho(ReadF64(fp, inp1), ReadF64(fp, inp2), ReadF64(fp, inp3), ReadF64(fp, inp4), ReadF64(fp, inp5), ReadF64(fp, inp6))
}

func op_gl_Scalef(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Scalef(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func op_gl_TexCoord2d(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.TexCoord2d(ReadF64(fp, inp1), ReadF64(fp, inp2))
}

func op_gl_TexCoord2f(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.TexCoord2f(ReadF32(fp, inp1), ReadF32(fp, inp2))
}

// gl_1_0
func op_gl_CullFace(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.CullFace(uint32(ReadI32(fp, inp1)))
}

func op_gl_Hint(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.Hint(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_Scissor(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.Scissor(ReadI32(fp, inp1), ReadI32(fp, inp2), ReadI32(fp, inp3), ReadI32(fp, inp4))
}

func op_gl_TexParameteri(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.TexParameteri(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)), ReadI32(fp, inp3))
}

func op_gl_TexImage2D(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6, inp7, inp8 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5], expr.Inputs[6], expr.Inputs[7]
	gl.TexImage2D(uint32(ReadI32(fp, inp1)), ReadI32(fp, inp2), ReadI32(fp, inp3), ReadI32(fp, inp4), ReadI32(fp, inp5), ReadI32(fp, inp6), uint32(ReadI32(fp, inp7)), uint32(ReadI32(fp, inp8)), nil)
}

func op_gl_Clear(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.Clear(uint32(ReadI32(fp, inp1)))
}

func op_gl_ClearColor(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.ClearColor(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3), ReadF32(fp, inp4))
}

func op_gl_ClearStencil(expr *CXExpression, fp int) {
	inp0 := expr.Inputs[0]
	gl.ClearStencil(ReadI32(fp, inp0))
}

func op_gl_ClearDepth(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.ClearDepth(ReadF64(fp, inp1))
}

func op_gl_StencilMask(expr *CXExpression, fp int) {
	inp0 := expr.Inputs[0]
	gl.StencilMask(uint32(ReadI32(fp, inp0)))
}

func op_gl_ColorMask(expr *CXExpression, fp int) {
	inp0, inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.ColorMask(ReadBool(fp, inp0), ReadBool(fp, inp1), ReadBool(fp, inp2), ReadBool(fp, inp3))
}

func op_gl_DepthMask(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.DepthMask(ReadBool(fp, inp1))
}

func op_gl_Disable(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.Disable(uint32(ReadI32(fp, inp1)))
}

func op_gl_Enable(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.Enable(uint32(ReadI32(fp, inp1)))
}

func op_gl_BlendFunc(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.BlendFunc(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_StencilFunc(expr *CXExpression, fp int) {
	inp0, inp1, inp2 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.StencilFunc(uint32(ReadI32(fp, inp0)), ReadI32(fp, inp1), uint32(ReadI32(fp, inp2)))
}

func op_gl_StencilOp(expr *CXExpression, fp int) {
	inp0, inp1, inp2 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.StencilOp(uint32(ReadI32(fp, inp0)), uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_DepthFunc(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.DepthFunc(uint32(ReadI32(fp, inp1)))
}

func op_gl_GetError(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI32(int32(gl.GetError()))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_GetTexLevelParameteriv(expr *CXExpression, fp int) {
	inp1, inp2, inp3, out1 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Outputs[0]
	var outValue int32 = 0
	gl.GetTexLevelParameteriv(uint32(ReadI32(fp, inp1)), ReadI32(fp, inp2), uint32(ReadI32(fp, inp3)), &outValue)
	WriteMemory(GetFinalOffset(fp, out1), FromI32(outValue))
}

func op_gl_DepthRange(expr *CXExpression, fp int) {
	inp0, inp1 := expr.Inputs[0], expr.Inputs[1]
	gl.DepthRange(ReadF64(fp, inp0), ReadF64(fp, inp1))
}

func op_gl_Viewport(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.Viewport(ReadI32(fp, inp1), ReadI32(fp, inp2), ReadI32(fp, inp3), ReadI32(fp, inp4))
}

// gl_1_1
func op_gl_DrawArrays(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.DrawArrays(uint32(ReadI32(fp, inp1)), ReadI32(fp, inp2), ReadI32(fp, inp3))
}

func op_gl_BindTexture(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.BindTexture(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_DeleteTextures(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	tmp := uint32(ReadI32(fp, inp2))
	gl.DeleteTextures(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
}

func op_gl_GenTextures(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	tmp := uint32(ReadI32(fp, inp2))
	gl.GenTextures(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
	outB1 := FromI32(int32(tmp))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// gl_1_3
func op_gl_ActiveTexture(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.ActiveTexture(uint32(ReadI32(fp, inp1)))
}

// gl_1_5
func op_gl_BindBuffer(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.BindBuffer(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_DeleteBuffers(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	tmp := uint32(ReadI32(fp, inp2))
	gl.DeleteBuffers(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
}

func op_gl_GenBuffers(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	tmp := uint32(ReadI32(fp, inp2))
	gl.GenBuffers(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
	outB1 := FromI32(int32(tmp))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_BufferData(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.BufferData(uint32(ReadI32(fp, inp1)), int(ReadI32(fp, inp2)), gl.Ptr(ReadF32Data(fp, inp3)), uint32(ReadI32(fp, inp4)))
}

func op_gl_BufferSubData(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.BufferSubData(uint32(ReadI32(fp, inp1)), int(ReadI32(fp, inp2)), int(ReadI32(fp, inp3)), gl.Ptr(ReadF32Data(fp, inp4)))
}

// gl_2_0
func op_gl_StencilOpSeparate(expr *CXExpression, fp int) {
	inp0, inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.StencilOpSeparate(uint32(ReadI32(fp, inp0)), uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)), uint32(ReadI32(fp, inp3)))
}

func op_gl_StencilFuncSeparate(expr *CXExpression, fp int) {
	inp0, inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.StencilFuncSeparate(uint32(ReadI32(fp, inp0)), uint32(ReadI32(fp, inp1)), ReadI32(fp, inp2), uint32(ReadI32(fp, inp3)))
}

func op_gl_StencilMaskSeparate(expr *CXExpression, fp int) {
	inp0, inp1 := expr.Inputs[0], expr.Inputs[1]
	gl.StencilMaskSeparate(uint32(ReadI32(fp, inp0)), uint32(ReadI32(fp, inp1)))
}

func op_gl_AttachShader(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.AttachShader(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_BindAttribLocation(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	xstr := cSources[ReadStr(fp, inp3)]
	gl.BindAttribLocation(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)), *xstr)
}

func op_gl_CompileShader(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	shad := uint32(ReadI32(fp, inp1))
	gl.CompileShader(shad)

	var status int32
	gl.GetShaderiv(shad, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shad, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shad, logLength, nil, gl.Str(log))

		fmt.Printf("failed to compile: %v", log)
	}
}

func op_gl_CreateProgram(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI32(int32(gl.CreateProgram()))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_CreateShader(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(gl.CreateShader(uint32(ReadI32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_DeleteProgram(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.DeleteShader(uint32(ReadI32(fp, inp1)))
}

func op_gl_DeleteShader(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.DeleteShader(uint32(ReadI32(fp, inp1)))
}

func op_gl_DetachShader(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.DetachShader(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_EnableVertexAttribArray(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.EnableVertexAttribArray(uint32(ReadI32(fp, inp1)))
}

func op_gl_GetAttribLocation(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	xstr := cSources[ReadStr(fp, inp2)]
	outB1 := FromI32(gl.GetAttribLocation(uint32(ReadI32(fp, inp1)), *xstr))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_GetShaderiv(expr *CXExpression, fp int) {
	// pointers
	panic("gl.GetShaderiv")
}

func op_gl_GetUniformLocation(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	xstr := cSources[ReadStr(fp, inp2)]
	outB1 := FromI32(gl.GetUniformLocation(uint32(ReadI32(fp, inp1)), *xstr))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_LinkProgram(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.LinkProgram(uint32(ReadI32(fp, inp1)))
}

func op_gl_ShaderSource(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	xstr := cSources[ReadStr(fp, inp3)]
	gl.ShaderSource(uint32(ReadI32(fp, inp1)), ReadI32(fp, inp2), xstr, nil)
}

func op_gl_UseProgram(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.UseProgram(uint32(ReadI32(fp, inp1)))
}

func op_gl_Uniform1f(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.Uniform1f(ReadI32(fp, inp1), ReadF32(fp, inp2))
}

func op_gl_Uniform1i(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.Uniform1i(ReadI32(fp, inp1), ReadI32(fp, inp2))
}

func op_gl_VertexAttribPointer(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4]
	gl.VertexAttribPointer(uint32(ReadI32(fp, inp1)), ReadI32(fp, inp2), uint32(ReadI32(fp, inp3)), ReadBool(fp, inp4), ReadI32(fp, inp5), nil)
}

func op_gl_VertexAttribPointerI32(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]
	gl.VertexAttribPointer(uint32(ReadI32(fp, inp1)), ReadI32(fp, inp2), uint32(ReadI32(fp, inp3)), ReadBool(fp, inp4), ReadI32(fp, inp5), gl.PtrOffset(int(ReadI32(fp, inp6))))
}

// gl_3_0
func op_gl_BindRenderbuffer(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.BindRenderbuffer(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_DeleteRenderbuffers(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	tmp := uint32(ReadI32(fp, inp2))
	gl.DeleteRenderbuffers(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
}

func op_gl_GenRenderbuffers(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	tmp := uint32(ReadI32(fp, inp2))
	gl.GenRenderbuffers(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
	outB1 := FromI32(int32(tmp))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_RenderbufferStorage(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.RenderbufferStorage(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)), ReadI32(fp, inp3), ReadI32(fp, inp4))
}

func op_gl_BindFramebuffer(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.BindFramebuffer(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_DeleteFramebuffers(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	tmp := uint32(ReadI32(fp, inp2))
	gl.DeleteFramebuffers(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
}

func op_gl_GenFramebuffers(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	tmp := uint32(ReadI32(fp, inp2))
	gl.GenFramebuffers(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
	outB1 := FromI32(int32(tmp))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_CheckFramebufferStatus(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(gl.CheckFramebufferStatus(uint32(ReadI32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_FramebufferTexture2D(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4]
	gl.FramebufferTexture2D(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)), uint32(ReadI32(fp, inp3)), uint32(ReadI32(fp, inp4)), ReadI32(fp, inp5))
}

func op_gl_FramebufferRenderbuffer(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.FramebufferRenderbuffer(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)), uint32(ReadI32(fp, inp3)), uint32(ReadI32(fp, inp4)))
}

func op_gl_BindVertexArray(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	if runtime.GOOS == "darwin" {
		gl.BindVertexArrayAPPLE(uint32(ReadI32(fp, inp1)))
	} else {
		gl.BindVertexArray(uint32(ReadI32(fp, inp1)))
	}
}

func op_gl_DeleteVertexArrays(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	tmp := uint32(ReadI32(fp, inp2))
	if runtime.GOOS == "darwin" {
		gl.DeleteVertexArraysAPPLE(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
	} else {
		gl.DeleteVertexArrays(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
	}
}

func op_gl_GenVertexArrays(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	tmp := uint32(ReadI32(fp, inp2))
	if runtime.GOOS == "darwin" {
		gl.GenVertexArraysAPPLE(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
	} else {
		gl.GenVertexArrays(ReadI32(fp, inp1), &tmp) // will panic if inp1 > 1
	}
	outB1 := FromI32(int32(tmp))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
