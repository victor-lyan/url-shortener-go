package server

import (
	"embed"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

//go:embed static/*
var content embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name+".html", data)
}

func NewRenderer() *Template {
	temps, _ := template.ParseFS(content, "static/*.html")
	return &Template{templates: temps}
}
