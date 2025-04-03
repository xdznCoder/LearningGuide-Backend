package api

import (
	ChatProto "LearningGuide/file_api/proto/.ChatProto"
	"encoding/json"
	"google.golang.org/grpc"
	"io"
	"regexp"
	"strings"
)

type TemplateType int

const (
	TemplateTypeUserQuery TemplateType = iota + 1
	TemplateTypeExerciseGenerate
	TemplateTypeMindMapGenerate
	TemplateTypeFileDescribeGenerate
	TemplateTypeNounExplainGenerate
)

func transResultToStringJSON(result string) string {
	pattern := "```json\\s*({[\\s\\S]*?})\\s*```"

	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(result)

	var output string

	if matches == nil {
		output = result
	} else {
		output = matches[1]
	}

	return output
}

func transResultToExercise(result string) (exercise, error) {
	pattern := "```json\\s*({[\\s\\S]*?})\\s*```"

	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(result)

	var output string

	if matches == nil {
		output = result
	} else {
		output = matches[1]
	}

	var question exercise

	err := json.Unmarshal([]byte(output), &question)
	if err != nil {
		return exercise{}, err
	}

	return question, nil
}

type exercise struct {
	Question string      `json:"question"`
	Sections SectionsSet `json:"sections"`
	Answer   string      `json:"answer"`
	Reason   string      `json:"reason"`
}

type SectionsSet struct {
	A string `json:"A"`
	B string `json:"B"`
	C string `json:"C"`
	D string `json:"D"`
}

func ToString(stream grpc.ServerStreamingClient[ChatProto.ChatModelResponse]) (string, error) {
	var result strings.Builder
	for {
		output, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		result.WriteString(output.Content)
	}
	return result.String(), nil
}
