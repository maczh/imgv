package service

import (
	"github.com/fishtailstudio/imgo"
	"imgv/utils"
)

var format = map[string]string{
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"bmp":  "image/bmp",
	"gif":  "image/gif",
	"tiff": "image/tiff",
}

func Format(img *imgo.Image, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("bright params: " + utils.ToJSON(params))
	contentType := format[params["value"]]
	if contentType == "" {
		contentType = "image/png"
	}
	return contentType, img, nil
}
