package controller

import (
	"fmt"
	"github.com/fishtailstudio/imgo"
	"github.com/gin-gonic/gin"
	"github.com/sadlil/gologger"
	"imgv/service"
	"imgv/utils"
	"net/url"
)

var logger = gologger.GetLogger()

func ImageProcess(c *gin.Context) (string, *imgo.Image, error) {
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
	if imgProcessParams == "" {
		logger.Error("params missing x-oss-process")
		return "", nil, fmt.Errorf("params missing x-oss-process")
	}
	action, actionParams, err := service.SplitImageProcessParameters(imgProcessParams)
	if err != nil {
		return "", nil, err
	}
	values := u.Query()
	for k,v := range c.Request.URL.Query(){
		if k != "url" && k != "x-oss-process" {
			logger.Debug(fmt.Sprintf("添加query参数: %s: %s",k, v[0]))
			values.Add(k, v[0])
		}
	}
	u.RawQuery = values.Encode()
	imgUrl := u.String()
	logger.Debug("图片源地址: " + imgUrl)
	switch action {
	case "resize":
		return service.Resize(imgUrl,actionParams)
	case "corp":

	case "rotate":

	case "format":

	case "rounded-corners":

	case "blur":

	case "sharpen":

	case "bright":

	case "contrast":

	case "info":

	}
	return "", nil, nil
}
