package entity

type Summary struct {
	Id        int32
	OwnerId   string
	Skills    string
	Bio       string
	Languages string
}

type GetRequestSummary struct {
	Id      int32
	OwnerId string
}

type GetAllRequestSummary struct {
	Page  int32
	Limit int32
	Field string
	Value string
}

type GetAllResponseSummary struct {
	Summary []*Summary
}

type StatusSummary struct {
	Status bool
}
