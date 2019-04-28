package driver

import (
	"github.com/HFO4/gbc-in-cloud/util"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"log"
	"os"
)

type LCD struct {
	pixels *[160][144][3]uint8
	window *pixelgl.Window

	pixelMap *pixel.PictureData

	inputStatus *byte
}

func (lcd *LCD) Init(pixels *[160][144][3]uint8) {
	lcd.pixels = pixels
	log.Println("[Display] Initialize GUI display")
	lcd.pixelMap = pixel.MakePictureData(pixel.R(0, 0, 160, 144))

}

func (lcd *LCD) InitStatus(statusPointer *byte) {
	lcd.inputStatus = statusPointer
}

func (lcd *LCD) UpdateInput() bool {
	// Mapping from keys to GB index.
	// ref :https://github.com/Humpheh/goboy/blob/master/pkg/gbio/iopixel/pixels.go
	var keyMap = map[pixelgl.Button]byte{
		// A button
		pixelgl.KeyZ: 4,
		// B button
		pixelgl.KeyX: 5,
		// SELECT button
		pixelgl.KeyBackspace: 6,
		// START button
		pixelgl.KeyEnter: 7,
		// RIGHT button
		pixelgl.KeyRight: 0,
		// LEFT button
		pixelgl.KeyLeft: 1,
		// UP button
		pixelgl.KeyUp: 2,
		// DOWN button
		pixelgl.KeyDown: 3,
	}
	var requestInterrupt bool
	var statusCopy byte
	statusCopy = *lcd.inputStatus
	for key, offset := range keyMap {
		if lcd.window.JustPressed(key) {
			statusCopy = util.ClearBit(statusCopy, uint(offset))
			requestInterrupt = true
		}
		if lcd.window.JustReleased(key) {
			statusCopy = util.SetBit(statusCopy, uint(offset))
			requestInterrupt = false
		}
	}

	*lcd.inputStatus = statusCopy
	return requestInterrupt
}

func (lcd *LCD) Run(drawSignal chan bool) {
	cfg := pixelgl.WindowConfig{
		Title:  "TETRIS [FPS:60] [CLOCK:4194304]",
		Bounds: pixel.R(0, 0, 160*3, 144*3),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	lcd.window = win

	go func() {
		for !win.Closed() {
		}
		os.Exit(1)
	}()

	for {
		<-drawSignal
		for y := 0; y < 144; y++ {
			for x := 0; x < 160; x++ {
				colour := color.RGBA{R: lcd.pixels[x][y][0], G: lcd.pixels[x][y][1], B: lcd.pixels[x][y][2], A: 0xFF}
				lcd.pixelMap.Pix[(143-y)*160+x] = colour
			}
		}

		graph := pixel.NewSprite(pixel.Picture(lcd.pixelMap), pixel.R(0, 0, 160, 144))
		mat := pixel.IM
		mat = mat.Moved(win.Bounds().Center())
		mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(3, 3))
		graph.Draw(lcd.window, mat)
		win.Update()
	}

}