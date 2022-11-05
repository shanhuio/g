package paypal

import (
	"shanhu.io/pub/errcode"
)

func capturedOrderFromResponse(id string, r *captureOrderResponse) (
	*CapturedOrder, error,
) {
	const completed = "COMPLETED"
	if r.Status != completed {
		return nil, errcode.Internalf(
			"status is not completed but %q", r.Status,
		)
	}

	captured := &CapturedOrder{
		OrderID:          id,
		PayerCountryCode: payerCountry(r.Payer),
	}
	for _, unit := range r.PurchaseUnits {
		if unit.Payments == nil {
			return nil, errcode.Internalf(
				"purchase %q not captured", unit.ReferenceID,
			)
		}
		for _, capture := range unit.Payments.Captures {
			if capture != nil {
				captured.PaymentIDs = append(captured.PaymentIDs, capture.ID)
			}
		}
		if s := unit.Shipping; s != nil && s.Address != nil {
			address := *s.Address // make a copy
			shipping := &Shipping{
				ReferenceID: unit.ReferenceID,
				Address:     &address,
			}
			if s.Name != nil {
				shipping.To = s.Name.FullName
			}
			captured.Shippings = append(captured.Shippings, shipping)
		}
	}

	return captured, nil
}
