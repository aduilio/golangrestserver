package dto

type AccountRequest struct {
	Number string `json:"account_number"`
}
type AccountResponse struct {
	ID     string `json:"id"`
	Number string `json:"number"`
}
