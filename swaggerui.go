package swaggerui

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"testing/fstest"
)

//go:generate go run generate.go

//go:embed embed
var swagfs embed.FS

// SpecType represents one of the two types of Swagger spec
type SpecType string

// constants representing the type of spec your embedded swaggerui instance will be serving
const (
	SpecTypeJSON = `json`
	SpecTypeYAML = `yaml`
)

type swagWrapFS struct {
	static  fs.FS
	overlay fstest.MapFS
}

func (s *swagWrapFS) Open(name string) (fs.File, error) {
	if _, err := s.overlay.Stat(name); err == nil {
		return s.overlay.Open(name)
	}
	return s.static.Open(name)
}

// NewHandler returns a new handler that will serve swagger-ui with your spec
func NewHandler(spec []byte, spectype SpecType) http.Handler {
	//render the index template with the proper spec name inserted
	static, _ := fs.Sub(swagfs, "embed")
	tmpl, _ := template.ParseFS(static, "index.html")
	idxbuf := new(bytes.Buffer)
	tmpl.ExecuteTemplate(idxbuf, "index.html", struct{ SpecType string }{string(spectype)})
	m := fstest.MapFS{
		"index.html":                  {Data: idxbuf.Bytes()},
		"swagger." + string(spectype): {Data: spec},
	}
	return http.FileServer(http.FS(&swagWrapFS{
		static:  static,
		overlay: m,
	}))
}
