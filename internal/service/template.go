package service

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// Define default template functions that will be available in all templates
func (s *TemplatesService) generateDefaultFuncs() template.FuncMap {
	return template.FuncMap{

		// Math utility: adds two integers
		"add": func(a, b int) int { return a + b },

		// String utilities
		"toUpper": strings.ToUpper,   // Converts a string to uppercase
		"toLower": strings.ToLower,   // Converts a string to lowercase
		"trim":    strings.TrimSpace, // Removes leading and trailing whitespace
		"title": func(s string) string {
			return cases.Title(language.English).String(s)
		}, // Capitalizes the first letter of each word (deprecated but still used)

		// Date/time formatting
		"formatDate": func(t time.Time, layout string) string {
			return t.Format(layout) // Formats a time value using a Go layout string
		},

		// String check: returns true if string is empty or only whitespace
		"isEmpty": func(s string) bool {
			return strings.TrimSpace(s) == ""
		},

		// Marks a string as safe HTML so it's not escaped in the output
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},

		// Increment an integer (e.g., index in templates starting from 1)
		"inc": func(i int) int { return i + 1 },
	}
}

func (s *TemplatesService) RenderTemplateWithFuncs(templateFile string, data any, customFuncs []template.FuncMap) (string, error) {
	startTime := time.Now()

	// Load template
	tmplPath := filepath.Join(s.templatePath, templateFile)

	// Generate the default set of template functions
	funcs := s.generateDefaultFuncs()

	// If a custom FuncMap is provided (at least one), merge it into the default
	if len(customFuncs) > 0 {
		for k, v := range customFuncs[0] {
			funcs[k] = v // Override or add the user-defined function
		}
	}

	// Register functions in the template
	tmpl := template.New(filepath.Base(tmplPath)).Funcs(funcs)

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
