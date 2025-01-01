package save

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/leedinh/telebot/bitly/internal/lib/api/response"
	"github.com/leedinh/telebot/bitly/internal/lib/bloomfilter"
	"github.com/leedinh/telebot/bitly/internal/lib/hasher"
	"github.com/leedinh/telebot/bitly/internal/lib/logger/sl"
	"github.com/leedinh/telebot/bitly/internal/storage"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 7

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, bf *bloomfilter.BloomFilter, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.With("op", op),
			slog.With("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("empty request body")

			render.JSON(w, r, resp.Error("empty request body"))

			return
		}

		if err != nil {
			log.Error("failed to decode the request", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		// log.Info("request has been decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := hasher.GenerateAlias(aliasLength, req.URL, bf)

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrorURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, Response{
				Response: resp.Error("URL already exists"),
			})

			return
		}

		if err != nil {
			log.Error("failed to save the URL", sl.Err(err))

			render.JSON(w, r, Response{
				Response: resp.Error("failed to save the URL"),
			})

			return
		}

		bf.Add([]byte(alias))
		log.Info("URL has been saved", id)

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
