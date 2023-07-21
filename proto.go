package goroot

import (
	"bytes"
	"os"
	"reflect"
	"text/template"
	"time"
)

// Create a template for the proto content.
const protoTemplate = `
syntax = "proto3";

{{ range $struct := .Types }}
message {{ getTypeName $struct }} {
	{{ range $fieldIndex, $field := getFieldNamesAndTypes $struct }}
	{{ $field.Name }} {{ convertType $field.Type }} = {{ $field.Index }};
	{{ end }}
}
{{ end }}

service {{ .ServiceName }} {
	{{ range $handler := .Handlers }}
	rpc {{ getHandlerName $handler }} ({{ getHandlerRequestType $handler }}) returns ({{ getHandlerResponseType $handler }});
	{{ end }}
}

`

func extractStructsFromFields(t reflect.Type) []any {
	structs := make([]any, 0)
	if t.Kind() != reflect.Struct {
		return structs
	}

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i).Type

		if fieldType.Kind() == reflect.Struct {
			structs = append(structs, reflect.New(fieldType).Elem().Interface())
			childStructs := extractStructsFromFields(fieldType)
			if len(childStructs) > 0 {
				structs = append(structs, childStructs...)
			}
		} else if fieldType.Kind() == reflect.Array || fieldType.Kind() == reflect.Slice {
			if fieldType.Elem().Kind() == reflect.Struct {
				structs = append(structs, reflect.New(fieldType.Elem()).Elem().Interface())
				childStructs := extractStructsFromFields(fieldType.Elem())
				if len(childStructs) > 0 {
					structs = append(structs, childStructs...)
				}
			}
		} else if fieldType.Kind() == reflect.Map {
			if fieldType.Elem().Kind() == reflect.Struct {
				structs = append(structs, reflect.New(fieldType.Elem()).Elem().Interface())
				childStructs := extractStructsFromFields(fieldType.Elem())
				if len(childStructs) > 0 {
					structs = append(structs, childStructs...)
				}
			}
		}
	}

	return structs
}

func isStructInSlice(t reflect.Type, structs []interface{}) bool {
	for _, s := range structs {
		if reflect.TypeOf(s) == t {
			return true
		}
	}
	return false
}

func ExtractAllStructsFromHandlers(handlers []Handler) []any {
	var structs []any
	for _, handler := range handlers {
		requestStruct := handler.GetRequestType()
		responseStruct := handler.GetResponseType()

		// extract fields of structs and if they are structs, add them to structs
		// if they are not structs, ignore them
		requestStructs := extractStructsFromFields(reflect.TypeOf(requestStruct))
		for _, s := range requestStructs {
			if !isStructInSlice(reflect.TypeOf(s), structs) {
				structs = append(structs, s)
			}
		}

		responseStructs := extractStructsFromFields(reflect.TypeOf(responseStruct))
		for _, s := range responseStructs {
			if !isStructInSlice(reflect.TypeOf(s), structs) && !mustIgnoreStruct(reflect.TypeOf(s)) {
				structs = append(structs, s)
			}
		}

		// add the request and response structs to structs
		structs = append(structs, requestStruct, responseStruct)
	}
	return structs
}

func mustIgnoreStruct(t reflect.Type) bool {
	switch t.Name() {
	case "Time", "Duration":
		return true
	default:
		return false
	}
}

// GenerateProtoContent Generate the proto content using the template and reflection.
func GenerateProtoContent(serviceName string, handlers []Handler, allTypes []interface{}) (string, error) {
	tmpl := template.Must(template.New("protoTemplate").Funcs(template.FuncMap{
		"getTypeName":            getTypeName,
		"getFieldNamesAndTypes":  getFieldNamesAndTypes,
		"convertType":            convertType,
		"getHandlerName":         getHandlerName,
		"getHandlerRequestType":  getHandlerRequestType,
		"getHandlerResponseType": getHandlerResponseType,
	}).Parse(protoTemplate))

	var generatedContent bytes.Buffer
	if err := tmpl.Execute(&generatedContent, struct {
		Types       []interface{}
		ServiceName string
		Handlers    []Handler
	}{
		Types:       allTypes,
		ServiceName: serviceName,
		Handlers:    handlers,
	}); err != nil {
		return "", err
	}

	return generatedContent.String(), nil
}

// Write the content to a file.
func writeToFile(fileName, content string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions for reflection.
func getTypeName(i interface{}) string {
	reflectType := reflect.TypeOf(i)
	if reflectType.Kind() == reflect.Array || reflectType.Kind() == reflect.Slice {
		return getTypeName(reflect.TypeOf(i).Elem())
	}
	if reflectType.Kind() == reflect.Map {
		return getTypeName(reflect.TypeOf(i).Elem())
	}
	return reflectType.Name()
}

func getFieldNamesAndTypes(i interface{}) []struct {
	Name  string
	Type  reflect.Type
	Index int
} {
	var fields []struct {
		Name  string
		Type  reflect.Type
		Index int
	}
	t := reflect.TypeOf(i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fields = append(fields, struct {
			Name  string
			Type  reflect.Type
			Index int
		}{
			field.Name,
			field.Type,
			i + 1,
		})
	}

	return fields
}

func convertType(goType reflect.Type) string {

	if goType.Kind() == reflect.Array || goType.Kind() == reflect.Slice {
		return "repeated " + convertType(goType.Elem())
	}

	if goType.Kind() == reflect.Map {
		keyType := goType.Key()
		valueType := goType.Elem()
		return "map<" + convertType(keyType) + ", " + convertType(valueType) + ">"
	}

	switch goType {
	case reflect.TypeOf(""):
		return "string"
	case reflect.TypeOf(int(0)):
		return "int32"
	case reflect.TypeOf(int8(0)):
		return "int32"
	case reflect.TypeOf(int16(0)):
		return "int32"
	case reflect.TypeOf(int32(0)):
		return "int32"
	case reflect.TypeOf(int64(0)):
		return "int64"
	case reflect.TypeOf(uint(0)):
		return "uint32"
	case reflect.TypeOf(uint8(0)):
		return "uint32"
	case reflect.TypeOf(uint16(0)):
		return "uint32"
	case reflect.TypeOf(uint32(0)):
		return "uint32"
	case reflect.TypeOf(uint64(0)):
		return "uint64"
	case reflect.TypeOf(float32(0)):
		return "float"
	case reflect.TypeOf(float64(0)):
		return "double"
	case reflect.TypeOf(bool(false)):
		return "bool"
	case reflect.TypeOf(time.Time{}):
		return "google.protobuf.Timestamp"
	default:
		return goType.Name()
	}
}

func getHandlerName(handler Handler) string {
	return getTypeName(handler)
}

func getHandlerRequestType(handler Handler) string {
	return getTypeName(handler.GetRequestType())
}

func getHandlerResponseType(handler Handler) string {
	return getTypeName(handler.GetResponseType())
}