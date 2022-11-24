package service

import (
	"github.com/disintegration/imaging"
	"github.com/fishtailstudio/imgo"
	"imgv/utils"
	"strconv"
)

func Sharpen(img *imgo.Image, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("sharpen params: " + utils.ToJSON(params))
	contentType := img.Mimetype()
	v := params["value"]
	if v == "" {
		v = "100"
	}
	s, _ := strconv.Atoi(v)
	im := imaging.Sharpen(img.ToImage(), float64(s)/float64(100))
	return contentType, imgo.LoadFromImage(im), nil
}
