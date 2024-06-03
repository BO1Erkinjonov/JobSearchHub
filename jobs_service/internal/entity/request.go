package entity

type Request struct {
	JobId           string
	ClientId        string
	SummaryId       int32
	StatusResp      string
	DescriptionResp string
}

type GetAllReq struct {
	Page  int32
	Limit int32
	Field string
	Value string
}

type StatusReq struct {
	Status bool
}

type GetRequest struct {
	JobId    string
	ClientId string
}
