package cli

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
)

// PinMatrix displays a 3x3 PIN entry matrix using tcell
// Returns the entered PIN as a string of position numbers (1-9)
// The matrix layout matches the device display:
//
//	7 8 9
//	4 5 6
//	1 2 3
//
// User navigates with arrow keys, selects with Enter/Space, backspace to delete, Escape to cancel
func PinMatrix() (string, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return "", fmt.Errorf("failed to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		return "", fmt.Errorf("failed to init screen: %v", err)
	}
	defer screen.Fini()

	// Style definitions
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Background(tcell.NewRGBColor(0, 0, 139)).Foreground(tcell.NewRGBColor(255, 255, 255))
	selectedStyle := tcell.StyleDefault.Background(tcell.NewRGBColor(0, 128, 0)).Foreground(tcell.NewRGBColor(0, 0, 0))
	titleStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.NewRGBColor(255, 255, 0)).Bold(true)
	pinStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.NewRGBColor(0, 128, 0))
	helpStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.NewRGBColor(128, 128, 128))

	screen.SetStyle(defStyle)
	screen.Clear()

	// Current selection position (0-8, where 0=top-left, 8=bottom-right)
	cursorPos := 4 // Start in center

	// Entered PIN digits
	var pin string

	// Position to digit mapping (matches device layout):
	// Position 0,1,2 (top row) = digits 7,8,9
	// Position 3,4,5 (middle row) = digits 4,5,6
	// Position 6,7,8 (bottom row) = digits 1,2,3
	posToDigit := []byte{'7', '8', '9', '4', '5', '6', '1', '2', '3'}

	draw := func() {
		screen.Clear()
		w, h := screen.Size()

		// Calculate center position
		boxWidth := 21  // 3 cells * 5 chars + 2 borders + 4 separators
		boxHeight := 11 // 3 cells * 3 lines + 2 borders
		startX := (w - boxWidth) / 2
		startY := (h - boxHeight) / 2
		if startX < 0 {
			startX = 0
		}
		if startY < 2 {
			startY = 2
		}

		// Draw title
		title := "PIN Entry Matrix"
		titleX := (w - len(title)) / 2
		for i, r := range title {
			screen.SetContent(titleX+i, startY-2, r, nil, titleStyle)
		}

		// Draw instructions
		help := "Look at device screen for digit positions"
		helpX := (w - len(help)) / 2
		for i, r := range help {
			screen.SetContent(helpX+i, startY-1, r, nil, helpStyle)
		}

		// Draw the 3x3 grid
		// Top border
		drawHLine(screen, startX, startY, boxWidth, boxStyle)

		for row := 0; row < 3; row++ {
			y := startY + 1 + row*3

			// Cell content rows (3 lines per cell)
			for line := 0; line < 3; line++ {
				screen.SetContent(startX, y+line, '|', nil, boxStyle)
				for col := 0; col < 3; col++ {
					pos := row*3 + col
					cellX := startX + 1 + col*7

					style := boxStyle
					if pos == cursorPos {
						style = selectedStyle
					}

					// Draw cell content (5 chars wide, 3 lines tall)
					for cx := 0; cx < 5; cx++ {
						if line == 1 && cx == 2 {
							// Center of cell - show a dot or the position hint
							screen.SetContent(cellX+cx, y+line, '*', nil, style)
						} else {
							screen.SetContent(cellX+cx, y+line, ' ', nil, style)
						}
					}

					// Draw separator
					screen.SetContent(cellX+5, y+line, '|', nil, boxStyle)
				}
			}

			// Horizontal separator (except after last row)
			if row < 2 {
				drawHLine(screen, startX, y+3, boxWidth, boxStyle)
			}
		}

		// Bottom border
		drawHLine(screen, startX, startY+boxHeight-1, boxWidth, boxStyle)

		// Draw entered PIN (as asterisks)
		pinLabel := "PIN: "
		pinY := startY + boxHeight + 1
		pinX := (w - len(pinLabel) - 9) / 2
		for i, r := range pinLabel {
			screen.SetContent(pinX+i, pinY, r, nil, defStyle)
		}
		for i := 0; i < len(pin); i++ {
			screen.SetContent(pinX+len(pinLabel)+i, pinY, '*', nil, pinStyle)
		}
		for i := len(pin); i < 9; i++ {
			screen.SetContent(pinX+len(pinLabel)+i, pinY, '_', nil, helpStyle)
		}

		// Draw help text
		help2 := "[Arrows] Move  [Enter/Space] Select  [Backspace] Delete  [Esc] Cancel"
		help2X := (w - len(help2)) / 2
		for i, r := range help2 {
			screen.SetContent(help2X+i, pinY+2, r, nil, helpStyle)
		}

		screen.Show()
	}

	// Main event loop
	for {
		draw()

		ev := <-screen.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()

		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return "", fmt.Errorf("canceled")

			case tcell.KeyEnter, tcell.KeyCtrlJ:
				if len(pin) < 9 {
					pin += string(posToDigit[cursorPos])
				}

			case tcell.KeyBackspace, tcell.KeyDelete:
				if len(pin) > 0 {
					pin = pin[:len(pin)-1]
				}

			case tcell.KeyUp:
				if cursorPos >= 3 {
					cursorPos -= 3
				}

			case tcell.KeyDown:
				if cursorPos < 6 {
					cursorPos += 3
				}

			case tcell.KeyLeft:
				if cursorPos%3 > 0 {
					cursorPos--
				}

			case tcell.KeyRight:
				if cursorPos%3 < 2 {
					cursorPos++
				}

			case tcell.KeyCtrlC:
				return "", fmt.Errorf("interrupted")

			case tcell.KeyRune:
				ch := ev.Str()
				if len(ch) > 0 {
					switch ch[0] {
					case ' ':
						// Space also selects
						if len(pin) < 9 {
							pin += string(posToDigit[cursorPos])
						}
					case 'q', 'Q':
						return "", fmt.Errorf("canceled")
					case '1', '2', '3', '4', '5', '6', '7', '8', '9':
						// Direct digit entry (for testing/convenience)
						if len(pin) < 9 {
							pin += string(ch[0])
						}
					}
				}
			}
		}
	}
}

