package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/flowchartsman/swaggerui"
)

//go:embed spec/petstore.yml
var spec []byte

func main() {
	log.SetFlags(0)
	mux := http.NewServeMux()
	mux.HandleFunc("/pet/", petHandler)
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
