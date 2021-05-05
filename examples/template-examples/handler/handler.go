package handler

import (
	"bytes"
	"image"
	"image/png"
	"os"

	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/template"
)

func Template(c *gofr.Context) (interface{}, error) {
	data := struct {
		Title string
		Items []string
	}{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
			"My pages",
		},
	}

	return template.Template{Directory: c.TemplateDir, File: "test.html", Data: data, Type: template.HTML}, nil
}

// Image handler demonstrates how to use `template.File` for responding with any Content-Type,
// in this example we respond with a PNG image
func Image(c *gofr.Context) (interface{}, error) {
	f, _ := os.Open(c.TemplateDir + "/gopher.png")

	defer f.Close()

	i, _, _ := image.Decode(f)

	b := new(bytes.Buffer)
	png.Encode(b, i)

	return template.File{
		Content:     b.Bytes(),
		ContentType: "image/png",
	}, nil
}