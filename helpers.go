package protomask

import "google.golang.org/protobuf/reflect/protoreflect"

type fieldMask struct {
	paths []string
}

func (mask *fieldMask) IsValid(message protoreflect.ProtoMessage) bool {
	if mask == nil || message == nil {
		return false
	}

	ref := message.ProtoReflect()
	for _, path := range mask.paths {
		field := ref.Descriptor().Fields().ByName(protoreflect.Name(path))
		if field == nil {
			return false
		}
		if !ref.Has(field) {
			return false
		}
	}
	return true
}

func (mask *fieldMask) GetPaths() []string {
	if mask == nil {
		return nil
	}
	return mask.paths
}

// All returns a field mask that contains all populated fields from the given message.
// It will only contain shallow fields from the root of the message, no nested ones.
func All[T protoreflect.ProtoMessage](message T) FieldMask {
	paths := []string{}
	message.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		paths = append(paths, string(fd.Name()))
		return true
	})
	return &fieldMask{
		paths: paths,
	}
}
