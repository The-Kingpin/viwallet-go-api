package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"gitlab.com/the-kingpin/viwallet/internal/config"
	"gitlab.com/the-kingpin/viwallet/internal/models"
)

var app *config.AppConfig

const PathToTemplates = "./templates"

func NewRenderer(a *config.AppConfig) {
	app = a
}

// RenderTemplate renders template by given template name. tmpl is name of the desired template to be renderd
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		// avoid using the cache on every request for development purpose
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Println(tc)
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
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

// CreateTemplateCache creates a template cache as a map by finding the base name of page set it as key and
// for values assigns the pointer to the template
func CreateTemplateCache() (map[string]*template.Template, error) {
	tCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", PathToTemplates))

	if err != nil {
		return tCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", PathToTemplates))
		if err != nil {
			return tCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", PathToTemplates))
			if err != nil {
				return tCache, err
			}
		}

		tCache[name] = ts
	}

	return tCache, nil
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "username") {
		td.IsAuthenticated = true
	}
	return td
}
