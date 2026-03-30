package templates

import (
	"html/template"
	"net/http"
	"net/url"
)

var tmpl *template.Template

func Load() {
	funcMap := template.FuncMap{
		"favicon": func(rawURL string) string {
			u, err := url.Parse(rawURL)
			if err != nil {
				return ""
			}
			return "https://www.google.com/s2/favicons?domain=" + u.Hostname() + "&sz=32"
		},
	}

	var err error
	tmpl, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}
}

func Render(w http.ResponseWriter, name string, data any) {
	t, err := template.New("").Funcs(template.FuncMap{
		"favicon": func(rawURL string) string {
			u, err := url.Parse(rawURL)
			if err != nil {
				return ""
			}
			return "https://www.google.com/s2/favicons?domain=" + u.Hostname() + "&sz=32"
		},
	}).ParseFiles("templates/layout.html", "templates/"+name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
