package docker

import (
	"fmt"

	"github.com/docker/docker/client"
)

// Client wraps the official Docker Go SDK client
type Client struct {
	cli *client.Client
}

// NewClient initializes a Docker Engine API client from environment variables
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	return &Client{cli: cli}, nil
}

// RawClient returns the underlying *client.Client
func (c *Client) RawClient() *client.Client {
	return c.cli
}

// Close closes the client transport
func (c *Client) Close() error {
	if c.cli != nil {
		return c.cli.Close()
	}
	return nil
}
