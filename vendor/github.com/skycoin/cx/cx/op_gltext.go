// +build extra full

package base

import (
	"unicode/utf8"

	"github.com/go-gl/gltext"
)

var fonts map[string]*gltext.Font = make(map[string]*gltext.Font, 0)

func op_gltext_LoadTrueType(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]

	if theFont, err := gltext.LoadTruetype(openFiles[ReadStr(fp, inp2)], ReadI32(fp, inp3), rune(ReadI32(fp, inp4)), rune(ReadI32(fp, inp5)), gltext.Direction(ReadI32(fp, inp6))); err == nil {
		fonts[ReadStr(fp, inp1)] = theFont
	} else {
		panic(err)
	}
}

func op_gltext_Printf(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]

	if err := fonts[ReadStr(fp, inp1)].Printf(ReadF32(fp, inp2), ReadF32(fp, inp3), ReadStr(fp, inp4)); err != nil {
		panic(err)
	}
}

func op_gltext_Metrics(expr *CXExpression, fp int) {
	inp1, inp2, out1, out2 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0], expr.Outputs[1]

	width, height := fonts[ReadStr(fp, inp1)].Metrics(ReadStr(fp, inp2))

	WriteMemory(GetFinalOffset(fp, out1), FromI32(int32(width)))
	WriteMemory(GetFinalOffset(fp, out2), FromI32(int32(height)))
}

func op_gltext_Texture(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	WriteMemory(GetFinalOffset(fp, out1), FromI32(int32(fonts[ReadStr(fp, inp1)].Texture())))
}

func op_gltext_NextGlyph(expr *CXExpression, fp int) { // refactor
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	out1, out2, out3, out4, out5, out6, out7 := expr.Outputs[0], expr.Outputs[1], expr.Outputs[2], expr.Outputs[3], expr.Outputs[4], expr.Outputs[5], expr.Outputs[6]
	font := fonts[ReadStr(fp, inp1)]
	str := ReadStr(fp, inp2)
	var index int = int(ReadI32(fp, inp3))
	var runeValue rune = -1
	var width int = -1
	var x int = 0
	var y int = 0
	var w int = 0
	var h int = 0
	var advance int = 0
	if index < len(str) {
		runeValue, width = utf8.DecodeRuneInString(str[index:])
		g := font.Glyphs()[runeValue-font.Low()]
		x = g.X
		y = g.Y
		w = g.Width
		h = g.Height
		advance = g.Advance
	}

	WriteMemory(GetFinalOffset(fp, out1), FromI32(int32(runeValue-font.Low())))
	WriteMemory(GetFinalOffset(fp, out2), FromI32(int32(width)))
	WriteMemory(GetFinalOffset(fp, out3), FromI32(int32(x)))
	WriteMemory(GetFinalOffset(fp, out4), FromI32(int32(y)))
	WriteMemory(GetFinalOffset(fp, out5), FromI32(int32(w)))
	WriteMemory(GetFinalOffset(fp, out6), FromI32(int32(h)))
	WriteMemory(GetFinalOffset(fp, out7), FromI32(int32(advance)))
}

func op_gltext_GlyphBounds(expr *CXExpression, fp int) {
	inp1, out1, out2 := expr.Inputs[0], expr.Outputs[0], expr.Outputs[1]
	font := fonts[ReadStr(fp, inp1)]
	var maxGlyphWidth, maxGlyphHeight int = font.GlyphBounds()
	WriteMemory(GetFinalOffset(fp, out1), FromI32(int32(maxGlyphWidth)))
	WriteMemory(GetFinalOffset(fp, out2), FromI32(int32(maxGlyphHeight)))
}

func op_gltext_GlyphMetrics(expr *CXExpression, fp int) { // refactor
	inp1, inp2, out1, out2 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0], expr.Outputs[1]

	width, height := fonts[ReadStr(fp, inp1)].GlyphMetrics(uint32(ReadI32(fp, inp2)))

	WriteMemory(GetFinalOffset(fp, out1), FromI32(int32(width)))
	WriteMemory(GetFinalOffset(fp, out2), FromI32(int32(height)))
}

func op_gltext_GlyphInfo(expr *CXExpression, fp int) { // refactor
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	out1, out2, out3, out4, out5 := expr.Outputs[0], expr.Outputs[1], expr.Outputs[2], expr.Outputs[3], expr.Outputs[4]
	font := fonts[ReadStr(fp, inp1)]
	glyph := ReadI32(fp, inp2)
	var x int = 0
	var y int = 0
	var w int = 0
	var h int = 0
	var advance int = 0
	g := font.Glyphs()[glyph]
	x = g.X
	y = g.Y
	w = g.Width
	h = g.Height
	advance = g.Advance

	WriteMemory(GetFinalOffset(fp, out1), FromI32(int32(x)))
	WriteMemory(GetFinalOffset(fp, out2), FromI32(int32(y)))
	WriteMemory(GetFinalOffset(fp, out3), FromI32(int32(w)))
	WriteMemory(GetFinalOffset(fp, out4), FromI32(int32(h)))
	WriteMemory(GetFinalOffset(fp, out5), FromI32(int32(advance)))
}
