package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string) error {

	t, err := template.New(templateName).Funcs(nil).ParseFiles(fmt.Sprintf("./templates/%s", templateName))
	if err != nil {
		log.Fatal(err)
	}

	t, err = t.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", "./templates"))
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}
