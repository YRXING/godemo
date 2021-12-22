package concurrentmap

import "errors"

func newIllegalParameterError(text string) error {
	return errors.New(text)
}
