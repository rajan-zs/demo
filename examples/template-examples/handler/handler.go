package handler

import (
	"bytes"
	"image"
	"image/png"
	"os"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/template"
)

func Template(ctx *gofr.Context) (interface{}, error) {
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

	return template.Template{Directory: ctx.TemplateDir, File: "test.html", Data: data, Type: template.HTML}, nil
}

// Image handler demonstrates how to use `template.File` for responding with any Content-Type,
// in this example we respond with a PNG image
func Image(ctx *gofr.Context) (interface{}, error) {
	f, _ := os.Open(ctx.TemplateDir + "/gopher.png")

	defer f.Close()

	i, _, _ := image.Decode(f)

	b := new(bytes.Buffer)

	err := png.Encode(b, i)
	if err != nil {
		return nil, err
	}

	return template.File{
		Content:     b.Bytes(),
		ContentType: "image/png",
	}, nil
}
