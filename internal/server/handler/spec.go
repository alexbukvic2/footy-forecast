package handler

import (
	"net/http"

	"github.com/alexbukvic2/footy-forecast/docs"
)

const swaggerUIHTML = `<!DOCTYPE html>
<html>
<head>
  <title>footy-forecast API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({ url: "/openapi.json", dom_id: "#swagger-ui" });
  </script>
</body>
</html>
`

// Spec serves the OpenAPI spec and Swagger UI.
type Spec struct{}

// NewSpec constructs a Spec handler.
func NewSpec() *Spec {
	return &Spec{}
}

// ServeJSON handles GET /openapi.json.
func (h *Spec) ServeJSON(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(docs.SpecJSON)
}

// ServeUI handles GET /docs.
func (h *Spec) ServeUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(swaggerUIHTML))
}
