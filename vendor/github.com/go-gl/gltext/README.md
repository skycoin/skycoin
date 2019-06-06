## gltext

**Note**: This package is experimental and subject to change.
Use at your own discretion.

The gltext package offers a simple set of text rendering utilities for OpenGL
programs. It deals with TrueType and Bitmap (raster) fonts. Text can be
rendered in various directions (Left-to-right, right-to-left and top-to-bottom).
This allows for correct display of text for various languages.

The package supports the full set of unicode characters, provided the loaded
font does as well.


### TODO

* Have a look at Valve's 'Signed Distance Field` techniques to render
  sharp font textures are different zoom levels.

  * [SIGGRAPH2007_AlphaTestedMagnification.pdf](http://www.valvesoftware.com/publications/2007/SIGGRAPH2007_AlphaTestedMagnification.pdf)
  * [Youtube video](http://www.youtube.com/watch?v=CGZRHJvJYIg)
  
  More links to info in the youtube video description.
  An alternative might be a port of [GLyphy](http://code.google.com/p/glyphy/)


### Known bugs

* Determining the height of truetype glyphs is not entirely accurate.
  It is unclear at this point how to get to this information reliably.
  Specifically the parts in `LoadTruetype` at truetype.go#L54+.
  The vertical glyph bounds computed by freetype-go are not correct for
  certain fonts. Right now we manually offset the value by added `4` to
  the height. This is an unreliable hack and should be fixed.
* `freetype-go` does not expose `AdvanceHeight` for vertically rendered fonts.
  This may mean that the Advance size for top-to-bottom fonts is incorrect.


### Dependencies

This packages uses [freetype-go](https://code.google.com/p/freetype-go) which is licensed 
under GPLv2 e FTL licenses. You can choose which one is a better fit for your 
use case but FTL requires you to give some form of credit to Freetype.org

You can read the [GPLv2](https://code.google.com/p/freetype-go/source/browse/licenses/gpl.txt)
and [FTL](https://code.google.com/p/freetype-go/source/browse/licenses/ftl.txt)
licenses for more information about the requirements.

### Usage

    go get github.com/go-gl/gltext

Refer to [go-gl/examples/gltext][ex] for usage examples.

[ex]: https://github.com/go-gl/examples/tree/64b743f99c4e9151c09563e9be3339441eb9296b/gltext


### License

Copyright 2012 The go-gl Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.

