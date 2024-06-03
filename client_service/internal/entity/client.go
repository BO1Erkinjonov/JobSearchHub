package entity

import "time"

type Client struct {
	Id           string
	Role         string
	FirstName    string
	LastName     string
	Email        string
	Password     string
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type GetRequest struct {
	ClientId string
	IsActive bool
}
type GetAllRequest struct {
	Page  int32
	Limit int32
	Field string
	Value string
}

type GetAllResponse struct {
	Clients []Client
	Count   int32
}

type DeleteReq struct {
	ClientId      string
	IsActive      bool
	IsHardDeleted bool
}

type Status struct {
	Status bool
}
