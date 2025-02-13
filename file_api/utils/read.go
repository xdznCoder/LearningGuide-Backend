package utils

import (
	"code.sajari.com/docconv"
	"errors"
	"io"
	"path"
)

var ErrInvalidType = errors.New("invalid type")

func ReadFile(file io.Reader, fileName string) (string, error) {
	var result *docconv.Response

	var err error

	switch path.Ext(fileName) {
	case ".md":
		content, err := io.ReadAll(file)

		if err != nil {
			return "", err
		}
		return string(content), nil

	case ".docx":
		result, err = docconv.Convert(file, docconv.MimeTypeByExtension(fileName), true)
	case ".doc":
		result, err = docconv.Convert(file, docconv.MimeTypeByExtension(fileName), true)
	case ".pdf":
		result, err = docconv.Convert(file, docconv.MimeTypeByExtension(fileName), true)
	case ".pptx":
		result, err = docconv.Convert(file, docconv.MimeTypeByExtension(fileName), true)
	case ".txt":
		result, err = docconv.Convert(file, docconv.MimeTypeByExtension(fileName), true)
	default:
		return "", ErrInvalidType
	}

	if err != nil {
		return "", err
	}

	return result.Body, nil
}
