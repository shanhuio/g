// Copyright (C) 2023  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
