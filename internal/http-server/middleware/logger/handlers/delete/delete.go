package delete

import (
	"errors"
	"log/slog"
	"net/http"

	resp "urlrest/internal/lib/api/response"
	"urlrest/internal/lib/logger/sl"
	"urlrest/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UrlDelete interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter UrlDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		queryAlias := r.URL.Query().Get("alias")
		if queryAlias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("alias is required"))
			return
		}

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := urlDeleter.DeleteURL(queryAlias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", "alias", queryAlias)
				render.JSON(w, r, resp.Error("not found"))
			} else {
				log.Error("failed to delete url", sl.Err(err))
				render.JSON(w, r, resp.Error("internal error"))
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}

}
