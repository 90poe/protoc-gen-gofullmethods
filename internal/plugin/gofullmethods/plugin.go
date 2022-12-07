package gofullmethods

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	// for excluding google proto files
	googleProtoFilePrefix = "google"
)

// Plugin represents gofullmethods plugin
type Plugin struct {
	input  io.Reader
	output io.Writer
}

// NewPlugin creates new gofullmethods plugin instance
func NewPlugin(input io.Reader, output io.Writer) *Plugin {
	return &Plugin{
		input:  input,
		output: output,
	}
}

// getRequest retrieves data from provided io.Reader and returns plugin reques object
func (p *Plugin) getRequest() (*pluginpb.CodeGeneratorRequest, error) {
	input, err := io.ReadAll(p.input)
	if err != nil {
		return nil, fmt.Errorf("cannot read plugin input: %w", err)
	}

	var request pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(input, &request); err != nil {
		return nil, fmt.Errorf("cannot unmrashal proto code generator request: %w", err)
	}

	return &request, nil
}

// getDataObject converts protogen file to data object
func (p *Plugin) getDataObject(protoFile *protogen.File) data {
	services := make([]service, 0, len(protoFile.Services))
	for _, grpcService := range protoFile.Services {
		grpcMethods := make([]string, 0, len(grpcService.Methods))
		for _, grpcMethod := range grpcService.Methods {
			grpcMethods = append(grpcMethods, grpcMethod.GoName)
		}

		services = append(services, service{
			Name:    grpcService.GoName,
			Methods: grpcMethods,
		})
	}

	return data{
		FileName:  *protoFile.Proto.Name,
		GoPackage: string(protoFile.GoPackageName),
		Package:   *protoFile.Proto.Package,
		Services:  services,
	}
}

// processProtoFile writes generated golang code to protogen plugin
func (p *Plugin) processProtoFile(protogenPlugin *protogen.Plugin, protoFile *protogen.File) error {
	// skip google proto files
	if strings.HasPrefix(*protoFile.Proto.Name, googleProtoFilePrefix) {
		return nil
	}

	// skip if protofile without services
	if len(protoFile.Services) == 0 {
		return nil
	}

	// prepare golang code
	dataObj := p.getDataObject(protoFile)
	buffer := &bytes.Buffer{}
	if err := golangCodeTemplate.Execute(buffer, &dataObj); err != nil {
		return fmt.Errorf("cannot execute golang template: %w", err)
	}

	golangCodeBytes := buffer.Bytes()
	formattedCode, err := format.Source(golangCodeBytes)
	if err != nil {
		return fmt.Errorf("cannot format generated go code: %w", err)
	}

	// protogenPlugin has slice of files
	generatedFile := protogenPlugin.NewGeneratedFile(getGeneratedFileName(protoFile), ".")
	_, err = generatedFile.Write(formattedCode)
	if err != nil {
		return fmt.Errorf("cannot write generated code to protogen plugin: %w", err)
	}

	return nil
}

func supportedCodeGeneratorFeatures() uint64 {
	// Enable support for optional keyword in proto3.
	return uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
}

// SetSupportedFeaturesOnPluginGen sets supported proto3 features
// on protogen.Plugin.
func SetSupportedFeaturesOnPluginGen(gen *protogen.Plugin) {
	gen.SupportedFeatures = supportedCodeGeneratorFeatures()
}

// SetSupportedFeaturesOnCodeGeneratorResponse sets supported proto3 features
// on pluginpb.CodeGeneratorResponse.
func SetSupportedFeaturesOnCodeGeneratorResponse(resp *pluginpb.CodeGeneratorResponse) {
	sf := supportedCodeGeneratorFeatures()
	resp.SupportedFeatures = &sf
}

// Run executes gofullmethods plugin
func (p *Plugin) Run() error {
	// get plugin request
	request, err := p.getRequest()
	if err != nil {
		return err
	}

	// default options
	protogenOptions := protogen.Options{}
	protogenPlugin, err := protogenOptions.New(request)
	if err != nil {
		return fmt.Errorf("cannot create protogen plugin: %w", err)
	}
	SetSupportedFeaturesOnPluginGen(protogenPlugin)

	// generate golang code and write it to protogen plugin
	for _, protoFile := range protogenPlugin.Files {
		if processErr := p.processProtoFile(protogenPlugin, protoFile); processErr != nil {
			return processErr
		}
	}

	// write response to output
	protogenPluginOutput := protogenPlugin.Response()
	SetSupportedFeaturesOnCodeGeneratorResponse(protogenPluginOutput)

	out, err := proto.Marshal(protogenPluginOutput)
	if err != nil {
		return fmt.Errorf("cannot marshal protogen plugin output: %w", err)
	}

	fmt.Fprint(p.output, string(out))
	return nil
}
