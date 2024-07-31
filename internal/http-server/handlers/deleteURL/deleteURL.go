package deleteURL

import (
	"errors"
	"log/slog"
	"net/http"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate mockery --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("empty alias"))
			return
		}

		err := urlDeleter.DeleteURL(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not exists", "alias", alias)
			render.JSON(w, r, "url not exists")
			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, "internal error")
			return

		}

		log.Info("url deleted", slog.String("alias", alias))
		render.JSON(w, r, "url deleted")

	}
}
