package iamcel

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func collectDependencies(messages ...protoreflect.MessageDescriptor) (*protoregistry.Files, error) {
	fileDescriptors := make(map[string]protoreflect.FileDescriptor, len(messages))
	for _, message := range messages {
		messageParentFile := message.ParentFile()
		fileDescriptors[messageParentFile.Path()] = messageParentFile
		// Initialize list of transitive imports.
		transitiveImports := make([]protoreflect.FileImport, 0, messageParentFile.Imports().Len())
		for i := 0; i < messageParentFile.Imports().Len(); i++ {
			transitiveImports = append(transitiveImports, messageParentFile.Imports().Get(i))
		}
		// Expand list of transitive imports.
		for i := 0; i < len(transitiveImports); i++ {
			currImport := transitiveImports[i]
			if _, ok := fileDescriptors[currImport.Path()]; ok {
				continue
			}
			fileDescriptors[currImport.Path()] = currImport.FileDescriptor
			for j := 0; j < currImport.FileDescriptor.Imports().Len(); j++ {
				transitiveImports = append(transitiveImports, currImport.FileDescriptor.Imports().Get(j))
			}
		}
	}
	var registry protoregistry.Files
	for _, fileDescriptor := range fileDescriptors {
		if err := registry.RegisterFile(fileDescriptor); err != nil {
			return nil, err
		}
	}
	return &registry, nil
}
