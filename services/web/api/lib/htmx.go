package lib

import (
	"html/template"
	"net/http"
)

// IsHTMXRequest checks if the request is made via HTMX
func IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// RenderTemplate renders the appropriate template based on whether the request is from HTMX or a regular browser request
// If it's an HTMX request, it renders only the content template
// If it's a regular request, it renders the full layout with the content template embedded
func RenderTemplate(w http.ResponseWriter, r *http.Request, templateFiles []string, data interface{}) error {
	var tmpl *template.Template
	var err error
	var templateName string

	// Parse all template files
	tmpl, err = template.ParseFiles(templateFiles...)
	if err != nil {
		return err
	}

	// Determine which template to render
	if IsHTMXRequest(r) {
		templateName = "content"
	} else {
		templateName = "layout"
	}

	// Execute the appropriate template
	return tmpl.ExecuteTemplate(w, templateName, data)
}
