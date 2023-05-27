package client

import "context"

type key string

var clientKey = key("client_context")

// Inject the client into the current context
func Inject(ctx context.Context, data *Client) context.Context {
	return context.WithValue(ctx, clientKey, data)
}

// Extract the client from the specified context
//
// Return error if client is not exists in the context
func FromContext(ctx context.Context) (*Client, error) {
	currentClient, ok := ctx.Value(clientKey).(*Client)
	if !ok {
		return nil, ErrInvalidClient
	}
	if currentClient == nil {
		return nil, ErrInvalidClient
	}
	if currentClient.Xid == "" {
		return nil, ErrInvalidClient
	}

	return currentClient, nil
}
