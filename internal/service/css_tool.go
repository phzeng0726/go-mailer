package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/vanng822/go-premailer/premailer"
)

type CSSToolsService struct {
	tmplPath string
	tmplSvc  Templates
}

func NewCSSToolsService(tmplPath string, tmplSvc Templates) *CSSToolsService {
	return &CSSToolsService{
		tmplPath: tmplPath,
		tmplSvc:  tmplSvc,
	}
}

func (s *CSSToolsService) RenderTemplateWithCSS(templateFile, cssFile string, data any) (string, error) {
	return s.renderWithCSS(
		cssFile,
		func() (string, error) { return s.tmplSvc.RenderTemplate(templateFile, data) },
		"RenderTemplateWithCSS",
	)
}

func (s *CSSToolsService) RenderTemplateWithFuncsAndCSS(templateFile, cssFile string, data any) (string, error) {
	return s.renderWithCSS(
		cssFile,
		func() (string, error) { return s.tmplSvc.RenderTemplateWithFuncs(templateFile, data) },
		"RenderTemplateWithFuncsAndCSS",
	)
}

func (s *CSSToolsService) insertCSSIntoHead(htmlStr, cssContent string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Find <head> and insert <style>
	head := doc.Find("head")
	if head.Length() == 0 {
		// If no head tag exists, create one and add it before html or body
		doc.Find("html").PrependHtml("<head></head>")
		head = doc.Find("head")
	}

	head.AppendHtml("<style>" + cssContent + "</style>")

	// Convert back to HTML string
	html, err := doc.Html()
	if err != nil {
		return "", fmt.Errorf("failed to render HTML: %w", err)
	}

	return html, nil
}

func (s *CSSToolsService) convertCSSIntoInline(html string) (string, error) {
	opts := premailer.NewOptions()
	opts.KeepBangImportant = true
	opts.RemoveClasses = true

	prem, err := premailer.NewPremailerFromString(html, opts)
	if err != nil {
		return "", fmt.Errorf("failed to create premailer: %w", err)
	}

	inlineHTML, err := prem.Transform()
	if err != nil {
		return "", fmt.Errorf("failed to transform inline styles: %w", err)
	}

	return inlineHTML, nil
}
func (s *CSSToolsService) renderWithCSS(
	cssFile string,
	renderFunc func() (string, error),
	logName string,
) (string, error) {
	// Start timing
	startTime := time.Now()

	// Template rendering and CSS loading can be done in parallel
	type result struct {
		html string
		css  []byte
		err  error
	}

	htmlCh := make(chan result, 1)
	cssCh := make(chan result, 1)

	// Async template rendering
	go func() {
		html, err := renderFunc()
		htmlCh <- result{html: html, err: err}
	}()

	// Async CSS file reading
	go func() {
		cssPath := filepath.Join(s.tmplPath, cssFile)
		css, err := os.ReadFile(cssPath)
		cssCh <- result{css: css, err: err}
	}()

	// Get results
	htmlRes := <-htmlCh
	if htmlRes.err != nil {
		return "", fmt.Errorf("failed to render template: %w", htmlRes.err)
	}

	cssRes := <-cssCh
	if cssRes.err != nil {
		return "", fmt.Errorf("failed to read CSS file %s: %w", cssFile, cssRes.err)
	}

	// Insert CSS into head
	htmlWithCSS, err := s.insertCSSIntoHead(htmlRes.html, string(cssRes.css))
	if err != nil {
		return "", fmt.Errorf("failed to insert CSS into head: %w", err)
	}

	// Convert to inline styles
	inlineHTML, err := s.convertCSSIntoInline(htmlWithCSS)
	if err != nil {
		return "", fmt.Errorf("failed to convert to inline styles: %w", err)
	}

	// Calculate and log total execution time
	executionTime := time.Since(startTime)
	log.Printf("%s completed in %v", logName, executionTime)

	return inlineHTML, nil
}
