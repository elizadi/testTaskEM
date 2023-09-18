package domain

import "errors"

type Errs []error

func (e Errs) String() string {
	res := ""
	for _, err := range e {
		res += err.Error() + ", "
	}
	return res
}

var ErrEmptyMessage = errors.New("empty message")
