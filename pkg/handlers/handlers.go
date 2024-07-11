package handlers

import (
	"net/http"

	"github.com/ChunHou23/booking-service/pkg/config"
	"github.com/ChunHou23/booking-service/pkg/models"
	"github.com/ChunHou23/booking-service/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AboutMe(w http.ResponseWriter, r *http.Request) {
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "aboutMe.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
