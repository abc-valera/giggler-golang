package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var (
	openapiDocsPort = strings.TrimSpace(os.Getenv("OPENAPI_DOCS_PORT"))

	//go:embed scalar_docs.html
	docsHTML string

	//go:embed bundle.yaml
	bundle []byte
)

func main() {
	if openapiDocsPort == "" {
		log.Fatal("SQLITE_PATH env is required")
	}

	var docsBytes bytes.Buffer
	if err := template.Must(template.New("docs").Parse(docsHTML)).Execute(&docsBytes, "http://localhost:3001/bundle"); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(docsBytes.Bytes())
	})

	mux.HandleFunc("/bundle", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(bundle)
	})

	fmt.Println("Serving OpenAPI docs at http://localhost:3001/")
	if err := http.ListenAndServe(":3001", mux); err != nil {
		panic(err)
	}
}
