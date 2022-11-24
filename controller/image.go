package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sadlil/gologger"
	"image"
	"imgv/service"
	"imgv/utils"
	"net/url"
)

var logger = gologger.GetLogger()

func ImageProcess(c *gin.Context) (string, image.Image, error) {
	oriUrl := c.Query("url")
	if oriUrl == "" {
		logger.Error("params missing url")
		return "", nil, fmt.Errorf("params missing url")
	}
	logger.Debug("url = " + oriUrl)
	u, err := url.Parse(oriUrl)
	if err != nil {
		logger.Error("url parse error: " + err.Error())
		return "", nil, err
	}
	logger.Debug("u = " + utils.ToJSON(u))
	imgProcessParams := c.Query("x-oss-process")
	values := u.Query()
	if imgProcessParams == "" {
		imgProcessParams = values.Get("x-oss-process")
		if imgProcessParams == "" {
			logger.Error("params missing x-oss-process")
			return "", nil, fmt.Errorf("params missing x-oss-process")
		}
	}
	actions, actionParams, err := service.SplitImageProcessParameters(imgProcessParams)
	if err != nil {
		return "", nil, err
	}
	values.Del("x-oss-process")
	for k, v := range c.Request.URL.Query() {
		if k != "url" && k != "x-oss-process" {
			logger.Debug(fmt.Sprintf("添加query参数: %s: %s", k, v[0]))
			values.Add(k, v[0])
		}
	}
	u.RawQuery = values.Encode()
	imgUrl := u.String()
	logger.Debug("图片源地址: " + imgUrl)
	var contentType string
	img := service.LoadImage(imgUrl)
	for i, action := range actions {
		switch action {
		case "resize":
			contentType, img, err = service.Resize(img, actionParams[i])
		case "corp":
			contentType, img, err = service.Corp(img, actionParams[i])
		case "rotate":
			contentType, img, err = service.Rotate(img, actionParams[i])
		case "format":

		case "rounded-corners":
			contentType, img, err = service.RoundedCorners(img, actionParams[i])
		case "circle":
			contentType, img, err = service.Circle(img, actionParams[i])
		case "blur":
			contentType, img, err = service.Blur(img, actionParams[i])
		case "sharpen":

		case "bright":

		case "contrast":

		case "info":

		}

	}
	if err != nil {
		return "", nil, err
	}
	return contentType, img.ToImage(), nil
}
