package service

import (
	"errors"

	"github.com/razorpay/razorpay-go"
)

type RazorpayService struct {
	Client *razorpay.Client
}

func NewRazorpayService(key, secret string) *razorpay.Client {
	return razorpay.NewClient(key, secret)
}

func (r *RazorpayService) ExecutePayment(amount int) (string, error) {

	data := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  "101",
	}

	body, err := r.Client.Order.Create(data, nil)
	if err != nil {
		return "", errors.New("payment not initiated")
	}
	razorId, _ := body["id"].(string)
	return razorId, nil
}
