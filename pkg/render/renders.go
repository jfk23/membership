package render

import (
	// "bytes"
	// "fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/jfk23/gobookings/pkg/config"
	"github.com/jfk23/gobookings/pkg/model"
)

var app *config.AppConfig

func LoadTemplate(a *config.AppConfig) {
	app = a
}

func ApplyTemplateData(t *model.TemplateData) *model.TemplateData {
	return t
}

//RenderTemplate is rendering.
func RenderTemplate(rw http.ResponseWriter, tmpl string, tempData *model.TemplateData) {

	var appcache map[string]*template.Template

	if app.UseCache {
		appcache = app.CachedTemplate
	} else {
		appcache, _ = CreateTemplate()
	}

	// ts, err := CreateTemplate()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if r, ok := appcache[tmpl]; ok {

		//buf := new(bytes.Buffer)

		tempData = ApplyTemplateData(tempData)

		e := r.Execute(rw, tempData)
		//fmt.Println(buf.Len())
		
		if e != nil {
			log.Fatal(e)
		}
		// _, e = buf.WriteTo(rw)
		// if e != nil {
		// 	fmt.Println("There is error in writing template to buffer.")
		// }

	}
}

var funcmap template.FuncMap

func CreateTemplate() (map[string]*template.Template, error) {
	cached := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return cached, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		//fmt.Println("current file being processed is: ", name)

		ts, err := template.New(name).Funcs(funcmap).ParseFiles(page)

		if err != nil {
			return cached, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")

		if err != nil {
			return cached, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return cached, err
			}
		}
		cached[name] = ts

	}

	return cached, nil
}
