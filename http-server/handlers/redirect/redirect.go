package redirect

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/yokoshima228/url-shortener/internal/lib/logger/api"
	"github.com/yokoshima228/url-shortener/internal/lib/logger/sl"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const position = "handlers.redirect.New"

		log = log.With(
			slog.String("op", position),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias was empty")
			render.JSON(w, r, "invalid request")
			return
		}

		log.Info("Got alias from user:", sl.Info(alias))

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			log.Error("Error getting url", sl.Err(err))
			render.JSON(w, r, api.Error("Error redirecting"))
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
