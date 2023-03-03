package proto

import (
	"log"
	"os"

	"github.com/emicklei/proto"
)

func Parse(filePath string) (Service, error) {

	reader, err := os.Open(filePath)
	if err != nil {
		return Service{}, err
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return Service{}, err
	}

	var protoInterface Service

	proto.Walk(
		definition,
		proto.WithEnum(enumHandler(&protoInterface)),
		proto.WithImport(importHandler(&protoInterface)),
		proto.WithMessage(messageHandler(&protoInterface)),
		proto.WithOneof(oneofHandler(&protoInterface)),
		proto.WithOption(optionHandler(&protoInterface)),
		proto.WithPackage(packageHandler(&protoInterface)),
		proto.WithRPC(rpcHandler(&protoInterface)),
		proto.WithService(serviceHandler(&protoInterface)),
	)

	protoInterface = markNestedMessages(protoInterface)

	return protoInterface, nil
}

func enumHandler(ms *Service) func(*proto.Enum) {

	return func(*proto.Enum) {
	}
}

func importHandler(ms *Service) func(*proto.Import) {

	return func(*proto.Import) {
	}
}

func messageHandler(ms *Service) func(*proto.Message) {

	return func(m *proto.Message) {

		message := Message{
			Name: m.Name,
		}

		for _, f := range m.Elements {
			field, ok := f.(*proto.NormalField)
			if !ok {
				log.Fatal("Only normal fields supported")
			}

			message.Fields = append(message.Fields, Field{
				Name: field.Name,
				Type: field.Type,
			})
		}

		ms.Messages = append(ms.Messages, message)
	}
}

func oneofHandler(ms *Service) func(*proto.Oneof) {

	return func(*proto.Oneof) {

		log.Fatal("OneOf not supported")
	}
}

func optionHandler(ms *Service) func(*proto.Option) {

	return func(*proto.Option) {
	}
}

func packageHandler(ms *Service) func(*proto.Package) {

	return func(p *proto.Package) {

		ms.ProtoPackage = p.Name
	}
}

func rpcHandler(ms *Service) func(*proto.RPC) {

	return func(r *proto.RPC) {

		ms.Methods = append(ms.Methods, Method{
			Name: r.Name,
			RequestMessage: r.RequestType,
			ResponseMessage: r.ReturnsType,
		})
	}
}

func serviceHandler(ms *Service) func(*proto.Service) {

	return func(s *proto.Service) {
		ms.ServiceName = s.Name
	}
}

func markNestedMessages(i Service) Service {

	for iM := range i.Messages {
		for iF := range i.Messages[iM].Fields {
			if i.IsMessage(i.Messages[iM].Fields[iF].Type) {
				i.Messages[iM].Fields[iF].IsNestedMessage = true
			}
		}
	}
	return i
}

