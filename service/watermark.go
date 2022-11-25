package service

import (
	"encoding/base64"
	"fmt"
	"github.com/fishtailstudio/imgo"
	"image"
	"image/color"
	"image/draw"
	"imgv/utils"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"unicode"
)

var fonts = map[string]string{
	"方正仿宋":   "/fonts/FZFSK.ttf",
	"华文宋体":   "/fonts/Songti.ttc",
	"方正书宋":   "/fonts/FZSSK.ttf",
	"方正楷体":   "/fonts/FZKTK.ttf",
	"文泉驿正黑":  "/fonts/文泉驿正黑.ttf",
	"文泉驿微米黑": "/fonts/文泉驿微米黑.ttf",
}

type Position struct {
	Grid string //九宫格位置
	Dx   int    //距离左(第左、中列)、右(右列)边距点数
	Dy   int    //距离上(第上、中行)、下(下行)边距点数
}

type markText struct {
	Text     string //文字内容
	FontName string //字体名称
	Color    string //文字颜色RGB,例如：000000表示黑色，FFFFFF表示白色
	Size     int    //字体大小 单位pt
	Alpha    uint32 //透明度
	Position
}

func (t *markText) FontFile() string {
	p, _ := filepath.Abs(path.Dir(os.Args[0]))
	fontFile := p + fonts[t.FontName]
	logger.Debug("字体文件名: " + fontFile)
	return fontFile
}

func (t *markText) length() float64 {
	var l float64 = 0.0
	for _, r := range []rune(t.Text) {
		if unicode.Is(unicode.Han, r) {
			l++
		} else {
			l += 0.5
		}
	}
	return l
}

func (t *markText) Width() int {
	logger.Debug(fmt.Sprintf("文字个数: %d", t.length()))
	w := int((t.length() * float64(t.Size)) / 0.75)
	logger.Debug(fmt.Sprintf("文字宽度: %d px", w))
	return w
}

func (t *markText) Height() int {
	h := int(float64(t.Size) / 0.75)
	logger.Debug(fmt.Sprintf("文字高度: %d px", h))
	return h
}

func (t *markText) TextColor() color.Color {
	r, _ := strconv.ParseUint(t.Color[:2], 16, 8)
	g, _ := strconv.ParseUint(t.Color[2:4], 16, 8)
	b, _ := strconv.ParseUint(t.Color[4:], 16, 8)
	a := uint8(float64(t.Alpha) * 2.55)
	return color.NRGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: a,
	}
}

type markImage struct {
	Url   string  //水印图片地址
	Scale float64 //缩放比例 1-100
	Alpha float64 //透明度 0-100
	Position
}

func calcWaterMarkLeftTop(pos Position, width, height, w, h int) (int, int) {
	//blockWidth := int(math.Round(float64(width) / float64(3)))
	//blockHeight := int(math.Round(float64(height) / float64(3)))
	switch pos.Grid {
	case "nw":
		return pos.Dx, pos.Dy
	case "north":
		return pos.Dx + (width-w)/2, pos.Dy
	case "ne":
		return width - w - pos.Dx, pos.Dy
	case "west":
		return pos.Dx, pos.Dy + (height-h)/2
	case "center":
		return pos.Dx + (width-w)/2, pos.Dy + (height-h)/2
	case "east":
		return width - w - pos.Dx, pos.Dy + (height-h)/2
	case "sw":
		return pos.Dx, height - h - pos.Dy
	case "south":
		return pos.Dx + (width-w)/2, height - h - pos.Dy
	case "se":
		return width - w - pos.Dx, height - h - pos.Dy
	default:
		return width - w - pos.Dx, height - h - pos.Dy
	}
}

