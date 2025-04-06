package client

import "context"

type AgeClient interface {
	GetAgeByName(ctx context.Context, name string) (*int, error)
}
