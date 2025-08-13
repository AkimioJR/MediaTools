package schemas

type ErrResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
