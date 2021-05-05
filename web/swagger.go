package web

import (
	"embed"

	"github.com/zopsmart/gofr/pkg/gofr/template"
)

//go:embed swagger/*
var fs embed.FS //nolint:gochecknoglobals // This has to be declared as global as per embed package implementation

func GetSwaggerFile(fileName string) (data []byte, contentType string, err error) {
	t := template.Template{}
	if fileName == "" {
		t.File = "index.html"
		t.Type = template.HTML
	} else {
		t.File = fileName
		t.Type = template.FILE
	}

	data, err = fs.ReadFile("swagger/" + t.File)
	if err != nil {
		return
	}

	contentType = t.ContentType()

	return
}
