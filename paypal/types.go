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

// Address is a payer's address.
type Address struct {
	Line1       string `json:"address_line_1,omitempty"`
	Line2       string `json:"address_line_2,omitempty"`
	City        string `json:"admin_area_2,omitempty"`
	State       string `json:"admin_area_1,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}

type payerName struct {
	First string `json:"given_name"`
	Last  string `json:"surname"`
}

type phoneNumber struct {
	Number string `json:"national_number"`
}

type phone struct {
	Type   string `json:"phone_type"`
	Number *phoneNumber
}

type payer struct {
	Name    *payerName `json:"name"`
	Email   string     `json:"email_address"`
	ID      string     `json:"payer_id"`
	Phone   *payer     `json:"phone,omitempty"`
	Address *Address   `json:"address"`
}

type purchaseUnit struct {
	ReferenceID string `json:"reference_id,omitempty"`
	InvoiceID   string `json:"invoice_id,omitempty"`
	Description string `json:"description,omitempty"`

	Shipping *shipping       `json:"shipping,omitempty"`
	Payments *payments       `json:"payments,omitempty"`
	Amount   *purchaseAmount `json:"amount,omitempty"`

	// Paypal generated ID
	ID string `json:"id,omitempty"`
}

type purchaseAmount struct {
	CurrencyCode string          `json:"currency_code"`
	Value        string          `json:"value"`
	Breakdown    *amountBreadown `json:"breakdown,omitempty"`
}

type shippingName struct {
	FullName string `json:"full_name,omitempty"`
}

type shipping struct {
	Name    *shippingName `json:"name,omitempty"`
	Address *Address      `json:"address"`
}

type payments struct {
	Captures []*capture `json:"captures"`
}

type capture struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type appContext struct {
	BrandName  string `json:"brand_name"`
	UserAction string `json:"user_action,omitempty"`
	ReturnURL  string `json:"return_url"`
	CancelURL  string `json:"cancel_url"`
}

type link struct {
	URL    string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type createOrderRequest struct {
	Intent        string          `json:"intent"`
	Description   string          `json:"description,omitempty"`
	PurchaseUnits []*purchaseUnit `json:"purchase_units"`
	AppContext    *appContext     `json:"application_context"`
}

type createOrderResponse struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Links  []*link `json:"links"`
}

type captureOrderRequest struct{}

type captureOrderResponse struct {
	ID            string          `json:"id"`
	Payer         *payer          `json:"payer"`
	Status        string          `json:"status"`
	PurchaseUnits []*purchaseUnit `json:"purchase_units"`
}

type getOrderResponse struct {
	ID            string          `json:"id"`
	Status        string          `json:"status"`
	Intent        string          `json:"intent"`
	PurchaseUnits []*purchaseUnit `json:"purchase_units"`
	Payer         *payer          `json:"payer"`
	CreateTime    string          `json:"create_time"`
	UpdateTime    string          `json:"update_time"`
	Links         []*link         `json:"links"`
}

type amountBreadown struct {
	ItemTotal *money `json:"item_total,omitempty"`
	TaxTotal  *money `json:"tax_total,omitempty"`
	Shipping  *money `json:"shipping,omitempty"`
}

type money struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type createSubscriptionRequest struct {
	PlanID     string      `json:"plan_id"`
	AppContext *appContext `json:"application_context"`
	CustomID   string      `json:"custom_id,omitempty"`
}

type createSubscriptionResponse struct {
	Status     string  `json:"status"`
	ID         string  `json:"id"`
	CreateTime string  `json:"create_time"`
	Links      []*link `json:"links"`
}

type getSubscriptionResponse struct {
	ID          string       `json:"id"`
	PlanID      string       `json:"plan_id"`
	CustomID    string       `json:"custom_id"`
	Status      string       `json:"status"`
	BillingInfo *billingInfo `json:"billing_info"`
}

type billingInfo struct {
	OutstandingBalance *billingBalance `json:"outstanding_balance"`
}

type billingBalance struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type cancelSubscriptionRequest struct {
	Reason string `json:"reason"`
}
