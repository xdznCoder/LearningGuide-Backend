package tpl

import "github.com/cloudwego/eino/schema"

type TemplateType int

const (
	TemplateTypeUserQuery TemplateType = iota + 1
	TemplateTypeExerciseGenerate
	TemplateTypeMindMapGenerate
)

type internalMap struct {
	m map[TemplateType]*schema.Message
}

var Map internalMap

func init() {
	Map = internalMap{m: make(map[TemplateType]*schema.Message)}
}

func (m internalMap) Get(t TemplateType) *schema.Message {
	return m.m[t]
}
