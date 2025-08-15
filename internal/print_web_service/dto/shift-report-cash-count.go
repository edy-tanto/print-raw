package dto

type Shift struct {
	Id             uint    `json:"id"`
	ShiftStartAt   string  `json:"shift_start_at"`
	ShiftEndAt     string  `json:"shift_end_at"`
	UserName       string  `json:"user_name"`
	OpeningBalance float32 `json:"opening_balance"`
}

type PaymentSummary struct {
	PaymentMethod string  `json:"payment_method"`
	GrandTotal    float32 `json:"grand_total"`
}

type CashCounts struct {
	Nominal float32 `json:"nominal"`
	Qty     uint    `json:"qty"`
	Total   float32 `json:"total"`
}

type ShiftReportCashCount struct {
	Shift                Shift            `json:"shift"`
	TotalPayment         float32          `json:"total_payment"`
	PaymentSummary       []PaymentSummary `json:"payment_summary"`
	CashCounts           []CashCounts     `json:"cash_counts"`
	CashCountsGrandTotal float32          `json:"cash_counts_grand_total"`
	CashCountsQtyTotal   uint             `json:"cash_counts_qty_total"`
}

type PrintShiftReportCashCountRequestBody struct {
	ShiftReportCashCount ShiftReportCashCount `json:"report_cash_count"`
}
