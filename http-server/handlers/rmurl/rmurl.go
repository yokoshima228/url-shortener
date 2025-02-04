package rmurl

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/yokoshima228/url-shortener/internal/lib/logger/sl"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const position = "handlers.delete.New"

		log = log.With(
			slog.String("op", position),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("No alias was provided")
			render.JSON(w, r, "invalid request")
			return
		}

		log.Info("Got alias from user", sl.Info(alias))

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Error("Error deleting URL", sl.Err(err))
			render.JSON(w, r, "Error deleting url")
			return
		}

		log.Info("Deleted URL successfully")
		render.JSON(w, r, "Deleted URL")
	}
}
