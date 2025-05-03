package dto

type CashRefundDetail struct {
	Item     string  `json:"item"`
	Qty      uint    `json:"qty"`
	Subtotal float32 `json:"subtotal"`
}

type CashRefund struct {
	Id                uint               `json:"id"`
	Op                string             `json:"op"`
	Date              string             `json:"date"`
	TotalRefund       float32            `json:"total_refund"`
	CashRefundDetails []CashRefundDetail `json:"cash_refund_details"`
}

type PrintCashRefundRequestBody struct {
	CashRefund CashRefund `json:"cash_refund"`
}