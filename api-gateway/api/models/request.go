package models

type Request struct {
	JobId     string `json:"job_id"`
	SummaryId int32  `json:"summary_id"`
}

type RequestResp struct {
	JobId     string `json:"job_id"`
	ClientId  string `json:"client_id"`
	SummaryId int32  `json:"summary_id"`
}

type ListRequest struct {
	Requests []RequestListResp `json:"requests"`
}

type RequestListResp struct {
	JobId           string `json:"job_id"`
	ClientId        string `json:"client_id"`
	SummaryId       int32  `json:"summary_id"`
	StatusResp      string `json:"status_resp"`
	DescriptionResp string `json:"description_resp"`
}

type RequestResponse struct {
	JobId           string `json:"job_id"`
	StatusResp      string `json:"status_resp"`
	DescriptionResp string `json:"description_resp"`
}

type RequestReq struct {
	JobId           string `json:"job_id"`
	ClientId        string `json:"client_id"`
	StatusResp      string `json:"status_resp"`
	DescriptionResp string `json:"description_resp"`
}
