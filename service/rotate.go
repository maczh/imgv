package service

import (
	"github.com/fishtailstudio/imgo"
	"imgv/utils"
	"strconv"
)

func RotateUrl(imgUrl string, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("Rotate image url: " + imgUrl)
	img := LoadImage(imgUrl)
	if img.Error != nil {
		return "", img, img.Error
	}
	return Rotate(img, params)
}

func Rotate(img *imgo.Image, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("rotate params: " + utils.ToJSON(params))
	contentType := img.Mimetype()
	v := params["value"]
	degree, _ := strconv.Atoi(v)
	return contentType, img.Rotate(degree), nil
}
