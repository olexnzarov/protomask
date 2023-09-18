package protomask

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var ErrInvalidFieldMask = errors.New("invalid field mask")
var ErrInvalidPath = errors.New("invalid path")

type FieldMask interface {
	GetPaths() []string
	IsValid(protoreflect.ProtoMessage) bool
}

// Update updates the targetMessage with values from the updateMessage, updateMask specifies which fields need to be updated.
//
// FieldMask should contain field names like in .proto file. Nested paths are supported (e.g. "foo.bar.xyz").
// If nested field's parent is a nil value, it will be initialized with a default value.
// Nested paths inside a map are not supported.
//
// See the following package for a FieldMask implementation: https://google.golang.org/protobuf/types/known/fieldmaskpb
func Update[T protoreflect.ProtoMessage](targetMessage T, updateMessage T, updateMask FieldMask) error {
	if updateMask == nil || !updateMask.IsValid(targetMessage) {
		return ErrInvalidFieldMask
	}

	targetMessageRef, updateMessageRef := targetMessage.ProtoReflect(), updateMessage.ProtoReflect()
	for _, path := range updateMask.GetPaths() {
		target, targetField, err := populateMessageProperty(targetMessageRef, path)
		if err != nil {
			return err
		}
		update, updateField, err := populateMessageProperty(updateMessageRef, path)
		if err != nil {
			return err
		}

		value := update.Get(updateField)

		if isNil(value) {
			target.Clear(targetField)
		} else {
			target.Set(
				targetField,
				value,
			)
		}
	}

	return nil
}

func getFieldByName(message protoreflect.Message, fieldName string) (protoreflect.FieldDescriptor, error) {
	field := message.Descriptor().Fields().ByName(protoreflect.Name(fieldName))
	if field == nil {
		return nil, fmt.Errorf("unknown field: '%s'", fieldName)
	}
	return field, nil
}

func toMessage(value protoreflect.Value) protoreflect.Message {
	switch v := value.Interface().(type) {
	case protoreflect.Message:
		return v
	default:
		return nil
	}
}

func isNil(value protoreflect.Value) bool {
	switch v := value.Interface().(type) {
	case protoreflect.Message:
		return !v.IsValid()
	default:
		return !value.IsValid()
	}
}

func populateMessageProperty(message protoreflect.Message, path string) (protoreflect.Message, protoreflect.FieldDescriptor, error) {
	fields := strings.Split(path, ".")
	switch len(fields) {
	case 0:
		return nil, nil, ErrInvalidPath
	case 1:
		field, err := getFieldByName(message, fields[0])
		if err != nil {
			return nil, nil, err
		}
		return message, field, nil
	}

	for i := 0; i < len(fields)-1; i++ {
		// These nil checks are redundant if we use Google's fieldmaskpb package.
		// But it's better to be safe than sorry and account for other implementations not implementing their "IsValid" function correctly.
		messageField := message.Descriptor().Fields().ByName(protoreflect.Name(fields[i]))
		if messageField == nil {
			return nil, nil, fmt.Errorf("unknown field: '%s'", strings.Join(fields[:i+1], "."))
		}

		nextMessage := toMessage(message.Get(messageField))
		if nextMessage == nil {
			return nil, nil, fmt.Errorf("unsupported nested type: '%s'", strings.Join(fields[:i+1], "."))
		}

		// We need to make sure the message value is not nil.
		// For example, we have a path of "a.b.c". It would be possible to get a "c" of "b" even if the value of "b" is nil.
		// Therefore, before we can set a value of "c", we need to initialize "b", or we'll get a nil pointer exception.
		if !nextMessage.IsValid() {
			value := message.NewField(messageField)
			message.Set(messageField, value)
			nextMessage = toMessage(value)
		}

		message = nextMessage
	}

	field, err := getFieldByName(message, fields[len(fields)-1])
	if err != nil {
		return nil, nil, err
	}
	return message, field, nil
}
