package convoy_go

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotListSubscriptionResponse = errors.New("invalid list subscription response")
	ErrNotSubscriptionResponse     = errors.New("invalid subscription response")
)

type Subscription struct {
	client *Client
}

type CreateSubscriptionRequest struct {
	Name       string `json:"name"`
	SourceID   string `json:"source_id"`
	EndpointID string `json:"endpoint_id"`

	AlertConfig  *AlertConfiguration  `json:"alert_config"`
	RetryConfig  *RetryConfiguration  `json:"retry_config"`
	FilterConfig *FilterConfiguration `json:"filter_config"`
}

type AlertConfiguration struct {
	Count     int    `json:"count"`
	Threshold string `json:"threshold"`
}

type RetryConfiguration struct {
	Type       string `json:"type"`
	Duration   string `json:"duration"`
	RetryCount int    `json:"retry_count"`
}

type FilterConfiguration struct {
	EventTypes []string `json:"event_types" bson:"event_types,omitempty"`
}

type SubscriptionResponse struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`

	Source   *SourceResponse   `json:"source_metadata,omitempty"`
	Endpoint *EndpointResponse `json:"endpoint_metadata,omitempty"`

	// subscription config
	AlertConfig  *AlertConfiguration  `json:"alert_config,omitempty"`
	RetryConfig  *RetryConfiguration  `json:"retry_config,omitempty"`
	FilterConfig *FilterConfiguration `json:"filter_config,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type SubscriptionQueryParam struct {
	GroupID    string `url:"groupId"`
	EndpointID string `url:"endpointId"`
	PerPage    int    `url:"per_page"`
	Page       int    `url:"page"`
}

type ListSubscriptionResponse struct {
	Content    []SubscriptionResponse `json:"content"`
	Pagination Pagination             `json:"pagination"`
}

func newSubscription(client *Client) *Subscription {
	return &Subscription{
		client: client,
	}
}

func (s *Subscription) All(ctx context.Context, query *SubscriptionQueryParam) (*ListSubscriptionResponse, error) {
	url, err := addOptions(s.generateUrl(), query)
	if err != nil {
		return nil, err
	}

	respPtr := &ListSubscriptionResponse{}
	err = getResource(ctx, s.client.apiKey, url, s.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (s *Subscription) Create(ctx context.Context, body *CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	url, err := addOptions(s.generateUrl(), nil)
	if err != nil {
		return nil, err
	}

	respPtr := &SubscriptionResponse{}
	err = postJSON(ctx, s.client.apiKey, url, body, s.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (s *Subscription) Find(ctx context.Context, subscriptionId string) (*SubscriptionResponse, error) {
	url, err := addOptions(s.generateUrl()+"/"+subscriptionId, nil)
	if err != nil {
		return nil, err
	}

	respPtr := &SubscriptionResponse{}
	err = getResource(ctx, s.client.apiKey, url, s.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (s *Subscription) Update(ctx context.Context, subscriptionId string, body *CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	url, err := addOptions(s.generateUrl()+"/"+subscriptionId, nil)
	if err != nil {
		return nil, err
	}

	respPtr := &SubscriptionResponse{}
	err = postJSON(ctx, s.client.apiKey, url, body, s.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (s *Subscription) Delete(ctx context.Context, subscriptionId string) error {
	url, err := addOptions(s.generateUrl()+"/"+subscriptionId, nil)
	if err != nil {
		return err
	}

	err = deleteResource(ctx, s.client.apiKey, url, s.client.client, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Subscription) generateUrl() string {
	return fmt.Sprintf("%s/projects/%s/subscriptions", s.client.baseURL, s.client.projectID)
}
