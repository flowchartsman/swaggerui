package swaggerui

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"testing/fstest"
)

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
	overlay fstest.MapFS
}

func (s *swagWrapFS) Open(name string) (fs.File, error) {
	if _, err := s.overlay.Stat(name); err == nil {
		return s.overlay.Open(name)
	}
	return swagfs.Open(name)
}

func (s *swagWrapFS) ReadFile(name string) ([]byte, error) {
	if _, err := s.overlay.Stat(name); err == nil {
		return s.overlay.ReadFile(name)
	}
	return swagfs.ReadFile(name)
}

// I don't know how important it is that I throw in a synthesized entry for the specfile, so for now I don't
func (s *swagWrapFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return swagfs.ReadDir(name)
}

func addprefix(prefix string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = prefix + r.URL.Path
		next.ServeHTTP(w, r)
	})
}

// NewHandler returns a new handler that will serve swagger-ui with your spec
func NewHandler(spec []byte, spectype SpecType) http.Handler {
	//render the index template with the proper spec name inserted
	tmpl, _ := template.ParseFS(swagfs, "embed/index.html")
	idxbuf := new(bytes.Buffer)
	tmpl.ExecuteTemplate(idxbuf, "index.html", struct{ SpecType string }{string(spectype)})
	m := fstest.MapFS{
		"embed/index.html":                  {Data: idxbuf.Bytes()},
		"embed/swagger." + string(spectype): {Data: spec},
	}
	return addprefix("embed/", http.FileServer(http.FS(&swagWrapFS{overlay: m})))
}
