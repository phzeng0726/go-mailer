package service

import (
	"bytes"
	"log"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type TemplatesService struct {
	templatePath string
}

func NewTemplatesService(templatePath string) *TemplatesService {
	return &TemplatesService{
		templatePath: templatePath,
	}
}

func (s *TemplatesService) RenderTemplate(templateFile string, data any) (string, error) {
	startTime := time.Now()

	// Load template
	tmplPath := filepath.Join(s.templatePath, templateFile)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	// Calculate and log total execution time
	executionTime := time.Since(startTime)
	log.Printf("RenderTemplate completed in %v", executionTime)

	return buf.String(), nil
}

func (s *TemplatesService) RenderTemplateWithFuncs(templateFile string, data any) (string, error) {
	startTime := time.Now()

	// Load template
	tmplPath := filepath.Join(s.templatePath, templateFile)

	// Register the "add" function in the template so that the index can start from 1
	tmpl := template.New(filepath.Base(tmplPath)).Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
	})

	tmplWithAdd, err := tmpl.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmplWithAdd.Execute(&buf, data); err != nil {
		return "", err
	}

	// Calculate and log total execution time
	executionTime := time.Since(startTime)
	log.Printf("RenderTemplateWithFuncs completed in %v", executionTime)

	return buf.String(), nil
}
