package paypal

import (
	"fmt"
)

func valueStr(cents int) string {
	minor := cents % 100
	major := cents / 100
	if major == 0 {
		return fmt.Sprintf("0.%02d", minor)
	}
	return fmt.Sprintf("%d.%02d", major, minor)
}

func newMoney(cents int, currency string) *money {
	if cents == 0 {
		return nil
	}
	if currency == "" {
		currency = "USD"
	}
	return &money{
		CurrencyCode: currency,
		Value:        valueStr(cents),
	}
}

// CreateOrderRequest is the request for creating an order.
type CreateOrderRequest struct {
	ReferenceID   string
	InvoiceID     string
	Description   string
	Cents         int
	TaxCents      int
	ShippingCents int
	Currency      string
	ReturnURL     string
	CancelURL     string
}

func (r *CreateOrderRequest) amount() *purchaseAmount {
	currency := "USD" // default currency
	if r.Currency != "" {
		currency = r.Currency
	}

	total := r.Cents + r.TaxCents + r.ShippingCents
	totalValue := valueStr(total)

	amount := &purchaseAmount{
		CurrencyCode: currency,
		Value:        totalValue,
	}
	if total != r.Cents {
		amount.Breakdown = &amountBreadown{
			ItemTotal: newMoney(r.Cents, currency),
			TaxTotal:  newMoney(r.TaxCents, currency),
			Shipping:  newMoney(r.ShippingCents, currency),
		}
	}
	return amount
}

func (r *CreateOrderRequest) toPayPal(brand string) *createOrderRequest {
	return &createOrderRequest{
		Intent:      "CAPTURE",
		Description: r.Description,
		PurchaseUnits: []*purchaseUnit{{
			ReferenceID: r.ReferenceID,
			InvoiceID:   r.InvoiceID,
			Description: r.Description,
			Amount:      r.amount(),
		}},
		AppContext: &appContext{
			BrandName:  brand,
			UserAction: "PAY_NOW",
			ReturnURL:  r.ReturnURL,
			CancelURL:  r.CancelURL,
		},
	}

}

// CreateOrderResponse is the response for creating an order.
type CreateOrderResponse struct {
	OrderID    string
	ApproveURL string
}

// CapturedOrder is the response for capturing an order.
type CapturedOrder struct {
	OrderID          string
	PayerCountryCode string
	Shippings        []*Shipping
	PaymentIDs       []string
}

// Shipping contains shipping information.
type Shipping struct {
	ReferenceID string
	To          string
	Address     *Address
}

// OrderInfo contains basic order information.
type OrderInfo struct {
	OrderID          string
	Status           string
	PayerCountryCode string
}

func payerCountry(payer *payer) string {
	if payer == nil {
		return ""
	}
	if addr := payer.Address; addr != nil {
		return addr.CountryCode
	}
	return ""
}
