package service

import (
	"fmt"
	"github.com/fishtailstudio/imgo"
	"github.com/sadlil/gologger"
	"imgv/utils"
	"math"
	"strconv"
)

var logger = gologger.GetLogger()

func lfit(srcW, srcH, dstW, dstH int) (int, int) {
	if float64(srcW)/float64(dstW) == float64(srcH)/float64(dstH) {
		return dstW, dstH
	}
	if float64(srcW)/float64(dstW) > float64(srcH)/float64(dstH) {
		h := float64(srcH) / (float64(srcW)/float64(dstW))
		return dstW, int(math.Round(h))
	} else {
		w := float64(srcW) / (float64(srcH)/float64(dstH))
		return int(math.Round(w)), dstH
	}
}

func mfit(srcW, srcH, dstW, dstH int) (int, int) {
	if float64(srcW)/float64(dstW) == float64(srcH)/float64(dstH) {
		return dstW, dstH
	}
	if float64(srcW)/float64(dstW) < float64(srcH)/float64(dstH) {
		h := float64(srcH) / (float64(srcW)/float64(dstW))
		return dstW, int(math.Round(h))
	} else {
		w := float64(srcW) / (float64(srcH)/float64(dstH))
		return int(math.Round(w)), dstH
	}
}

func Resize(imgUrl string, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("Resize image url: " + imgUrl)
	logger.Debug("params: " + utils.ToJSON(params))
	img := imgo.LoadFromUrl(imgUrl)
	width := img.Width()
	height := img.Height()
	contentType := img.Mimetype()
	m, ws, hs := params["m"], params["w"], params["h"]
	if m == "" {
		m = "lfit"
	}
	if ws == "" {
		ws = fmt.Sprintf("%d", width)
	}
	w, _ := strconv.Atoi(ws)
	if hs == "" {
		hs = fmt.Sprintf("%d", height)
	}
	h, _ := strconv.Atoi(hs)
	switch m {
	case "lfit":
		w, h = lfit(width, height, w, h)
		logger.Debug(fmt.Sprintf("lfit: width=%d, height=%d, w=%d, h=%d", width, height, w, h))
		img = img.Resize(w, h)
	case "mfit":
		w, h = mfit(width, height, w, h)
		logger.Debug(fmt.Sprintf("mfit: width=%d, height=%d, w=%d, h=%d", width, height, w, h))
		img = img.Resize(w, h)
	case "fill":
		wr, hr := mfit(width, height, w, h)
		logger.Debug(fmt.Sprintf("fill-mfit: width=%d, height=%d, w=%d, h=%d", width, height, wr, hr))
		logger.Debug(fmt.Sprintf("corp: x=%d y=%d",(wr-w)/2, (hr-h)/2))
		img = img.Resize(wr, hr)
		img = img.Crop((wr-w)/2, (hr-h)/2, w, h)
	case "pad":

	case "fixed":
		img = img.Resize(w, h)
	}
	return contentType, img, nil
}
