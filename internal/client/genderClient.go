package client

import "context"

type GenderClient interface {
	GetGenderByName(ctx context.Context, name string) (*string, error)
}
