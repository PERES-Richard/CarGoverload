package utils

type BankAPI struct {
	Host      string    `json:"host"`
	Port       string    `json:"port"`
	PaymentEP string `json:"paymentEP"`
}

func (b BankAPI) PerformPayment(s string, f float64) bool {
	return true;
}