func WaterMark(img *imgo.Image, params map[string]string) (string, *imgo.Image, error) {
	logger.Debug("watermark params: " + utils.ToJSON(params))
	//txt,text := params["text"],""
	//if txt != "" {
	//	t, _ := base64.StdEncoding.DecodeString(txt)
	//	text = string(t)
	//}
	text := params["text"]
	font := params["type"]
	if font == "" {
		font = "文泉驿正黑"
	}
	s := params["size"]
	if s == "" {
		s = "40"
	}
	size, _ := strconv.Atoi(s)
	size = int(float64(size) * 0.75) //将px转成pt
	c := params["color"]
	if c == "" {
		c = "000000"
	}
	host := params["host"]
	ts, t := params["t"], 100
	if ts != "" {
		t, _ = strconv.Atoi(ts)
	}
	xx, yy, x, y := params["x"], params["y"], 0, 0
	x, _ = strconv.Atoi(xx)
	y, _ = strconv.Atoi(yy)
	g := params["g"]
	if g == "" {
		g = "se"
	}
	p := params["P"]
	if p == "" {
		p = "100"
	}
	scale, _ := strconv.Atoi(p)
	u, imgUrl := params["image"], ""
	if u != "" {
		iuu, _ := url.QueryUnescape(u)
		iu, _ := base64.StdEncoding.DecodeString(iuu)
		imgUrl = string(iu)
		if imgUrl[:1] == "/" {
			imgUrl = host + imgUrl
		}
	}
	logger.Debug("text: " + text)
	logger.Debug("imgUrl: " + imgUrl)
	if text != "" {
		txtMark := markText{
			Text:     text,
			FontName: font,
			Color:    c,
			Size:     size,
			Alpha:    uint32(t),
			Position: Position{
				Grid: g,
				Dx:   x,
				Dy:   y,
			},
		}
		return img.Mimetype(), textWaterMark(img, txtMark), nil
	}
	if imgUrl != "" {
		imageMark := markImage{
			Url:   imgUrl,
			Scale: float64(scale),
			Alpha: float64(t),
			Position: Position{
				Grid: g,
				Dx:   x,
				Dy:   y,
			},
		}
		return img.Mimetype(), imageWaterMark(img, imageMark), nil
	}
	return "", nil, fmt.Errorf("不是文字或图片类型水印")
}

func textWaterMark(img *imgo.Image, text markText) *imgo.Image {
	logger.Debug(fmt.Sprintf("文字宽度: %d, 高度: %d", text.Width(), text.Height()))
	x, y := calcWaterMarkLeftTop(text.Position, img.Width(), img.Height(), text.Width(), text.Height())
	logger.Debug(fmt.Sprintf("文字水印左上角: x=%d, y=%d", x, y))
	return img.Text(text.Text, x, y, text.FontFile(), text.TextColor(), float64(text.Size), 96)
}

//图片水印
func imageWaterMark(img *imgo.Image, mark markImage) *imgo.Image {
	markImg := imgo.LoadFromUrl(mark.Url)
	if markImg.Error != nil {
		logger.Error(markImg.Error.Error())
		return img
	}
	markImg = markImg.Resize(int(float64(markImg.Width())*mark.Scale/float64(100)), int(float64(markImg.Height())*mark.Scale/float64(100)))
	mi := adjustOpacity(toRGBA64(markImg.ToImage()), mark.Alpha/float64(100))
	x, y := calcWaterMarkLeftTop(mark.Position, img.Width(), img.Height(), markImg.Width(), markImg.Height())
	logger.Debug(fmt.Sprintf("水印左上角位置: x=%d,y=%d", x, y))
	return img.Insert(imgo.LoadFromImage(RGBA64toRGBA(mi)), x, y)
}

//Image转换为image.RGBA64
func toRGBA64(m image.Image) *image.RGBA64 {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA64(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, a := colorRgb.RGBA()
			nR := uint16(r)
			nG := uint16(g)
			nB := uint16(b)
			alpha := uint16(a)
			newRgba.SetRGBA64(i, j, color.RGBA64{R: nR, G: nG, B: nB, A: alpha})
		}
	}
	return newRgba
}

//将输入图像m的透明度变为原来的倍数。若原来为完成全不透明，则percentage = 0.5将变为半透明
func adjustOpacity(m *image.RGBA64, percentage float64) *image.RGBA64 {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA64(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, a := colorRgb.RGBA()
			opacity := uint16(float64(a) * percentage)
			//颜色模型转换，至关重要！
			v := newRgba.ColorModel().Convert(color.NRGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: opacity})
			//Alpha = 0: Full transparent
			rr, gg, bb, aa := v.RGBA()
			newRgba.SetRGBA64(i, j, color.RGBA64{R: uint16(rr), G: uint16(gg), B: uint16(bb), A: uint16(aa)})
		}
	}
	return newRgba
}

func RGBA64toRGBA(img *image.RGBA64) *image.RGBA {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Over)
	return rgba
}
