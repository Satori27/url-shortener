package save

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"url-shortener/internal/lib/random"

	"errors"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	validator "github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const (
	aliasLength   = 6
	aliasAttempts = 10
)

var (
	ErrMaximumAttempts = errors.New("the maximum number of attempts was completed")
)

//go:generate mockery --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (uint64, error)
	AliasExist(alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias, err = generateAlias(urlSaver)
			if err != nil {
				log.Error("failed to generate alias", sl.Err(err))
				render.JSON(w, r, resp.Error("failed to generate alias"))
				return
			}

		}

		err = urlSaver.AliasExist(alias)
		if err == storage.ErrURLExists {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Uint64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})

	}
}

func generateAlias(urlSaver URLSaver) (string, error) {
	const op = "handlers.url.save.Alias"

	attempt := 0
	alias := random.NewRandomAlias(aliasLength)
	for ; attempt < aliasAttempts; attempt++ {
		err := urlSaver.AliasExist(alias)
		if err == storage.ErrURLNotFound || err == nil {
			return alias, nil
		} else if err != storage.ErrURLExists {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		alias = random.NewRandomAlias(aliasLength)

	}
	if attempt == aliasAttempts {
		return "", ErrMaximumAttempts
	}
	return "", nil

}
