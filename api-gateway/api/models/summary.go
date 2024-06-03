package models

type Summary struct {
	Skills    string `json:"skills"`
	Bio       string `json:"bio"`
	Languages string `json:"languages"`
}

type SummaryResponse struct {
	Id        int32  `json:"id"`
	OwnerId   string `json:"owner_id"`
	Skills    string `json:"skills"`
	Bio       string `json:"bio"`
	Languages string `json:"languages"`
}

type SummaryUpdateRequest struct {
	Id        int32  `json:"id"`
	Skills    string `json:"skills"`
	Bio       string `json:"bio"`
	Languages string `json:"languages"`
}

type ListSummary struct {
	Summary []SummaryResponse `json:"summary"`
}
