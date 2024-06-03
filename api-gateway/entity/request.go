package entity

type Request struct {
	JobId           string
	ClientId        string
	SummaryId       int32
	StatusResp      string
	DescriptionResp string
}

type RequestResp struct {
	JobId     string `json:"job_id"`
	ClientId  string `json:"client_id"`
	SummaryId int32  `json:"summary_id"`
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

type GetRequestReq struct {
	JobId    string
	ClientId string
}
