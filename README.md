# swaggerui
Embedded, self-hosted [Swagger Ui](https://swagger.io/tools/swagger-ui/) for go servers

This module provides `swaggerui.Handler`, which you can use to serve an embedded copy of [Swagger UI](https://swagger.io/tools/swagger-ui/) as well as an embedded specification for your API.

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
	http.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
