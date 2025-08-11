package proxy

import (
	"io"
	"os"
)

func BufferToFile(content io.ReadCloser, tempFilePattern string) (*os.File, io.ReadCloser, error) {

	tempFile, err := os.CreateTemp("", "prefix")
	if err != nil {
		return nil, nil, err
	}

	rtnMe := io.NopCloser(io.TeeReader(content, tempFile))

	return tempFile, rtnMe, nil
}
