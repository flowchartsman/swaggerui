# swaggerui
Embedded, self-hosted [Swagger Ui](https://swagger.io/tools/swagger-ui/) for go servers

This module provides `swaggerui.Handler`, which you can use to serve an embedded copy of [Swagger UI](https://swagger.io/tools/swagger-ui/) as well as an embedded specification for your API.

## NOTE
This module depends on the `embed` package provided by the upcoming go 1.16 release, so it depends on having at least go 1.16 and, until that release lands, this project should be considered **experimental**. Once 1.16 lands, the repository will be version tagged to the latest major release of Swagger UI, using the included `generate.go` to keep up with future releases as necessary for tracking changes.

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
	http.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec, swaggerui.SpecTypeJSON)))
	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
