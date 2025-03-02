package translator

import "github.com/invopop/jsonschema"

type TranslateResponse struct {
	Markdown string `json:"markdown" jsonschema_description:"translated result"`
}

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

var TranslateMarkdownSchema = GenerateSchema[TranslateResponse]
