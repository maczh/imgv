package service

import (
	"github.com/fishtailstudio/imgo"
	"image"
	"image/color"
	"image/draw"
	"imgv/utils"
	"strconv"
)

func RoundedCorners(img *imgo.Image, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("rounded corners params: " + utils.ToJSON(params))
	contentType := img.Mimetype()
	r := params["r"]
	radius, _ := strconv.Atoi(r)
	i := imgo.LoadFromImage(img.ToImage())
	return contentType, borderRadius(i, radius), nil
}

func borderRadius(img *imgo.Image, r int) *imgo.Image {
	if r > img.Width()/2 || r > img.Height()/2 {
		return img
	}
	w, h := img.Width(), img.Height()
	c := radius{p: image.Point{X: w, Y: h}, r: r}
	dst := image.NewRGBA(img.Bounds())
	logger.Debug("bounds: " + utils.ToJSON(dst.Bounds()))
	logger.Debug("radius: " + utils.ToJSON(c.p))
	draw.DrawMask(dst, dst.Bounds(), img.ToImage(), image.Point{}, &c, image.Point{}, draw.Over)

	img = imgo.LoadFromImage(dst)
	return img
}

type radius struct {
	p image.Point // 矩形右下角位置
	r int
}

func (c *radius) ColorModel() color.Model {
	return color.AlphaModel
}
func (c *radius) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.p.X, c.p.Y)
}

func (c *radius) At(x, y int) color.Color {
	var xx, yy, rr float64
	var inArea bool

	// left up
	if x <= c.r && y <= c.r {
		xx, yy, rr = float64(c.r-x)+0.5, float64(y-c.r)+0.5, float64(c.r)
		inArea = true
	}

	// right up
	if x >= (c.p.X-c.r) && y <= c.r {
		xx, yy, rr = float64(x-(c.p.X-c.r))+0.5, float64(y-c.r)+0.5, float64(c.r)
		inArea = true
	}

	// left bottom
	if x <= c.r && y >= (c.p.Y-c.r) {
		xx, yy, rr = float64(c.r-x)+0.5, float64(y-(c.p.Y-c.r))+0.5, float64(c.r)
		inArea = true
	}

	// right bottom
	if x >= (c.p.X-c.r) && y >= (c.p.Y-c.r) {
		xx, yy, rr = float64(x-(c.p.X-c.r))+0.5, float64(y-(c.p.Y-c.r))+0.5, float64(c.r)
		inArea = true
	}

	if inArea && xx*xx+yy*yy >= rr*rr {
		return color.Alpha{}
	}
	return color.Alpha{A: 255}
}
