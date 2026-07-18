package model

import "time"

type UserStatus string

const (
	UserStatusUnspecified UserStatus = "UNSPECIFIED"
	UserStatusNew         UserStatus = "NEW"
	UserStatusActive      UserStatus = "ACTIVE"
	UserStatusDeleted     UserStatus = "DELETED"
	UserStatusBlocked     UserStatus = "BLOCKED"
)

type User struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Status      UserStatus
	FirstName   string
	LastName    string
	Patronymic  *string
	PhoneNumber *string
}

func NewUser(
	id string,
	firstName string,
	lastName string,
	patronymic *string,
	phoneNumber *string,
) *User {
	return &User{
		ID:          id,
		Status:      UserStatusNew,
		FirstName:   firstName,
		LastName:    lastName,
		Patronymic:  patronymic,
		PhoneNumber: phoneNumber,
	}
}
