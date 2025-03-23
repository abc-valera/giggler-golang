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
	port = strings.TrimSpace(os.Getenv("OPENAPI_PORT"))

	//go:embed scalar_docs.html
	scalarDocsHtml string

	//go:embed openapi_bundle.yaml
	openapiBundleYaml []byte
)

func main() {
	if port == "" {
		log.Fatal("OPENAPI_PORT env is required")
	}

	var docsBytes bytes.Buffer
	if err := template.Must(template.New("docs").Parse(scalarDocsHtml)).Execute(&docsBytes, port); err != nil {
		log.Fatalf("failed to parse docs template: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(docsBytes.Bytes())
	})

	mux.HandleFunc("/bundle", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(openapiBundleYaml)
	})

	fmt.Println("Serving OpenAPI docs at http://localhost:" + port + "/")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
