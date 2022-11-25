package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sadlil/gologger"
	"image"
	"imgv/service"
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
	imgProcessParams := c.Query("x-oss-process")
	values := u.Query()
	if imgProcessParams == "" {
		imgProcessParams = values.Get("x-oss-process")
		if imgProcessParams == "" {
			logger.Error("params missing x-oss-process")
			return "", nil, fmt.Errorf("params missing x-oss-process")
		}
	}
	cacheImage, err := service.LoadFromCache(u.Host+u.Path, imgProcessParams)
	if err == nil && cacheImage != nil {
		return cacheImage.Mimetype(), cacheImage.ToImage(), nil
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
	if img.Error != nil {
		return "", nil, fmt.Errorf("404 Not found")
	}
	for i, action := range actions {
		switch action {
		case "resize":
			contentType, img, err = service.Resize(img, actionParams[i])
		case "corp":
			contentType, img, err = service.Corp(img, actionParams[i])
		case "rotate":
			contentType, img, err = service.Rotate(img, actionParams[i])
		case "format":
			contentType, img, err = service.Format(img, actionParams[i])
		case "watermark":
			actionParams[i]["host"] = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
			contentType, img, err = service.WaterMark(img, actionParams[i])
		case "rounded-corners":
			contentType, img, err = service.RoundedCorners(img, actionParams[i])
		case "circle":
			contentType, img, err = service.Circle(img, actionParams[i])
		case "blur":
			contentType, img, err = service.Blur(img, actionParams[i])
		case "sharpen":
			contentType, img, err = service.Sharpen(img, actionParams[i])
		case "bright":
			contentType, img, err = service.Bright(img, actionParams[i])
		case "contrast":
			contentType, img, err = service.Contrast(img, actionParams[i])
		case "info":

		}

	}
	if err != nil {
		return "", nil, err
	}
	service.SaveCache(img, u.Host+u.Path, imgProcessParams)
	return contentType, img.ToImage(), nil
}
