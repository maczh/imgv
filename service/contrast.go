package service

import (
	"github.com/disintegration/imaging"
	"github.com/fishtailstudio/imgo"
	"imgv/utils"
	"strconv"
)

func Contrast(img *imgo.Image, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("contrast params: " + utils.ToJSON(params))
	contentType := img.Mimetype()
	v := params["value"]
	if v == "" {
		v = "0"
	}
	s, _ := strconv.Atoi(v)
	im := imaging.AdjustContrast(img.ToImage(), float64(s))
	return contentType, imgo.LoadFromImage(im), nil
}
