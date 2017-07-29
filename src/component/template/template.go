package app_template

import (
	"bytes"
	"html/template"
	"path/filepath"
)

func ParseEmailTemplate(templateFileName string, templateData interface{}) (string, error) {
	buf := new(bytes.Buffer)
	templates := make(map[string]*template.Template)
	files := []string{"src/component/template/layout/email_simple_layout.html", templateFileName}
	templates[filepath.Base("email_simple_layout.html")] = template.Must(template.ParseFiles(files...))
	if err := templates["email_simple_layout.html"].ExecuteTemplate(buf, "email_simple_layout.html", templateData); err != nil {
		return "", err
	}
	return buf.String(), nil
}
