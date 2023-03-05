package addons

import (
	"image/color"
)

type Pix struct {
	X uint8
	Y uint8
}

/*
var BL = color.RGBA{R: 0, G: 0, B: 0, A: 255}
var TL = color.RGBA{R: 0, G: 0, B: 0, A: 255}
var BR = color.RGBA{R: 0, G: 0, B: 50, A: 255}
var TR = color.RGBA{R: 0, G: 50, B: 0, A: 255}
*/

var Pixels = []Pix{
	// A
	{X: 14, Y: 1},
	{X: 13, Y: 2},
	{X: 12, Y: 3},
	{X: 11, Y: 4},
	{X: 10, Y: 5},
	{X: 9, Y: 6},
	{X: 8, Y: 7},
	{X: 7, Y: 8},
	{X: 6, Y: 9},
	{X: 5, Y: 10},
	{X: 4, Y: 11},
	{X: 3, Y: 12},
	{X: 2, Y: 13},
	{X: 1, Y: 14},
	// B
	{X: 0, Y: 16},
	{X: 0, Y: 18},
	{X: 0, Y: 20},
	{X: 0, Y: 22},
	{X: 0, Y: 24},
	{X: 0, Y: 26},
	{X: 0, Y: 28},
	{X: 0, Y: 30},
	{X: 0, Y: 32},
	{X: 0, Y: 34},
	{X: 0, Y: 36},
	{X: 0, Y: 38},
	{X: 0, Y: 40},
	{X: 0, Y: 42},
	// C
	{X: 1, Y: 42},
	{X: 2, Y: 41},
	{X: 3, Y: 40},
	{X: 4, Y: 39},
	{X: 5, Y: 38},
	{X: 6, Y: 37},
	{X: 7, Y: 36},
	{X: 8, Y: 35},
	{X: 9, Y: 34},
	{X: 10, Y: 33},
	{X: 11, Y: 32},
	{X: 12, Y: 31},
	{X: 13, Y: 30},
	{X: 14, Y: 29},
	// D
	{X: 15, Y: 27},
	{X: 15, Y: 25},
	{X: 15, Y: 23},
	{X: 15, Y: 21},
	{X: 15, Y: 19},
	{X: 15, Y: 17},
	{X: 15, Y: 15},
	{X: 15, Y: 13},
	{X: 15, Y: 11},
	{X: 15, Y: 9},
	{X: 15, Y: 7},
	{X: 15, Y: 5},
	{X: 15, Y: 3},
	{X: 15, Y: 1},
	// E
	{X: 16, Y: 1},
	{X: 17, Y: 2},
	{X: 18, Y: 3},
	{X: 19, Y: 4},
	{X: 20, Y: 5},
	{X: 21, Y: 6},
	{X: 22, Y: 7},
	{X: 23, Y: 8},
	{X: 24, Y: 9},
	{X: 25, Y: 10},
	{X: 26, Y: 11},
	{X: 27, Y: 12},
	{X: 28, Y: 13},
	{X: 29, Y: 14},
	// F
	{X: 30, Y: 16},
	{X: 30, Y: 17},
	{X: 30, Y: 18},
	{X: 30, Y: 19},
	{X: 30, Y: 20},
	{X: 30, Y: 21},
	{X: 30, Y: 22},
	{X: 30, Y: 23},
	{X: 30, Y: 24},
	{X: 30, Y: 25},
	{X: 30, Y: 26},
	{X: 30, Y: 27},
	{X: 30, Y: 28},
	{X: 30, Y: 29},
	// G
	{X: 29, Y: 44},
	{X: 28, Y: 45},
	{X: 27, Y: 46},
	{X: 26, Y: 47},
	{X: 25, Y: 48},
	{X: 24, Y: 49},
	{X: 23, Y: 50},
	{X: 22, Y: 51},
	{X: 21, Y: 52},
	{X: 20, Y: 53},
	{X: 19, Y: 54},
	{X: 18, Y: 55},
	{X: 17, Y: 56},
	{X: 16, Y: 57},
	// H
	{X: 14, Y: 57},
	{X: 13, Y: 56},
	{X: 12, Y: 55},
	{X: 11, Y: 54},
	{X: 10, Y: 53},
	{X: 9, Y: 52},
	{X: 8, Y: 51},
	{X: 7, Y: 50},
	{X: 6, Y: 49},
	{X: 5, Y: 48},
	{X: 4, Y: 47},
	{X: 3, Y: 46},
	{X: 2, Y: 45},
	{X: 1, Y: 44},
	// I
	{X: 16, Y: 29},
	{X: 17, Y: 30},
	{X: 18, Y: 31},
	{X: 19, Y: 32},
	{X: 20, Y: 33},
	{X: 21, Y: 34},
	{X: 22, Y: 35},
	{X: 23, Y: 36},
	{X: 24, Y: 37},
	{X: 25, Y: 38},
	{X: 26, Y: 39},
	{X: 27, Y: 40},
	{X: 28, Y: 41},
	{X: 29, Y: 42},
	// J
	{X: 31, Y: 42},
	{X: 32, Y: 41},
	{X: 33, Y: 40},
	{X: 34, Y: 39},
	{X: 35, Y: 38},
	{X: 36, Y: 37},
	{X: 37, Y: 36},
	{X: 38, Y: 35},
	{X: 39, Y: 34},
	{X: 40, Y: 33},
	{X: 41, Y: 32},
	{X: 42, Y: 31},
	{X: 43, Y: 30},
	{X: 44, Y: 29},
	// K
	{X: 45, Y: 27},
	{X: 45, Y: 25},
	{X: 45, Y: 23},
	{X: 45, Y: 21},
	{X: 45, Y: 19},
	{X: 45, Y: 17},
	{X: 45, Y: 15},
	{X: 45, Y: 13},
	{X: 45, Y: 11},
	{X: 45, Y: 9},
	{X: 45, Y: 7},
	{X: 45, Y: 5},
	{X: 45, Y: 3},
	{X: 45, Y: 1},
	// L
	{X: 46, Y: 1},
	{X: 47, Y: 2},
	{X: 48, Y: 3},
	{X: 49, Y: 4},
	{X: 50, Y: 5},
	{X: 51, Y: 6},
	{X: 52, Y: 7},
	{X: 53, Y: 8},
	{X: 54, Y: 9},
	{X: 55, Y: 10},
	{X: 56, Y: 11},
	{X: 57, Y: 12},
	{X: 58, Y: 13},
	{X: 59, Y: 14},
	// M
	{X: 60, Y: 16},
	{X: 60, Y: 18},
	{X: 60, Y: 20},
	{X: 60, Y: 22},
	{X: 60, Y: 24},
	{X: 60, Y: 26},
	{X: 60, Y: 28},
	{X: 60, Y: 30},
	{X: 60, Y: 32},
	{X: 60, Y: 34},
	{X: 60, Y: 36},
	{X: 60, Y: 38},
	{X: 60, Y: 40},
	{X: 60, Y: 42},
	// N
	{X: 59, Y: 42},
	{X: 58, Y: 41},
	{X: 57, Y: 40},
	{X: 56, Y: 39},
	{X: 55, Y: 38},
	{X: 54, Y: 37},
	{X: 53, Y: 36},
	{X: 52, Y: 35},
	{X: 51, Y: 34},
	{X: 50, Y: 33},
	{X: 49, Y: 32},
	{X: 48, Y: 31},
	{X: 47, Y: 30},
	{X: 46, Y: 29},
	// O
	{X: 44, Y: 1},
	{X: 43, Y: 2},
	{X: 42, Y: 3},
	{X: 41, Y: 4},
	{X: 40, Y: 5},
	{X: 39, Y: 6},
	{X: 38, Y: 7},
	{X: 37, Y: 8},
	{X: 36, Y: 9},
	{X: 35, Y: 10},
	{X: 34, Y: 11},
	{X: 33, Y: 12},
	{X: 32, Y: 13},
	{X: 31, Y: 14},
	// P
	{X: 31, Y: 44},
	{X: 32, Y: 45},
	{X: 33, Y: 46},
	{X: 34, Y: 47},
	{X: 35, Y: 48},
	{X: 36, Y: 49},
	{X: 37, Y: 50},
	{X: 38, Y: 51},
	{X: 39, Y: 52},
	{X: 40, Y: 53},
	{X: 41, Y: 54},
	{X: 42, Y: 55},
	{X: 43, Y: 56},
	{X: 44, Y: 57},
	// Q
	{X: 46, Y: 57},
	{X: 47, Y: 56},
	{X: 48, Y: 55},
	{X: 49, Y: 54},
	{X: 50, Y: 53},
	{X: 51, Y: 52},
	{X: 52, Y: 51},
	{X: 53, Y: 50},
	{X: 54, Y: 49},
	{X: 55, Y: 48},
	{X: 56, Y: 47},
	{X: 57, Y: 46},
	{X: 58, Y: 45},
	{X: 59, Y: 44},
}

