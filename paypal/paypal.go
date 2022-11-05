package paypal

import (
	"net/http"
	"net/url"
	"path"
	"strings"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
	"shanhu.io/pub/strutil"
)

// Config contiains configuration for the paypal module.
type Config struct {
	Account string
	ID      string
	Secret  string

	ReturnURL string `json:",omitempty"`
	CancelURL string `json:",omitempty"`
}

// PayPal is a web service module for handling paypal payments.
type PayPal struct {
	config *Config

	sandbox       bool
	returnURLBase *url.URL
	cancelURLBase *url.URL
	brand         string

	tokenClient *http.Client
}

func parseURL(u string) (*url.URL, error) {
	if u == "" {
		return nil, nil
	}
	return url.Parse(u)
}

// New creates a new Paypal module using the given config.
func New(config *Config, brand string) (*PayPal, error) {
	returnURL, err := parseURL(config.ReturnURL)
	if err != nil {
		return nil, errcode.Annotate(err, "parse return URL")
	}
	cancelURL, err := parseURL(config.CancelURL)
	if err != nil {
		return nil, errcode.Annotate(err, "parse cancel URL")
	}
	brand = strutil.Default(brand, "Shanhu Tech. Inc.")

	return &PayPal{
		config:        config,
		brand:         brand,
		sandbox:       strings.HasPrefix(config.Account, "sb-"),
		returnURLBase: returnURL,
		cancelURLBase: cancelURL,
		tokenClient:   http.DefaultClient,
	}, nil
}

func (p *PayPal) apiHost() string {
	if p.sandbox {
		return "api-m.sandbox.paypal.com"
	}
	return "api-m.paypal.com"
}

// IsSandBox returns true if this is a sandbox account.
func (p *PayPal) IsSandBox() bool { return p.sandbox }

func (p *PayPal) checkOutURL(base *url.URL, id string) string {
	if base == nil {
		return ""
	}

	q := make(url.Values)
	q.Set("ref", id)

	u := *base
	u.RawQuery = q.Encode()
	return u.String()
}

// ReturnURL returns the return URL for the given order.
func (p *PayPal) ReturnURL(id string) string {
	return p.checkOutURL(p.returnURLBase, id)
}

// CancelURL returns the cancel URL for the given order.
func (p *PayPal) CancelURL(id string) string {
	return p.checkOutURL(p.cancelURLBase, id)
}

func (p *PayPal) token() (string, error) {
	req := makeTokenRequest(p.apiHost(), p.config.ID, p.config.Secret)
	resp, err := p.tokenClient.Do(req)
	if err != nil {
		return "", errcode.Annotate(err, "fetch token")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", httputil.RespError(resp)
	}
	return tokenFromResponse(resp.Body)
}

func (p *PayPal) client() (*httputil.Client, error) {
	token, err := p.token()
	if err != nil {
		return nil, errcode.Annotate(err, "get token")
	}

	server := &url.URL{
		Scheme: "https",
		Host:   p.apiHost(),
	}
	return &httputil.Client{
		Server:      server,
		TokenSource: httputil.NewStaticToken(token),
	}, nil
}

// CreateOrder creates a new order.
func (p *PayPal) CreateOrder(o *CreateOrderRequest) (
	*CreateOrderResponse, error,
) {
	if o.Cents <= 0 {
		return nil, errcode.InvalidArgf(
			"invalid order value: %d cents", o.Cents,
		)
	}

	c, err := p.client()
	if err != nil {
		return nil, err
	}

	req := o.toPayPal(p.brand)
	resp := new(createOrderResponse)
	if err := c.Call("/v2/checkout/orders", req, resp); err != nil {
		return nil, errcode.Annotate(err, "create order")
	}

	ret := &CreateOrderResponse{OrderID: resp.ID}
	ret.ApproveURL = findApproveURL(resp.Links)
	if ret.ApproveURL == "" {
		return nil, errcode.Internalf(
			"approve url missing for order %q", resp.ID,
		)
	}
	return ret, nil
}

