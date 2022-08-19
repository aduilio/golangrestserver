package dto

type TransferRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type TransferResponse struct {
	From AccountTransfer `json:"from"`
	To   AccountTransfer `json:"to"`
}

type AccountTransfer struct {
	Number  string  `json:"account_number"`
	Balance float64 `json:"balance"`
}
