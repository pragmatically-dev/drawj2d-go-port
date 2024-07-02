package remarkablepage

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func Hough(im image.Image, ntx, mry int) draw.Image {
	nimx := im.Bounds().Max.X
	mimy := im.Bounds().Max.Y

	him := image.NewGray(image.Rect(0, 0, ntx, mry))
	draw.Draw(him, him.Bounds(), image.NewUniform(color.White),
		image.Point{}, draw.Src)

	rmax := math.Hypot(float64(nimx), float64(mimy))
	dr := rmax / float64(mry/2)
	dth := math.Pi / float64(ntx)

	for jx := 0; jx < nimx; jx++ {
		for iy := 0; iy < mimy; iy++ {
			col := color.GrayModel.Convert(im.At(jx, iy)).(color.Gray)
			if col.Y == 255 {
				continue
			}
			for jtx := 0; jtx < ntx; jtx++ {
				th := dth * float64(jtx)
				r := float64(jx)*math.Cos(th) + float64(iy)*math.Sin(th)
				iry := mry/2 - int(math.Floor(r/dr+.5))
				col = him.At(jtx, iry).(color.Gray)
				if col.Y > 0 {
					col.Y--
					him.SetGray(jtx, iry, col)
				}
			}
		}
	}
	return him
}
