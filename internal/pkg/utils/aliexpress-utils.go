package utils

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	_ "github.com/ahmetb/go-linq/v3"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
	"html"
	"math/rand"
	"strings"
)

func ConvertSpecificationsToHTML(specs []model.Specification) string {
	var sb strings.Builder

	// Write the opening <ul> tag
	sb.WriteString("<ul>\n")

	// Iterate over each specification and add it as a list item
	for _, spec := range specs {
		sb.WriteString(fmt.Sprintf("  <li><strong>%s:</strong> %s</li>\n",
			html.EscapeString(spec.Name),
			html.EscapeString(spec.Value)))
	}

	// Write the closing </ul> tag
	sb.WriteString("</ul>")

	return sb.String()
}

func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func NewLink(link string) string {
	if len(link) < 4 {
		return link
	}

	return link[:len(link)-4] + "/" + RandomString(13) + ".jpg"
}

func GetSortedImages(images []model.Image) []string {
	var imgs []string
	linq.From(images).
		OrderBy(func(item interface{}) interface{} {
			return item.(model.Image).ZIndex
		}).
		ToSlice(&images)
	for _, img := range images {
		imgs = append(imgs, img.Url)
	}

	i := 0
	for len(imgs) < 7 {
		imgs = append(imgs, imgs[i])
		i++
	}

	return imgs
}

func Round2Decimals(f float64) float64 {
	return float64(int(f*100)) / 100
}
