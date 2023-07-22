package goroot

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
	"time"
)

// Create a template for the proto content.
const protoTemplate = `
// Code generated by GoRoot (github.com/hashemzargari/goroot). DO NOT EDIT.

syntax = "proto3";

package {{ .ProtoPackagePath }}.{{ .ServiceName }}.v{{ .ServiceVersion }};


import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

option go_package = "{{ .GoPackagePath }}/{{ .ServiceName }}/v{{ .ServiceVersion }};v{{ .ServiceVersion }}";
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    title: "{{ .ServiceName }}";
    version: "{{ .ServiceVersion }}.0";
  };
};


{{ range $struct := .Types }}
message {{ getTypeName $struct }} {
	{{ range $fieldIndex, $field := getFieldNamesAndTypes $struct }}
	{{ $field.Name }} {{ convertType $field.Type }} = {{ $field.Index }};
	{{ end }}
}
{{ end }}

service {{ .ServiceName }} {
	{{ range $handler := .Handlers }}
	rpc {{ getHandlerName $handler }} ({{ getHandlerRequestType $handler }}) returns ({{ getHandlerResponseType $handler }}) {
		option (google.api.http) = {
			{{ getHandlerMethod $handler }}: "{{ getHandlerApiRoute $handler }}";	
			body: "*";
		};
	}
	{{ end }}
}

`

// extractStructsFromFields extracts all structs from a given type
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

// isStructInSlice checks if a given struct is in a slice of structs
func isStructInSlice(t reflect.Type, structs []interface{}) bool {
	for _, s := range structs {
		if reflect.TypeOf(s) == t {
			return true
		}
	}
	return false
}

// ExtractAllStructsFromHandlers extracts all structs from a slice of handlers
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

// mustIgnoreStruct checks if a struct must be ignored in the proto generation.
func mustIgnoreStruct(t reflect.Type) bool {
	switch t.Name() {
	case "Time", "Duration":
		return true
	default:
		return false
	}
}

// GenerateProtoContent Generate the proto content using the template and reflection.
func GenerateProtoContent(
	goPackagePath string,
	protoPackagePath string,
	serviceName string,
	serviceVersion uint,
	handlers []Handler,
	allTypes []interface{},
) (string, error) {
	tmpl := template.Must(template.New("protoTemplate").Funcs(template.FuncMap{
		"getTypeName":            getTypeName,
		"getFieldNamesAndTypes":  getFieldNamesAndTypes,
		"convertType":            convertType,
		"getHandlerName":         getHandlerName,
		"getHandlerRequestType":  getHandlerRequestType,
		"getHandlerResponseType": getHandlerResponseType,
		"getHandlerMethod":       getHandlerMethod,
		"getHandlerApiRoute":     getHandlerApiRoute,
	}).Parse(protoTemplate))

	var generatedContent bytes.Buffer
	if err := tmpl.Execute(&generatedContent, struct {
		Types            []interface{}
		GoPackagePath    string
		ProtoPackagePath string
		ServiceName      string
		ServiceVersion   uint
		Handlers         []Handler
	}{
		Types:            allTypes,
		GoPackagePath:    goPackagePath,
		ProtoPackagePath: protoPackagePath,
		ServiceName:      serviceName,
		ServiceVersion:   serviceVersion,
		Handlers:         handlers,
	}); err != nil {
		return "", err
	}

	return generatedContent.String(), nil
}

// convertType converts the go type to proto type.
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

// getHandlerName returns the name of the handler.
func getHandlerName(handler Handler) string {
	return getTypeName(handler)
}

// getHandlerRequestType returns the name of the request type of the handler.
func getHandlerRequestType(handler Handler) string {
	return getTypeName(handler.GetRequestType())
}

// getHandlerResponseType returns the name of the response type of the handler.
func getHandlerResponseType(handler Handler) string {
	return getTypeName(handler.GetResponseType())
}

// getHandlerMethod returns the method of the handler.
func getHandlerMethod(handler Handler) string {
	return strings.ToLower(handler.GetMethod().String())
}

// getHandlerApiRoute returns the api route of the handler.
func getHandlerApiRoute(handler Handler) string {
	return strings.ToLower(handler.GetApiRoute())
}
