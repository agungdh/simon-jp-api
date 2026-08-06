package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"simon-jp-api/internal/service"
)

type Deps struct {
	Auth      *service.AuthService
	Dashboard *service.DashboardService
	Pegawai   *service.PegawaiService
	Diklat    *service.DiklatService
	Ppm       *service.PpmService
	Activity  *service.ActivityService
	MateriPpm *service.MateriPpmService
	Download  *service.DownloadService
}

func NewRouter(deps Deps, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(requestLogger(logger))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Get("/healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	authHandler := NewAuthHandler(deps.Auth)
	r.Post("/api/login", authHandler.Login)

	r.Route("/api", func(r chi.Router) {
		r.Use(authHandler.Middleware())

		r.Post("/logout", authHandler.Logout)
		r.Get("/user", authHandler.User)

		if deps.Dashboard != nil {
			r.Get("/dashboard", NewDashboardHandler(deps.Dashboard).Index)
		}

		if deps.Download != nil {
			r.Get("/download/{module}/{uuid}", NewDownloadHandler(deps.Download).Download)
		}

		if deps.Pegawai != nil {
			handler := NewPegawaiHandler(deps.Pegawai)
			r.Get("/pegawai", handler.List)
			r.Get("/pegawai/{uuid}", handler.Show)
		}

		if deps.Diklat != nil {
			handler := NewDiklatHandler(deps.Diklat)
			r.Get("/diklat", handler.List)
			r.Get("/diklat/{uuid}", handler.Show)
		}

		if deps.Ppm != nil {
			handler := NewPpmHandler(deps.Ppm)
			r.Get("/ppm", handler.List)
			r.Get("/ppm/{uuid}", handler.Show)
		}

		if deps.Activity != nil {
			handler := NewActivityHandler(deps.Activity)
			r.Get("/seminar", handler.ListSeminars)
			r.Get("/seminar/{uuid}", handler.ShowSeminar)
			r.Get("/webinar", handler.ListWebinars)
			r.Get("/webinar/{uuid}", handler.ShowWebinar)
			r.Get("/lc", handler.ListLcs)
			r.Get("/lc/{uuid}", handler.ShowLc)
			r.Get("/belajar-mandiri", handler.ListBelajarMandiris)
			r.Get("/belajar-mandiri/{uuid}", handler.ShowBelajarMandiri)
			r.Get("/mentoring", handler.ListMentorings)
			r.Get("/mentoring/{uuid}", handler.ShowMentoring)
			r.Get("/coaching", handler.ListCoachings)
			r.Get("/coaching/{uuid}", handler.ShowCoaching)
			r.Get("/workshop", handler.ListWorkshops)
			r.Get("/workshop/{uuid}", handler.ShowWorkshop)
		}

		if deps.MateriPpm != nil {
			handler := NewMateriPpmHandler(deps.MateriPpm)
			r.Get("/materi-ppm", handler.List)
			r.Get("/materi-ppm/{uuid}", handler.Show)
		}
	})

	return r
}
