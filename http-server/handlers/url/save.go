package url

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/yokoshima228/url-shortener/internal/lib/logger/api"
	"github.com/yokoshima228/url-shortener/internal/lib/logger/sl"
	"github.com/yokoshima228/url-shortener/internal/lib/random"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	api.Response
	Alias string `json:"alias,omitempty"`
}

const (
	// TODO: move to config
	aliasLength = 6
)

type URLSaver interface {
	SaveURL(url string, alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode JSON body", sl.Err(err))
			render.JSON(w, r, api.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Invalid request", sl.Err(err))
			render.JSON(w, r, api.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		render.JSON(w, r, Response{
			Response: *api.Ok(),
			Alias:    alias,
		})
	}
}
