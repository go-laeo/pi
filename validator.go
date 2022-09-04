package pi

import "context"

type Validator interface {
	Validate(ctx context.Context) error
}