// OrderInfo returns the order info of a particular ID.
func (p *PayPal) OrderInfo(id string) (*OrderInfo, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}

	callPath := path.Join("/v2/checkout/orders", id)
	resp := new(getOrderResponse)
	if err := c.JSONGet(callPath, resp); err != nil {
		return nil, err
	}
	return &OrderInfo{
		OrderID:          id,
		Status:           resp.Status,
		PayerCountryCode: payerCountry(resp.Payer),
	}, nil
}

// CaptureOrder captures the order of a particular ID.
func (p *PayPal) CaptureOrder(id string) (*CapturedOrder, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}

	callPath := path.Join("/v2/checkout/orders", id, "capture")
	req := &captureOrderRequest{}
	resp := new(captureOrderResponse)
	if err := c.Call(callPath, req, resp); err != nil {
		return nil, errcode.Annotate(err, "paypal call")
	}

	const completed = "COMPLETED"
	if resp.Status != completed {
		return nil, errcode.Internalf(
			"status is not completed but %q", resp.Status,
		)
	}
	return capturedOrderFromResponse(id, resp)
}

// CreateSubscription creates a pending subscription.
func (p *PayPal) CreateSubscription(req *CreateSubscriptionRequest) (
	*CreateSubscriptionResponse, error,
) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}

	callReq := &createSubscriptionRequest{
		PlanID:   req.PlanID,
		CustomID: req.CustomID,
		AppContext: &appContext{
			BrandName:  p.brand,
			UserAction: "SUBSCRIBE_NOW",
			ReturnURL:  req.ReturnURL,
			CancelURL:  req.CancelURL,
		},
	}
	callResp := new(createSubscriptionResponse)
	const callPath = "/v1/billing/subscriptions"
	if err := c.Call(callPath, callReq, callResp); err != nil {
		return nil, errcode.Annotate(err, "paypal call")
	}

	const pending = "APPROVAL_PENDING"
	if callResp.Status != pending {
		// TODO(h8liu): what if user has already subscribed?
		return nil, errcode.Internalf(
			"status is not pending but %q", callResp.Status,
		)
	}
	resp := &CreateSubscriptionResponse{ID: callResp.ID}
	resp.ApproveURL = findApproveURL(callResp.Links)
	if resp.ApproveURL == "" {
		return nil, errcode.Internalf(
			"approve url missing for sub %q", callResp.ID,
		)
	}
	return resp, nil
}

// CancelSubscription cancels a subscription.
func (p *PayPal) CancelSubscription(subID, reason string) error {
	c, err := p.client()
	if err != nil {
		return err
	}

	callPath := path.Join("/v1/billing/subscriptions", subID, "cancel")
	req := &cancelSubscriptionRequest{Reason: reason}
	if err := c.Call(callPath, req, nil); err != nil {
		return errcode.Annotate(err, "paypal call")
	}
	return nil
}

// SubscriptionInfo gets the info of a subscription.
func (p *PayPal) SubscriptionInfo(subID string) (*SubscriptionInfo, error) {
	c, err := p.client()
	if err != nil {
		return nil, err
	}

	callPath := path.Join("/v1/billing/subscriptions", subID)
	callResp := new(getSubscriptionResponse)
	if err := c.JSONGet(callPath, callResp); err != nil {
		return nil, err
	}

	return &SubscriptionInfo{
		ID:       subID,
		Status:   callResp.Status,
		CustomID: callResp.CustomID,
		PlanID:   callResp.PlanID,
	}, nil
}

func (p *PayPal) apiToken(c *aries.C) (string, error) { return p.token() }

// API returns the paypal's admin API stub.
func (p *PayPal) API() aries.Service {
	r := aries.NewRouter()
	r.Call("token", p.apiToken)
	return r
}
