package main

import (
	"html/template"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

var funcMap = template.FuncMap{
	"inc": func(i int) int { return i + 1 },
	"dec": func(i int) int { return i - 1 },
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob((templatesDir + "/layouts/*.html"))
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob((templatesDir + "/includes/*.html"))
	if err != nil {
		panic(err.Error())
	}

	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFilesFuncs(filepath.Base(include), funcMap, files...)
	}

	return r
}

func templateData(c *gin.Context, data any) map[string]any {
	return map[string]any{
		"context":  c,
		"flashmsg": FlashMessage(c),
		"data":     data,
	}
}
