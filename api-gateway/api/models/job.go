package models

type JobReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type JobUpdateReq struct {
	JobId       string `json:"job_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Responses   int32  `json:"number_of_people"`
}

type JobsOwner struct {
	Id          string         `json:"id"`
	OwnerId     string         `json:"owner_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Responses   int32          `json:"responses"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	DeletedAt   string         `json:"deleted_at"`
	Owners      ResponseClient `json:"owner"`
}

type JobResponse struct {
	Id          string `json:"id"`
	OwnerId     string `json:"owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Responses   int32  `json:"responses"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
}

type ListJobs struct {
	Jobs []JobsOwner `json:"jobs"`
}
