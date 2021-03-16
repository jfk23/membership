package handlers

import (
	"net/http"

	"github.com/jfk23/gobookings/pkg/config"
	"github.com/jfk23/gobookings/pkg/model"
	"github.com/jfk23/gobookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	ConfigSetting *config.AppConfig
}

func CreateNewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		ConfigSetting: a,
	}
}

func SetHandler(r *Repository) {
	Repo = r
}

//Home is page for /
func (re *Repository) Home(rw http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(rw, "home.page.html", &model.TemplateData{})

	remoteIP := r.RemoteAddr

	re.ConfigSetting.Session.Put(r.Context(), "remote_ip", remoteIP)

}

// About is page for /
func (re *Repository) About(rw http.ResponseWriter, r *http.Request) {
	remoteIP := re.ConfigSetting.Session.GetString(r.Context(), "remote_ip")

	var stringData = map[string]string{"test": "Hello Good day!", "remote_ip": remoteIP}
	render.RenderTemplate(rw, "about.page.html", &model.TemplateData{
		StringMap: stringData,
	})

}
