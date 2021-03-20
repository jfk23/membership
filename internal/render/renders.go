package render

import (
	// "bytes"
	// "fmt"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/jfk23/gobookings/internal/config"
	"github.com/jfk23/gobookings/internal/model"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func ApplyTemplateData(t *model.TemplateData, r *http.Request) *model.TemplateData {
	t.CSRFToken = nosurf.Token(r)
	t.FlashMsg = app.Session.PopString(r.Context(), "flash")
	t.ErrorMsg = app.Session.PopString(r.Context(), "error")
	t.WarningMsg = app.Session.PopString(r.Context(), "warning")

	if app.Session.Exists(r.Context(), "user_id") {
		t.IsAuthenticated = 1
	}
	
	return t
}

//Template is rendering.
func Template(rw http.ResponseWriter, hr *http.Request, tmpl string, tempData *model.TemplateData) error {

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
	r, ok := appcache[tmpl]
	if !ok {
		return errors.New("no matching Templates")
	} 
	

	//buf := new(bytes.Buffer)

	tempData = ApplyTemplateData(tempData, hr)

	e := r.Execute(rw, tempData)
	//fmt.Println(buf.Len())
	
	if e != nil {
		log.Fatal(e)
		return e
	}
	// _, e = buf.WriteTo(rw)
	// if e != nil {
	// 	fmt.Println("There is error in writing template to buffer.")
	// }

	
	return nil
}

var funcmap = template.FuncMap{
	"humanTime" : HumanTime,
	"formatTime" : FormatTime,
	"iterate" : Iterate,
}

func Iterate (count int) []int {
	var i []int
	var m int
	for m = 1; m <= count; m++ {
		i = append(i, m)
	}
	
	return i
}

func HumanTime (t time.Time) string{
	return t.Format("2006-01-02")
}

func FormatTime (t time.Time, f string) string {
	return t.Format(f)
}

var pathToTemplate = "./templates"

func CreateTemplate() (map[string]*template.Template, error) {
	
	cached := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplate))
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

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))

		if err != nil {
			return cached, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
			if err != nil {
				return cached, err
			}
		}
		cached[name] = ts

	}

	return cached, nil
}
