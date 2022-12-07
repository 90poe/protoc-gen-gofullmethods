package gofullmethods

import (
	"google.golang.org/protobuf/compiler/protogen"
)

// getGeneratedFileName returns filename for generated code
func getGeneratedFileName(protoFile *protogen.File) string {
	baseFileName := protoFile.GeneratedFilenamePrefix

	return baseFileName + ".fullmethods.pb.go"
}