var MaxX uint8 = 60
var MaxY uint8 = 58

func Fade(corners LEDWallFade) (retval []color.RGBA) {
	for _, pix := range Pixels {
		// Figure out the red value up the left side, then the right side, then
		// look at how far between those sides we are to get the final value
		red := scaleToRange2d(pix, corners.BL.R, corners.BR.R, corners.TL.R, corners.TR.R)
		green := scaleToRange2d(pix, corners.BL.G, corners.BR.G, corners.TL.G, corners.TR.G)
		blue := scaleToRange2d(pix, corners.BL.B, corners.BR.B, corners.TL.B, corners.TR.B)
		retval = append(retval, color.RGBA{R: red, G: green, B: blue, A: 255})
		//fmt.Println(pixnum, pix, red, green, blue)
		//fmt.Printf("/opt/homebrew/bin/mosquitto_pub -h 10.1.0.1 -t '/ledwall/1/request' -m \"{'action':'pixel','num':%d,'r':%d,'g':%d,'b':%d}\"\n", pixnum, red, green, blue)
	}
	return
}

func scaleToRange2d(pix Pix, bl, br, tl, tr uint8) uint8 {
	Left := uint8(scaleToRange(pix.Y, 0, MaxY, bl, tl))
	Right := uint8(scaleToRange(pix.Y, 0, MaxY, br, tr))
	val := uint8(scaleToRange(pix.X, 0, MaxX, Left, Right))
	return val
}

func scaleToRange(num, numMin, numMax, rangeMin, rangeMax uint8) float32 {
	// How far through the original range is num?
	var fraction = float32(int8(num)-int8(numMin)) / float32(int8(numMax)-int8(numMin))
	// How far through the output range is that?
	result := (float32(int8(rangeMax)-int8(rangeMin)) * fraction) + float32(rangeMin)
	return result
}
