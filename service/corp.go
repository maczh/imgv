package service

import (
	"fmt"
	"github.com/fishtailstudio/imgo"
	"imgv/utils"
	"math"
	"strconv"
)

func topLeftPoint(x, y, width, height int, g string) (int, int) {
	blockWidth := int(math.Round(float64(width) / float64(3)))
	blockHeight := int(math.Round(float64(height) / float64(3)))
	switch g {
	case "nw":
		return x, y
	case "north":
		return x + blockWidth, y
	case "ne":
		return x + blockWidth*2, y
	case "west":
		return x, y + blockHeight
	case "center":
		return x + blockWidth, y + blockHeight
	case "east":
		return x + blockWidth*2, y + blockHeight
	case "sw":
		return x, y + blockHeight*2
	case "south":
		return x + blockWidth, y + blockHeight*2
	case "se":
		return x + blockWidth*2, y + blockHeight*2
	default:
		return x, y
	}
}

func CorpUrl(imgUrl string, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("Corp image url: " + imgUrl)
	img := LoadImage(imgUrl)
	if img.Error != nil {
		return "", img, img.Error
	}
	return Corp(img, params)
}

func Corp(img *imgo.Image, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("corp params: " + utils.ToJSON(params))
	width := img.Width()
	height := img.Height()
	contentType := img.Mimetype()
	ws, hs, xp, yp, g := params["w"], params["h"], params["x"], params["y"], params["g"]
	w, h, x, y := width, height, 0, 0
	if ws != "" {
		w, _ = strconv.Atoi(ws)
	}
	if hs != "" {
		h, _ = strconv.Atoi(hs)
	}
	if xp != "" {
		x, _ = strconv.Atoi(xp)
	}
	if yp != "" {
		y, _ = strconv.Atoi(yp)
	}
	if g == "" {
		g = "nw"
	}
	x, y = topLeftPoint(x, y, width, height, g)
	if w >= width {
		w = width - x
	}
	if h >= height {
		h = height - y
	}
	logger.Debug(fmt.Sprintf("Crop参数: x=%d y=%d w=%d h=%d", x, y, w, h))
	img = img.Crop(x, y, w, h)
	return contentType, img, nil
}
