package tpl

import (
	"context"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type TemplateType int

const (
	TemplateTypeUserQuery TemplateType = iota + 1
	TemplateTypeExerciseGenerate
	TemplateTypeMindMapGenerate
	TemplateTypeFileDescribeGenerate
	TemplateTypeNounExplainGenerate
)

type internalMap struct {
	m map[TemplateType]string
}

var Map internalMap

func init() {
	Map = internalMap{m: make(map[TemplateType]string)}
	Map.m[TemplateTypeUserQuery] = userQueryPrompt
	Map.m[TemplateTypeExerciseGenerate] = exercisePrompt
	Map.m[TemplateTypeMindMapGenerate] = mindMapPrompt
	Map.m[TemplateTypeFileDescribeGenerate] = fileDescribePrompt
	Map.m[TemplateTypeNounExplainGenerate] = nounExplainPrompt
}

type TemplateMessage struct {
	Type      TemplateType
	Documents []*schema.Document
	Query     string
	History   []*schema.Message
	File      []*schema.Document
}

func (m internalMap) Get(ctx context.Context, tm TemplateMessage) ([]*schema.Message, error) {

	if tm.Type == TemplateTypeUserQuery {
		temp := prompt.FromMessages(schema.GoTemplate, []schema.MessagesTemplate{
			schema.SystemMessage(m.m[tm.Type]),
			schema.MessagesPlaceholder("history", true),
			schema.UserMessage("{{.question}}"),
		}...)
		return temp.Format(ctx, map[string]any{
			"documents": tm.Documents,
			"file":      tm.File,
			"history":   tm.History,
			"question":  tm.Query,
		})
	}

	temp := prompt.FromMessages(schema.GoTemplate, []schema.MessagesTemplate{
		schema.SystemMessage(m.m[tm.Type]),
		schema.UserMessage("{{.question}}"),
	}...)
	return temp.Format(ctx, map[string]any{
		"documents": tm.Documents,
		"question":  tm.Query,
	})
}

var userQueryPrompt = `
# Role: Student Learning Assistant

# Language: Chinese

- When providing assistance:
  • Be clear and concise
  • Include practical examples when relevant
  • Reference documentation when helpful
  • Suggest improvements or next steps if applicable

- Attention: if files and documents are all provided, explain files first, ignore the documents!

here's documents searched for you:
==== doc start ====
	  {{.documents}}
==== doc end ====

here's files the user send now or before
==== file start ====
		{{.file}}
==== file end ====
`

var exercisePrompt = `
# Role: Student Learning Assistant

# Language: Chinese

- Your Generation must obey the following JSON format:
JSON format:

{
	"question": "问题的内容",
	"sections": {
		"A": "a选项",
		"B": "b选项",
		"C": "c选项",
		"D": "d选项"
	},
	"answer": "本题目的答案",
	"reason": "选择该答案的原因"
}
- You MUST NOT Generate anything else, only generate the JSON content!!!

here's documents provided for you:
==== doc start ====
	  {{.documents}}
==== doc end ====
`

var mindMapPrompt = `
# Role: Student Learning Assistant

# Language: Chinese

- Your Should Generate JSMind format JSON content, the format should be like the following:
{
    "meta": {
        "name": "示例思维导图",
        "author": "学海导航助手",
        "version": "1.0"
    },
    "format": "node_tree",
    "data": {
        "id": "root",
        "topic": "主题",
        "children": [
			{
				"id": "子节点的ID" ,
				"topic": "<UNK>",		
				children: "<UNK>",
			}
			
			...剩余的同样的树的节点
		]
    }
}
- You MUST NOT Generate anything else, only generate the JSON content!!!

here's documents provided for you:
==== doc start ====
	  {{.documents}}
==== doc end ====
`

var fileDescribePrompt = `
# Role: Student Learning Assistant

# Language: Chinese

- You should Describe the File's Content to User
- When describing:
	• Be clear and concise
	• Reference documentation when helpful
	• The describe should be short but detailed

here's the file for you to describe:
==== doc start ====
	  {{.documents}}
==== doc end ====
`

var nounExplainPrompt = `
# Role: Student Learning Assistant

# Language: Chinese

- You should Explain the following Noun to User
- When Explaining the noun:
  • Be clear and concise
  • Include practical examples when relevant
  • Reference documentation when helpful

here's documents searched for you:
==== doc start ====
	  {{.documents}}
==== doc end ====
`
