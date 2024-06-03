package models

type ReqClient struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RespClient struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

type ResponseClient struct {
	Id        string `json:"id"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type ResponseError struct {
	Error interface{} `json:"error"`
}

// ServerError ...
type ServerError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
