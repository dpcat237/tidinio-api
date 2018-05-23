package api_controller

import (
	"net/http"
	"html/template"
	"bytes"
)

func Documentation(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/component/template/views/api/documentation_v2.html")
	if err != nil {
	}
	buf := new(bytes.Buffer)
	templateData := struct {
	}{
	}
	if err = t.Execute(buf, templateData); err != nil {
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.Error(w, buf.String(), http.StatusOK)
}
