package helpers

import (
	"bytes"
	"io"
)

func CloneReadCloser(input io.ReadCloser) ([]byte, io.ReadCloser, error) {
	clonedInput, mErr := io.ReadAll(input)
	if mErr != nil {
		return nil, nil, mErr
	}

	rtnMe := io.NopCloser(bytes.NewBuffer(clonedInput))

	return clonedInput, rtnMe, nil

}
