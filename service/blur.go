package service

import (
	"github.com/fishtailstudio/imgo"
	"imgv/utils"
	"strconv"
)

func Blur(img *imgo.Image, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("blur params: " + utils.ToJSON(params))
	contentType := img.Mimetype()
	r, s := params["r"], params["s"]
	if r == "" {
		r = "1"
	}
	if s == "" {
		s = "1"
	}
	radius, _ := strconv.Atoi(r)
	sigma, _ := strconv.Atoi(s)
	return contentType, img.GaussianBlur(radius, float64(sigma)), nil
}
