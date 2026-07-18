package user_http_handler

import (
	"net/http"

	"github.com/Kiveri/newsletter-subscribe/internal/adapter/api/http/http_transport"
)

type (
	Handler struct {
		userUseCase userUseCase
	}

	userUseCase interface{}
)

func NewHandler(userUseCase userUseCase) *Handler {
	return &Handler{
		userUseCase: userUseCase,
	}
}

func (h *Handler) Routes() []http_transport.Route {
	return []http_transport.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.Create,
		},
	}
}
