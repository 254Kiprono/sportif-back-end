package services

type PaymentProvider interface {
	InitiatePayment(amount float64, currency string, description string) (string, error)
	VerifyPayment(reference string) (bool, error)
}

type StripePaymentProvider struct{}

func (s *StripePaymentProvider) InitiatePayment(amount float64, currency string, description string) (string, error) {
	// Implement Stripe Logic
	return "stripe_session_id", nil
}
func (s *StripePaymentProvider) VerifyPayment(reference string) (bool, error) {
	// Implement Stripe Verification
	return true, nil
}

type MpesaPaymentProvider struct{}

func (m *MpesaPaymentProvider) InitiatePayment(amount float64, currency string, description string) (string, error) {
	// Implement M-Pesa Logic
	return "mpesa_checkout_id", nil
}
func (m *MpesaPaymentProvider) VerifyPayment(reference string) (bool, error) {
	// Implement M-Pesa Verification
	return true, nil
}
