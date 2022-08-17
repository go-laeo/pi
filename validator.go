package ezy

import "context"

type Validator interface {
	Validate(ctx context.Context) error
}
