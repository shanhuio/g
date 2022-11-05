// Copyright (C) 2022  Shanhu Tech Inc.
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
