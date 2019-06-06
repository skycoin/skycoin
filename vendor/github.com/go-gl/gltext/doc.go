// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
The gltext package offers a set of text rendering utilities for OpenGL
programs. It deals with TrueType and Bitmap (raster) fonts.

Text can be rendered in predefined directions (Left-to-right, right-to-left and
top-to-bottom). This allows for correct display of text for various languages.

This package supports the full set of unicode characters, provided the loaded
font does as well.

This packages uses freetype-go (code.google.com/p/freetype-go) which is licensed 
under GPLv2 e FTL licenses. You can choose which one is a better fit for your 
use case but FTL requires you to give some form of credit to Freetype.org

You can read the GPLv2 (https://code.google.com/p/freetype-go/source/browse/licenses/gpl.txt)
and FTL (https://code.google.com/p/freetype-go/source/browse/licenses/ftl.txt)
licenses for more information about the requirements.
*/
package gltext
