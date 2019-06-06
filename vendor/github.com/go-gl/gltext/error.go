// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gltext

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
)

// checkGLError returns an opengl error if one exists.
func checkGLError() error {
	errno := gl.GetError()
	if errno == gl.NO_ERROR {
		return nil
	}
	return fmt.Errorf("GL error: %d", errno)
}
