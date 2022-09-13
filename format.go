package pi

import (
	"encoding/json"
)

func Format[T any](ctx Context, p *T) error {
	_, r := ctx.Raw()
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}

	if v, ok := any(p).(Validator); ok {
		if err = v.Validate(r.Context()); err != nil {
			return err
		}
	}

	return nil
}
