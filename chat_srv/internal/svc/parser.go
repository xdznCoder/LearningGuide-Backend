package svc

import (
	"code.sajari.com/docconv"
	"context"
	"errors"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"io"
	"path"
)

const ContextFileNameKey = "CostumedFileName"

type Parser struct{}

func (p *Parser) fileName(ctx context.Context) string {
	fileName, ok := ctx.Value(ContextFileNameKey).(string)
	if !ok {
		return ""
	}
	return fileName
}

func (p *Parser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	fileName := p.fileName(ctx)

	var docs []*schema.Document
	var result *docconv.Response

	var err error

	switch path.Ext(fileName) {
	case ".md":
		content, err := io.ReadAll(reader)

		if err != nil {
			return docs, err
		}
		docs = append(docs, &schema.Document{
			Content: string(content),
			ID:      fileName,
		})
		return docs, nil

	case ".docx":
		result, err = docconv.Convert(reader, docconv.MimeTypeByExtension(fileName), true)
	case ".doc":
		result, err = docconv.Convert(reader, docconv.MimeTypeByExtension(fileName), true)
	case ".pdf":
		result, err = docconv.Convert(reader, docconv.MimeTypeByExtension(fileName), true)
	case ".pptx":
		result, err = docconv.Convert(reader, docconv.MimeTypeByExtension(fileName), true)
	case ".txt":
		result, err = docconv.Convert(reader, docconv.MimeTypeByExtension(fileName), true)
	default:
		return docs, errors.New("file type not support: " + fileName)
	}

	if err != nil {
		return docs, err
	}

	docs = append(docs, &schema.Document{
		Content: result.Body,
		ID:      fileName,
	})

	return docs, nil
}
