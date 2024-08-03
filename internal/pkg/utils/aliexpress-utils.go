package utils

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	_ "github.com/ahmetb/go-linq/v3"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
	"html"
	"math/rand"
	"strings"
	"time"
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
	// Seed the random number generator to ensure different results each run
	rand.Seed(time.Now().UnixNano())

	// Define the range of characters (a-z, A-Z, 0-9)
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a byte slice of length n
	b := make([]byte, n)

	// Populate the byte slice with random characters
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	// Return the random string
	return string(b)
}

func NewLink(link string) string {
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
		newUrl := NewLink(img.Url)
		imgs = append(imgs, newUrl)
	}

	i := 0
	for len(imgs) < 7 {
		imgs = append(imgs, imgs[i])
		i++
	}

	return imgs
}
