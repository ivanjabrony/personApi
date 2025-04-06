package client

import "context"

type NationalityClient interface {
	GetNationalityByName(ctx context.Context, name string) (*string, error)
}
