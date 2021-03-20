package render

import (
	// "fmt"
	// "io"
	"net/http"
	"testing"

	"github.com/jfk23/gobookings/internal/config"
	"github.com/jfk23/gobookings/internal/model"
)


func TestApplyTemplateData(t *testing.T) {

	var td model.TemplateData

	r, err := getSession()
	
	if err != nil  {
		t.Error(err)
	}

	sessionManager.Put(r.Context(), "flash", "123")

	result := ApplyTemplateData(&td, r)

	// _ = result

	if result.FlashMsg != "123" {
		t.Error("There is error with rendering")
	}

}

func TestRenderTemplate(t *testing.T) {
	pathToTemplate = "./../../templates"

	tc, er := CreateTemplate()

	if er != nil {
		t.Error(er)
	}

	app.CachedTemplate = tc


	var mw MyWriter
	r, err := getSession()
	
	if err != nil  {
		t.Error(err)
	}

	e := Template(&mw, r, "about.page.html", &model.TemplateData{})

	if e != nil {
		t.Errorf("supposed to be error-free")
	}

	e = Template(&mw, r, "non-exist.page.html", &model.TemplateData{})

	if e == nil {
		t.Errorf("supposed to be error")
	}
}

func TestCreateTemplate(t *testing.T) {

	pathToTemplate = "./../../templates"

	_, err := CreateTemplate()
	if err != nil {
		t.Error(err)
	}

}

func TestLoadTemplate(t *testing.T) {
	var b *config.AppConfig
	NewRenderer(b)
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some", nil)
		if err != nil {
			return nil, err
		}
	
	ctx := r.Context()

	c, _ := sessionManager.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(c)

	return r, nil

}


