package service

import (
	"bytes"
	"github.com/fishtailstudio/imgo"
	"github.com/golang/freetype"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"math"
)

type Position struct {
	Grid string //九宫格位置
	Dx   int    //距离左(第左、中列)、右(右列)边距点数
	Dy   int    //距离上(第上、中行)、下(下行)边距点数
}

type MarkText struct {
	Text      string    //文字内容
	FontFile  string    //字体路径
	Color     [4]uint8  //文字颜色RGBA 例子: [0,0,0,122]
	Size      vg.Length //字体大小 单位英寸
	Linewidth vg.Length //行高
	Angle     float64   //角度 0至1 的范围 即0到90度之间
	Space     string    //空格间隔
	Position
}

type MarkImage struct {
	Url   string  //水印图片地址
	Scale float64 //缩放比例 1-100
	Alpha float64 //透明度 0-100
	Position
}

func calcWaterMarkLeftTop(pos Position, width, height, w, h int) (int, int) {
	blockWidth := int(math.Round(float64(width) / float64(3)))
	blockHeight := int(math.Round(float64(height) / float64(3)))
	switch pos.Grid {
	case "nw":
		return pos.Dx, pos.Dy
	case "north":
		return pos.Dx + blockWidth, pos.Dy
	case "ne":
		return width - w - pos.Dx, pos.Dy
	case "west":
		return pos.Dx, pos.Dy + blockHeight
	case "center":
		return pos.Dx + blockWidth, pos.Dy + blockHeight
	case "east":
		return width - w - pos.Dx, pos.Dy + blockHeight
	case "sw":
		return pos.Dx, height - h - pos.Dy
	case "south":
		return pos.Dx + blockWidth, height - h - pos.Dy
	case "se":
		return width - w - pos.Dx, height - h - pos.Dy
	default:
		return width - w - pos.Dx, height - h - pos.Dy
	}
}

//图片水印
func imageWaterMark(img *imgo.Image, mark MarkImage) *imgo.Image {
	markImg := imgo.LoadFromUrl(mark.Url)
	if markImg.Error != nil {
		return img
	}
	markImg = markImg.Resize(int(float64(markImg.Width())*mark.Scale/float64(100)), int(float64(markImg.Height())*mark.Scale/float64(100)))
	mi := adjustOpacity(toRGBA64(markImg.ToImage()), mark.Alpha/float64(100))
	x, y := calcWaterMarkLeftTop(mark.Position, img.Width(), img.Height(), markImg.Width(), markImg.Height())
	return img.Insert(imgo.LoadFromImage(mi), x, y)
}

// WaterMark用于在图像上添加水印
func WaterMarkText(img image.Image, markText MarkText) (image.Image, error) {
	// 图片的长度设置画布的长度
	bounds := img.Bounds()
	w := vg.Length(bounds.Max.X) * vg.Inch / vgimg.DefaultDPI
	h := vg.Length(bounds.Max.Y) * vg.Inch / vgimg.DefaultDPI
	// 通过高和宽计算对角线
	diagonal := vg.Length(math.Sqrt(float64(w*w + h*h)))

	// 创建一个画布，宽度和高度是对角线
	c := vgimg.New(diagonal, diagonal)

	// 在画布中心绘制图像
	rect := vg.Rectangle{}
	// 计算中心位置,宽为w,高为h
	rect.Min.X = diagonal/2 - w/2
	rect.Min.Y = diagonal/2 - h/2
	rect.Max.X = diagonal/2 + w/2
	rect.Max.Y = diagonal/2 + h/2
	c.DrawImage(rect, img)

	// 制作一个 fontstyle ，宽度为英寸,字体 Courier 标准的等宽度字体
	// 读字体数据
	fontBytes, err := ioutil.ReadFile(markText.FontFile)
	if err != nil {
		log.Println("读取字体数据出错")
		log.Println(err)
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println("转换字体样式出错")
		log.Println(err)
		return nil, err
	}
	vg.AddFont("cn_font", font)
	fontStyle, err := vg.MakeFont("cn_font", vg.Inch*markText.Size)
	if err != nil {
		return nil, err
	}
	// 重复编写水印字体
	marktext := markText.Text
	unitText := marktext

	markTextWidth := fontStyle.Width(marktext)
	for markTextWidth <= diagonal {
		marktext += markText.Space + unitText
		markTextWidth = fontStyle.Width(marktext)
	}
	// 设置水印字体的颜色
	rgba := markText.Color
	c.SetColor(color.RGBA{rgba[0], rgba[1], rgba[2], rgba[3]})
	// 设置 0 到 π/2 之间的随机角度
	c.Rotate(markText.Angle * math.Pi / 2)

	// 设置每行水印的高度并添加水印
	// 一个字体的高度
	lineHeight := fontStyle.Extents().Height * markText.Linewidth
	for offset := -2 * diagonal; offset < 2*diagonal; offset += lineHeight {
		c.FillString(fontStyle, vg.Point{X: 0, Y: offset}, marktext)
	}

	// 画布写入新图片
	// 使用buffer去转换
	jc := vgimg.PngCanvas{Canvas: c}
	buff := new(bytes.Buffer)
	jc.WriteTo(buff)
	img, _, err = image.Decode(buff)
	if err != nil {
		return nil, err
	}

	// 得到图像的中心点
	ctp := int(diagonal * vgimg.DefaultDPI / vg.Inch / 2)

	// 切出打水印的图像
	size := bounds.Size()
	bounds = image.Rect(ctp-size.X/2, ctp-size.Y/2, ctp+size.X/2, ctp+size.Y/2)
	rv := image.NewRGBA(bounds)
	draw.Draw(rv, bounds, img, bounds.Min, draw.Src)
	return rv, nil
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
