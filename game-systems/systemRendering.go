package gamesystems

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var templates *template.Template

func mustLoadSystemTemplates(dir string) {
	wd, _ := os.Getwd()

	pageDir := filepath.Join(wd, dir)
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

func renderTemplateWithData(templateDir, name string, data interface{}) ([]byte, error) {
		if templates == nil {
		mustLoadSystemTemplates(templateDir)
	}

	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, name+".html", data)
	if err != nil {
		log.Println("template error", err)
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
