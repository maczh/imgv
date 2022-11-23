package service

import "github.com/fishtailstudio/imgo"

func LoadImage(imageUrl string) *imgo.Image {
	return imgo.LoadFromUrl(imageUrl)
}
