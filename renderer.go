package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
)


type Renderer struct {
	templates *template.Template
}

func NewRenderer() *Renderer {
	r := Renderer{}
	r.MustInitTemplates()
	return &r
}

func (r *Renderer) MustInitTemplates() {
	wd, _ := os.Getwd()
	d := filepath.Join(wd, "/templates")
	var allTemplateFiles []string
	err := filepath.WalkDir(d, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(p) == ".html" {
			allTemplateFiles = append(allTemplateFiles, p)
		}
		return nil
	})
	r.templates, err = templates.ParseFiles(allTemplateFiles...)
	if err != nil {
		log.Fatalf("Failed to parse template files %v\n", err)
	}
}

func (r Renderer) RenderPage(name string, data any) (string, error) {
	var buf bytes.Buffer
	err := r.templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Println("template error", err)
		return "", err
	}
	return buf.String(), nil
}

func (r Renderer) RenderErrorPage(msg string) (string, error) {
	epdata := struct {
		Msg string
	}{
		Msg: msg,
	}
	var buf bytes.Buffer
	err := r.templates.ExecuteTemplate(&buf, "error.html", epdata)
	if err != nil {
		log.Println("template error", err)
		return "", err
	}
	return buf.String(), nil
}
