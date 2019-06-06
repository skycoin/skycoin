// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gltext

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
)

// A Font allows rendering of text to an OpenGL context.
type Font struct {
	config         *FontConfig // Character set for this font.
	texture        uint32      // Holds the glyph texture id.
	listbase       uint32      // Holds the first display list id.
	maxGlyphWidth  int         // Largest glyph width.
	maxGlyphHeight int         // Largest glyph height.
}

// loadFont loads the given font data. This does not deal with font scaling.
// Scaling should be handled by the independent Bitmap/Truetype loaders.
// We therefore expect the supplied image and charset to already be adjusted
// to the correct font scale.
//
// The image should hold a sprite sheet, defining the graphical layout for
// every glyph. The config describes font metadata.
func loadFont(img *image.RGBA, config *FontConfig) (f *Font, err error) {
	f = new(Font)
	f.config = config

	// Resize image to next power-of-two.
	img = Pow2Image(img).(*image.RGBA)
	ib := img.Bounds()

	// Create the texture itself. It will contain all glyphs.
	// Individual glyph-quads display a subset of this texture.
	gl.GenTextures(1, &f.texture)
	gl.BindTexture(gl.TEXTURE_2D, f.texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(ib.Dx()), int32(ib.Dy()), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	// Create display lists for each glyph.
	f.listbase = gl.GenLists(int32(len(config.Glyphs)))

	texWidth := float32(ib.Dx())
	texHeight := float32(ib.Dy())

	for index, glyph := range config.Glyphs {
		// Update max glyph bounds.
		if glyph.Width > f.maxGlyphWidth {
			f.maxGlyphWidth = glyph.Width
		}

		if glyph.Height > f.maxGlyphHeight {
			f.maxGlyphHeight = glyph.Height
		}

		// Quad width/height
		vw := float32(glyph.Width)
		vh := float32(glyph.Height)

		// Texture coordinate offsets.
		tx1 := float32(glyph.X) / texWidth
		ty1 := float32(glyph.Y) / texHeight
		tx2 := (float32(glyph.X) + vw) / texWidth
		ty2 := (float32(glyph.Y) + vh) / texHeight

		// Advance width (or height if we render top-to-bottom)
		adv := float32(glyph.Advance)

		gl.NewList(f.listbase+uint32(index), gl.COMPILE)
		{
			gl.Begin(gl.QUADS)
			{
				gl.TexCoord2f(tx1, ty2)
				gl.Vertex2f(0, 0)
				gl.TexCoord2f(tx2, ty2)
				gl.Vertex2f(vw, 0)
				gl.TexCoord2f(tx2, ty1)
				gl.Vertex2f(vw, vh)
				gl.TexCoord2f(tx1, ty1)
				gl.Vertex2f(0, vh)
			}
			gl.End()

			switch config.Dir {
			case LeftToRight:
				gl.Translatef(adv, 0, 0)
			case RightToLeft:
				gl.Translatef(-adv, 0, 0)
			case TopToBottom:
				gl.Translatef(0, -adv, 0)
			}
		}
		gl.EndList()
	}

	err = checkGLError()
	return
}

// Dir returns the font's rendering orientation.
func (f *Font) Dir() Direction { return f.config.Dir }

// Low returns the font's lower rune bound.
func (f *Font) Low() rune { return f.config.Low }

// High returns the font's upper rune bound.
func (f *Font) High() rune { return f.config.High }

// Glyphs returns the font's glyph descriptors.
func (f *Font) Glyphs() Charset { return f.config.Glyphs }

// Release releases font resources.
// A font can no longer be used for rendering after this call completes.
func (f *Font) Release() {
	gl.DeleteTextures(1, &f.texture)
	gl.DeleteLists(f.listbase, int32(len(f.config.Glyphs)))
	f.config = nil
}

// Metrics returns the pixel width and height for the given string.
// This takes the scale and rendering direction of the font into account.
//
// Unknown runes will be counted as having the maximum glyph bounds as
// defined by Font.GlyphBounds().
func (f *Font) Metrics(text string) (int, int) {
	if len(text) == 0 {
		return 0, 0
	}

	gw, gh := f.GlyphBounds()

	if f.config.Dir == TopToBottom {
		return gw, f.advanceSize(text)
	}

	return f.advanceSize(text), gh
}

// advanceSize computes the pixel width or height for the given single-line
// input string. This iterates over all of its runes, finds the matching
// Charset entry and adds up the Advance values.
//
// Unknown runes will be counted as having the maximum glyph bounds as
// defined by Font.GlyphBounds().
func (f *Font) advanceSize(line string) int {
	gw, gh := f.GlyphBounds()
	glyphs := f.config.Glyphs
	low := f.config.Low
	indices := []rune(line)

	var size int
	for _, r := range indices {
		r -= low

		if r >= 0 && int(r) < len(glyphs) {
			size += glyphs[r].Advance
			continue
		}

		if f.config.Dir == TopToBottom {
			size += gh
		} else {
			size += gw
		}
	}

	return size
}

// Printf draws the given string at the specified coordinates.
// It expects the string to be a single line. Line breaks are not
// handled as line breaks and are rendered as glyphs.
//
// In order to render multi-line text, it is up to the caller to split
// the text up into individual lines of adequate length and then call
// this method for each line seperately.
func (f *Font) Printf(x, y float32, fs string, argv ...interface{}) error {
	indices := []rune(fmt.Sprintf(fs, argv...))

	if len(indices) == 0 {
		return nil
	}

	// Runes form display list indices.
	// For this purpose, they need to be offset by -FontConfig.Low
	low := f.config.Low
	for i := range indices {
		indices[i] -= low
	}

	var vp [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &vp[0])

	gl.PushAttrib(gl.TRANSFORM_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(float64(vp[0]), float64(vp[2]), float64(vp[1]), float64(vp[3]), 0, 1)
	gl.PopAttrib()

	gl.PushAttrib(gl.LIST_BIT | gl.CURRENT_BIT | gl.ENABLE_BIT | gl.TRANSFORM_BIT)
	{
		gl.MatrixMode(gl.MODELVIEW)
		gl.Disable(gl.LIGHTING)
		gl.Disable(gl.DEPTH_TEST)
		gl.Enable(gl.BLEND)
		gl.Enable(gl.TEXTURE_2D)

		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)
		gl.BindTexture(gl.TEXTURE_2D, f.texture)
		gl.ListBase(f.listbase)

		var mv [16]float32
		gl.GetFloatv(gl.MODELVIEW_MATRIX, &mv[0])

		gl.PushMatrix()
		{
			gl.LoadIdentity()

			mgw := float32(f.maxGlyphWidth)
			mgh := float32(f.maxGlyphHeight)

			switch f.config.Dir {
			case LeftToRight, TopToBottom:
				gl.Translatef(x, float32(vp[3])-y-mgh, 0)
			case RightToLeft:
				gl.Translatef(x-mgw, float32(vp[3])-y-mgh, 0)
			}

			gl.MultMatrixf(&mv[0])
			gl.CallLists(int32(len(indices)), gl.UNSIGNED_INT, unsafe.Pointer(&indices[0]))
		}
		gl.PopMatrix()
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}
	gl.PopAttrib()

	gl.PushAttrib(gl.TRANSFORM_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.PopAttrib()
	return checkGLError()
}

// GlyphBounds returns the largest width and height for any of the glyphs
// in the font. This constitutes the largest possible bounding box
// a single glyph will have.
func (f *Font) GlyphBounds() (int, int) {
	return f.maxGlyphWidth, f.maxGlyphHeight
}
