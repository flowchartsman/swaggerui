# swaggerui
Embedded, self-hosted swagger-ui for go servers

## NOTE
This module depends on the `embed` package provided by the upcoming go 1.16 release, so it depends on having at least go 1.16 and, until that release lands, this project should be considered **experimental**

## Example usage
```go
package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/flowchartsman/swaggerui"
)

//go:embed swagger.json
var spec []byte

func main() {
	log.SetFlags(0)
	http.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.NewHandler(spec, swaggerui.SpecTypeJSON)))
	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```