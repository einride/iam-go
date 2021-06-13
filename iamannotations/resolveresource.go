package iamannotations

import (
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func resolveResource(
	files *protoregistry.Files,
	startPackage protoreflect.FullName,
	resourceType string,
) (*annotations.ResourceDescriptor, bool) {
	var result *annotations.ResourceDescriptor
	var searchMessagesFn func(protoreflect.MessageDescriptors) bool
	searchMessagesFn = func(messages protoreflect.MessageDescriptors) bool {
		for i := 0; i < messages.Len(); i++ {
			message := messages.Get(i)
			if resource := proto.GetExtension(
				message.Options(), annotations.E_Resource,
			).(*annotations.ResourceDescriptor); resource != nil {
				if resource.Type == resourceType {
					result = resource
					return false
				}
			}
			if !searchMessagesFn(message.Messages()) {
				return false
			}
		}
		return true
	}
	searchFileFn := func(file protoreflect.FileDescriptor) bool {
		// Search file annotations.
		for _, resource := range proto.GetExtension(
			file.Options(), annotations.E_ResourceDefinition,
		).([]*annotations.ResourceDescriptor) {
			if resource.Type == resourceType {
				result = resource
				return false
			}
		}
		return searchMessagesFn(file.Messages())
	}
	// Start with a narrow search in the same package.
	files.RangeFilesByPackage(startPackage, searchFileFn)
	if result != nil {
		return result, true
	}
	// Fall back to a broad search of all files.
	files.RangeFiles(searchFileFn)
	return result, result != nil
}
