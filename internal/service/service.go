package service

type Templates interface {
	RenderTemplate(tmplFile string, data any) (string, error)
	RenderTemplateWithFuncs(tmplFile string, data any) (string, error)
}

type CSSTools interface {
	RenderTemplateWithCSS(tmplFile, cssFile string, data any) (string, error)
	RenderTemplateWithFuncsAndCSS(tmplFile, cssFile string, data any) (string, error)
}

type Services struct {
	Templates Templates
	CSSTools  CSSTools
}

type Deps struct {
	TmplPath string
	CSSPath  string
}

func NewServices(deps Deps) *Services {
	tmplSvc := NewTemplatesService(
		deps.TmplPath,
	)

	cssToolsSvc := NewCSSToolsService(
		deps.CSSPath,
		tmplSvc,
	)

	return &Services{
		Templates: tmplSvc,
		CSSTools:  cssToolsSvc,
	}
}
