package pi

import (
	"encoding/json"
)

// Format decodes request body as JSON object to *T.
func Format[T any](ctx Context, p *T) error {
	_, r := ctx.Raw()
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}

	return nil
}

// FormatValidator runs Format() on *P, then tries call (Validator).Validate() on it.
func FormatValidator[T any](ctx Context, p *T) error {
	err := Format(ctx, p)
	if err != nil {
		return err
	}

	if v, ok := any(p).(Validator); ok {
		if err = v.Validate(ctx.Context()); err != nil {
			return err
		}
	}

	return nil
}
