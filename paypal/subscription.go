package paypal

// CreateSubscriptionRequest is the request to create a subscription.
type CreateSubscriptionRequest struct {
	PlanID    string
	CustomID  string
	ReturnURL string
	CancelURL string
}

// CreateSubscriptionResponse is the response of a pending subscription
// creation.
type CreateSubscriptionResponse struct {
	ID         string
	ApproveURL string
}

// SubscriptionInfo returns the subscription info.
type SubscriptionInfo struct {
	ID       string
	PlanID   string
	CustomID string
	Status   string
}
