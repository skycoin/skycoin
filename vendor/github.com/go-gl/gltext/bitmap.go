// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gltext

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

// LoadBitmap loads a bitmap (raster) font from the given 
// sprite sheet and config files. It is optionally scaled by
// the given scale factor.
//
// A scale factor of 1 retains the original size. A factor of 2 doubles the 
// font size, etc. A scale factor of 0 is not valid and will default to 1.
//
// Supported image formats are 32-bit RGBA as PNG, JPEG and GIF.
func LoadBitmap(img, config io.Reader, scale int) (*Font, error) {
	pix, _, err := image.Decode(img)
	if err != nil {
		return nil, err
	}

	rgba := toRGBA(pix, scale)

	var fc FontConfig
	err = fc.Load(config)

	if err != nil {
		return nil, err
	}

	fc.Glyphs.Scale(scale)
	return loadFont(rgba, &fc)
}
