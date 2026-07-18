package user_http_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CreateReq struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Patronymic  *string `json:"patronymic"`
	PhoneNumber *string `json:"phone_number"`
}

type CreateResp struct {
	ID          string     `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Status      string     `json:"status"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Patronymic  *string    `json:"patronymic"`
	PhoneNumber *string    `json:"phone_number"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var createReq CreateReq
	err := json.NewDecoder(r.Body).Decode(&createReq)
	if err != nil {
		fmt.Println("ошибка")
	}
}