// drawHLine draws a horizontal line
func drawHLine(screen tcell.Screen, x, y, width int, style tcell.Style) {
	for i := 0; i < width; i++ {
		screen.SetContent(x+i, y, '-', nil, style)
	}
}

// PinMatrixSimple is a simpler version that collects the full PIN before returning
// Uses Enter to submit the complete PIN
func PinMatrixSimple() (string, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return "", fmt.Errorf("failed to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		return "", fmt.Errorf("failed to init screen: %v", err)
	}
	defer screen.Fini()

	// Enable mouse support
	screen.EnableMouse()

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Background(tcell.NewRGBColor(0, 0, 139)).Foreground(tcell.NewRGBColor(255, 255, 255))
	selectedStyle := tcell.StyleDefault.Background(tcell.NewRGBColor(0, 128, 0)).Foreground(tcell.NewRGBColor(0, 0, 0))
	titleStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.NewRGBColor(255, 255, 0)).Bold(true)
	pinStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.NewRGBColor(0, 128, 0))
	helpStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.NewRGBColor(128, 128, 128))

	screen.SetStyle(defStyle)

	cursorPos := 4
	var pin string
	posToDigit := []byte{'7', '8', '9', '4', '5', '6', '1', '2', '3'}

	submitted := false

	// Store grid position for mouse detection
	var gridStartX, gridStartY, cellWidth, cellHeight int

	for !submitted {
		screen.Clear()
		w, h := screen.Size()

		boxWidth := 21
		boxHeight := 11
		startX := (w - boxWidth) / 2
		startY := (h - boxHeight) / 2
		if startX < 0 {
			startX = 0
		}
		if startY < 2 {
			startY = 2
		}

		// Store for mouse hit detection
		gridStartX = startX
		gridStartY = startY
		cellWidth = 7  // 5 chars + 2 for borders
		cellHeight = 3 // 3 lines per cell

		// Title
		title := "PIN Entry - Look at device for numbers"
		titleX := (w - len(title)) / 2
		for i, r := range title {
			screen.SetContent(titleX+i, startY-2, r, nil, titleStyle)
		}

		// Grid top border
		for i := 0; i < boxWidth; i++ {
			screen.SetContent(startX+i, startY, '-', nil, boxStyle)
		}

		// Grid cells
		for row := 0; row < 3; row++ {
			y := startY + 1 + row*3
			for line := 0; line < 3; line++ {
				screen.SetContent(startX, y+line, '|', nil, boxStyle)
				for col := 0; col < 3; col++ {
					pos := row*3 + col
					cellX := startX + 1 + col*7
					style := boxStyle
					if pos == cursorPos {
						style = selectedStyle
					}
					for cx := 0; cx < 5; cx++ {
						ch := ' '
						if line == 1 && cx == 2 {
							ch = '*'
						}
						screen.SetContent(cellX+cx, y+line, ch, nil, style)
					}
					screen.SetContent(cellX+5, y+line, '|', nil, boxStyle)
				}
			}
			if row < 2 {
				for i := 0; i < boxWidth; i++ {
					screen.SetContent(startX+i, y+3, '-', nil, boxStyle)
				}
			}
		}

		// Bottom border
		for i := 0; i < boxWidth; i++ {
			screen.SetContent(startX+i, startY+boxHeight-1, '-', nil, boxStyle)
		}

		// PIN display
		pinY := startY + boxHeight + 1
		pinLabel := "PIN: "
		pinX := (w - len(pinLabel) - 9) / 2
		for i, r := range pinLabel {
			screen.SetContent(pinX+i, pinY, r, nil, defStyle)
		}
		for i := 0; i < len(pin); i++ {
			screen.SetContent(pinX+len(pinLabel)+i, pinY, '*', nil, pinStyle)
		}
		for i := len(pin); i < 9; i++ {
			screen.SetContent(pinX+len(pinLabel)+i, pinY, '_', nil, helpStyle)
		}

		// Help
		help := "[Click/Space] Add digit  [Enter] Submit  [Backspace] Delete  [Esc] Cancel"
		helpX := (w - len(help)) / 2
		if helpX < 0 {
			helpX = 0
		}
		for i, r := range help {
			if helpX+i < w {
				screen.SetContent(helpX+i, pinY+2, r, nil, helpStyle)
			}
		}

		screen.Show()

		ev := <-screen.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventMouse:
			mx, my := ev.Position()
			buttons := ev.Buttons()

			// Check for left click
			if buttons&tcell.Button1 != 0 {
				// Check if click is within the grid
				// Grid starts at gridStartY+1 and each cell is cellHeight tall
				// Grid starts at gridStartX+1 and each cell is cellWidth wide
				relX := mx - gridStartX - 1
				relY := my - gridStartY - 1

				if relX >= 0 && relY >= 0 {
					col := relX / cellWidth
					row := relY / cellHeight

					if col >= 0 && col < 3 && row >= 0 && row < 3 {
						clickedPos := row*3 + col
						if len(pin) < 9 {
							pin += string(posToDigit[clickedPos])
						}
					}
				}
			}
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return "", fmt.Errorf("canceled")
			case tcell.KeyEnter:
				if len(pin) > 0 {
					submitted = true
				}
			case tcell.KeyBackspace:
				if len(pin) > 0 {
					pin = pin[:len(pin)-1]
				}
			case tcell.KeyUp:
				if cursorPos >= 3 {
					cursorPos -= 3
				}
			case tcell.KeyDown:
				if cursorPos < 6 {
					cursorPos += 3
				}
			case tcell.KeyLeft:
				if cursorPos%3 > 0 {
					cursorPos--
				}
			case tcell.KeyRight:
				if cursorPos%3 < 2 {
					cursorPos++
				}
			case tcell.KeyRune:
				ch := ev.Str()
				if len(ch) > 0 {
					switch ch[0] {
					case ' ':
						if len(pin) < 9 {
							pin += string(posToDigit[cursorPos])
						}
					case 'q', 'Q':
						return "", fmt.Errorf("canceled")
					}
				}
			}
		}
	}

	return pin, nil
}
