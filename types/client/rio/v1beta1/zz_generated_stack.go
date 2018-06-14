package client

import (
	"github.com/rancher/norman/types"
)

const (
	StackType                      = "stack"
	StackFieldCreated              = "created"
	StackFieldDescription          = "description"
	StackFieldLabels               = "labels"
	StackFieldName                 = "name"
	StackFieldRemoved              = "removed"
	StackFieldSpaceId              = "spaceId"
	StackFieldState                = "state"
	StackFieldTemplates            = "templates"
	StackFieldTransitioning        = "transitioning"
	StackFieldTransitioningMessage = "transitioningMessage"
	StackFieldUuid                 = "uuid"
)

type Stack struct {
	types.Resource
	Created              string            `json:"created,omitempty" yaml:"created,omitempty"`
	Description          string            `json:"description,omitempty" yaml:"description,omitempty"`
	Labels               map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	Removed              string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	SpaceId              string            `json:"spaceId,omitempty" yaml:"spaceId,omitempty"`
	State                string            `json:"state,omitempty" yaml:"state,omitempty"`
	Templates            map[string]string `json:"templates,omitempty" yaml:"templates,omitempty"`
	Transitioning        string            `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string            `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	Uuid                 string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}
type StackCollection struct {
	types.Collection
	Data   []Stack `json:"data,omitempty"`
	client *StackClient
}

type StackClient struct {
	apiClient *Client
}

type StackOperations interface {
	List(opts *types.ListOpts) (*StackCollection, error)
	Create(opts *Stack) (*Stack, error)
	Update(existing *Stack, updates interface{}) (*Stack, error)
	Replace(existing *Stack) (*Stack, error)
	ByID(id string) (*Stack, error)
	Delete(container *Stack) error
}

func newStackClient(apiClient *Client) *StackClient {
	return &StackClient{
		apiClient: apiClient,
	}
}

func (c *StackClient) Create(container *Stack) (*Stack, error) {
	resp := &Stack{}
	err := c.apiClient.Ops.DoCreate(StackType, container, resp)
	return resp, err
}

func (c *StackClient) Update(existing *Stack, updates interface{}) (*Stack, error) {
	resp := &Stack{}
	err := c.apiClient.Ops.DoUpdate(StackType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *StackClient) Replace(obj *Stack) (*Stack, error) {
	resp := &Stack{}
	err := c.apiClient.Ops.DoReplace(StackType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *StackClient) List(opts *types.ListOpts) (*StackCollection, error) {
	resp := &StackCollection{}
	err := c.apiClient.Ops.DoList(StackType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *StackCollection) Next() (*StackCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &StackCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *StackClient) ByID(id string) (*Stack, error) {
	resp := &Stack{}
	err := c.apiClient.Ops.DoByID(StackType, id, resp)
	return resp, err
}

func (c *StackClient) Delete(container *Stack) error {
	return c.apiClient.Ops.DoResourceDelete(StackType, &container.Resource)
}
