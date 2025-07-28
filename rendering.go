package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
)


var templates *template.Template

func mustLoadTemplates() {
	wd, _ := os.Getwd()

	pageDir := filepath.Join(wd, "/templates/")
	var allFiles []string
	err := filepath.WalkDir(pageDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			allFiles = append(allFiles, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	templates = template.Must(template.ParseFiles(allFiles...))
}

func renderTemplateWithData(name string, data interface{}) ([]byte, error) {
		if templates == nil {
		mustLoadTemplates()
	}

	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, name+".html", data)
	if err != nil {
		log.Println("template error", err)
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func renderTemplate(name string) ([]byte, error) {
	if templates == nil {
		mustLoadTemplates()
	}

	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, name+".html", struct{}{})
	if err != nil {
		log.Println("template error", err)
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
