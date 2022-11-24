package service

import (
	"crypto/md5"
	"fmt"
	"github.com/fishtailstudio/imgo"
	"io"
)

var CacheDir string

func MD5Encode(content string) (md string) {
	h := md5.New()
	_, _ = io.WriteString(h, content)
	md = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func LoadFromCache(imgUrl, process string) (*imgo.Image, error) {
	u := fmt.Sprintf("%s?x-oss-process=%s", imgUrl, process)
	m := MD5Encode(u)
	img := imgo.LoadFromPath(fmt.Sprintf("%s/%s.png", CacheDir, m))
	if img.Error != nil {
		return nil, img.Error
	} else {
		return img, nil
	}
}

func SaveCache(img *imgo.Image, imgUrl, process string) {
	u := fmt.Sprintf("%s?x-oss-process=%s", imgUrl, process)
	m := MD5Encode(u)
	img.Save(fmt.Sprintf("%s/%s.png", CacheDir, m))
}
