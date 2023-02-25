package main

import (
	"image/color"
	"machine"
	"runtime/interrupt"
	"runtime/volatile"
	"unsafe"

	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
)

// KeyCode represents a GBA keypad code.
type KeyCode uint16

var (

	// Register display
	regDISPSTAT = (*volatile.Register16)(unsafe.Pointer(uintptr(0x4000004)))

	// Register keypad
	regKEYPAD = (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000130)))

	// Display from machine
	display = machine.Display

	// Screen resolution
	screenWidth  int16 = 240
	screenHeight int16 = 160

	version = "v0.1"

	// Colors
	black = color.RGBA{}
	white = color.RGBA{255, 255, 255, 255}
	green = color.RGBA{0, 255, 0, 255}
	red   = color.RGBA{255, 0, 0, 255}

	// Coordinates
	x int16 = screenWidth / 2
	y int16 = screenHeight / 2

	// Border width
	border int16 = 16

	// input represents the register for keypad input.
	//
	// The GBA is a bit strange and has "inverted" values for input, 0 = pressed
	// and 1 = unpressed.
	input = (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000130)))

	//KeyCodes / Buttons
	KeyA      = KeyCode(0x0001)
	KeyB      = KeyCode(0x0002)
	KeySelect = KeyCode(0x0004)
	KeyStart  = KeyCode(0x0008)
	KeyRight  = KeyCode(0x0010)
	KeyLeft   = KeyCode(0x0020)
	KeyUp     = KeyCode(0x0040)
	KeyDown   = KeyCode(0x0080)
	KeyR      = KeyCode(0x0100)
	KeyL      = KeyCode(0x0200)

	keyMask = uint16(0x03FF)
)

// KeyPressed returns true if the given key is pressed.
func KeyPressed(k KeyCode) bool {
	return (^input.Get()&keyMask)&uint16(k) > 0
}

// empty/reset canvas
func clearScreen() {
	tinydraw.FilledRectangle(
		display,
		int16(0), int16(0),
		screenWidth, screenHeight,
		red,
	)

	tinydraw.FilledRectangle(
		display,
		border, border,
		screenWidth-(2*border), screenHeight-(2*border),
		white,
	)

	tinydraw.FilledCircle(
		display,
		10, 150,
		7,
		white,
	)

	tinydraw.FilledCircle(
		display,
		230, 150,
		7,
		white,
	)

	tinyfont.WriteLine(display, &tinyfont.TomThumb, 90, border/2, "Etchy-Sketchy "+version, white)
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 6, 153, "< >", black)
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 229, 149, "A", black)
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 229, 157, "B", black)
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 97, 155, "RESET (START)", white)
}

// start screen showing github url and build version
func startScreen() {
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 5, 155, version, white)
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 145, 155, "github.com/Shellywell123", white)
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 5, 10, "GAME BOY Advance - TinyGO", white)
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 220, 10, "2023", white)
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 90, 80, "Etchy-Sketchy", red)
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 85, 90, "Press START button", green)
}

// prevent drawing outside of canvas
func withinBorders(next_x, next_y int16) bool {
	if (0+border < next_x) && (next_x < 240-border) &&
		(0+border < next_y) && (next_y < 160-border) {
		return true
	}
	return false
}

func update(interrupt.Interrupt) {

	// check for button inputs
	switch {

	// Start the "game"
	case (KeyPressed(KeyStart)):
		clearScreen()
		tinyfont.DrawChar(display, &tinyfont.TomThumb, x, y, '.', black)

	case (KeyPressed(KeyB) && KeyPressed(KeyRight)):
		if withinBorders(x+1, y-1) {
			x = x + 1
			y = y - 1
			tinyfont.DrawChar(display, &tinyfont.TomThumb, x, y, '.', black)
		}

	case (KeyPressed(KeyB) && KeyPressed(KeyLeft)):
		if withinBorders(x-1, y-1) {
			x = x - 1
			y = y - 1
			tinyfont.DrawChar(display, &tinyfont.TomThumb, x, y, '.', black)
		}

	case (KeyPressed(KeyA) && KeyPressed(KeyRight)):
		if withinBorders(x+1, y+1) {
			x = x + 1
			y = y + 1
			tinyfont.DrawChar(display, &tinyfont.TomThumb, x, y, '.', black)
		}

	case (KeyPressed(KeyA) && KeyPressed(KeyLeft)):
		if withinBorders(x-1, y+1) {
			x = x - 1
			y = y + 1
			tinyfont.DrawChar(display, &tinyfont.TomThumb, x, y, '.', black)
		}

	case (KeyPressed(KeyRight)):
		if withinBorders(x+1, y) {
			x = x + 1
			tinyfont.DrawChar(display, &tinyfont.TomThumb, x, y, '.', black)
		}

	case (KeyPressed(KeyLeft)):
		if withinBorders(x-1, y) {
			x = x - 1
			tinyfont.DrawChar(display, &tinyfont.TomThumb, x, y, '.', black)
		}

	case (KeyPressed(KeyA)):
		if withinBorders(x, y+1) {
			y = y + 1
			tinyfont.DrawChar(display, &tinyfont.TomThumb, x, y, '.', black)
		}

	case (KeyPressed(KeyB)):
		if withinBorders(x, y-1) {
			y = y - 1
			tinyfont.DrawChar(display, &tinyfont.TomThumb, x, y, '.', black)
		}
	}
}

func main() {
	// Set up the display
	display.Configure()

	// Register display status
	regDISPSTAT.SetBits(1<<3 | 1<<4)

	// render start screen
	startScreen()

	// Creates an interrupt that will call the "update" fonction below, hardware way to display things on the screen
	interrupt.New(machine.IRQ_VBLANK, update).Enable()

	// Infinite loop to avoid exiting the application
	for {
	}
}
