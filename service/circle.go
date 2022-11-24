package service

import (
	"github.com/fishtailstudio/imgo"
	"image"
	"image/color"
	"image/draw"
	"imgv/utils"
	"strconv"
)

func Circle(img *imgo.Image, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("circle params: " + utils.ToJSON(params))
	contentType := img.Mimetype()
	r := params["r"]
	radius, _ := strconv.Atoi(r)
	return contentType, borderCircle(img, radius), nil
}

func borderCircle(img *imgo.Image, r int) *imgo.Image {
	if r > img.Width()/2 || r > img.Height()/2 {
		return img
	}
	w, h := img.Width(), img.Height()
	c := circle{p: image.Point{X: w / 2, Y: h / 2}, r: r}
	dst := image.NewRGBA(img.Bounds())
	draw.DrawMask(dst, img.Bounds(), img.ToImage(), image.Point{}, &c, image.Point{}, draw.Over)
	img = imgo.LoadFromImage(dst)
	x, y := w/2-r, h/2-r
	return img.Crop(x, y, r*2, r*2)
}

type circle struct { // 这里需要自己实现一个圆形遮罩，实现接口里的三个方法
	p image.Point // 圆心位置
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		//noinspection GoStructInitializationWithoutFieldNames
		return color.Alpha{255}
	}
	//noinspection GoStructInitializationWithoutFieldNames
	return color.Alpha{0}
}
