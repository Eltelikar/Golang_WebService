package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gotest_26.08.25/internal/http-server/response"
	"gotest_26.08.25/internal/postgre"
)

type List interface {
	List(int64, int64) ([]postgre.RequestFields, int64, error)
}

type ListResponse struct {
	Status        string                  `json:"status"`
	Message       string                  `json:"message"`
	TotalRecords  int64                   `json:"total_records"`
	CurrentPage   int64                   `json:"current_page"`
	TotalPages    int64                   `json:"total_pages"`
	Subscriptions []postgre.RequestFields `json:"subscriptions"`
}

// NewList возвращает хендлер, возвращающий все подписки
//
// @Summary Получить список всех подписок
// @Description Возвращает все подписки
// @Tags subscriptions
// @Produce json
// @Param page query int false "Номер страницы (>=1)" default(1)
// @Param page_size query int false "Размер страницы (1..100)" default(20)
// @Success 200 {object} ListResponse
// @Failure 500 {object} response.Response
// @Router /api/v1/subscriptions [get]
func NewList(log *slog.Logger, storage List) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.NewList"

		const (
			defaultPage     = 1
			defaultPageSize = 20
			maxPageSize     = 100
		)

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("List handler started")

		q := r.URL.Query()
		page := defaultPage
		pageSize := int64(defaultPageSize)

		if value := q.Get("page"); value != "" {
			if n, err := strconv.Atoi(value); err != nil {
				log.Error("invalid page number", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid page number"))
				return
			} else if n > 0 {
				page = n
			}
		}

		if value := q.Get("page_size"); value != "" {
			if n, err := strconv.Atoi(value); err != nil {
				log.Error("invalid page size number", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid page sizenumber"))
				return
			} else if n > 0 && n <= maxPageSize {
				pageSize = int64(n)
			}
		}

		offset := int64(page-1) * pageSize

		subscriptions, count, err := storage.List(pageSize, offset)
		if err != nil {
			log.Error("Failed to list subscriptions", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		var totalPages int64
		if pageSize > 0 {
			totalPages = (count + pageSize - 1) / pageSize
		}

		log.Info("Subscriptions listed successfully")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, ListResponse{
			Status:        "success",
			Message:       "Subscriptions listed successfully",
			TotalRecords:  count,
			CurrentPage:   int64(page),
			TotalPages:    totalPages,
			Subscriptions: subscriptions,
		})
	}
}
