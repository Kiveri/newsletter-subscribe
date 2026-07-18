package http_transport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Kiveri/newsletter-subscribe/internal/pkg/zap_logger"
	"go.uber.org/zap"
)

type ResponseHandler struct {
	log *zap_logger.Logger
	w   http.ResponseWriter
}

func NewResponseHandler(
	log *zap_logger.Logger,
	w http.ResponseWriter,
) *ResponseHandler {
	return &ResponseHandler{
		log: log,
		w:   w,
	}
}

func (h *ResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.w.WriteHeader(http.StatusInternalServerError)

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	if err = json.NewEncoder(h.w).Encode(response); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}
