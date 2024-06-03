package entity

import "time"

type Job struct {
	Id          string
	Owner_id    string
	Title       string
	Description string
	Response    int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type GetReq struct {
	Id       string
	IsActive bool
}

type GetAll struct {
	Page  int32
	Limit int32
	Field string
	Value string
}

type StatusJob struct {
	Status bool
}

type DelReq struct {
	Id            string
	IsActive      bool
	IsHardDeleted bool
}
