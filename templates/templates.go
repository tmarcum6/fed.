package templates

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func Load() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func Render(w http.ResponseWriter, name string, data any) {
	// Parse layout + the specific page template together each time
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/"+name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Always execute layout.html — it has the DOCTYPE and calls {{ block "content" }}
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
