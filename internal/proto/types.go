package proto

import (
	"fmt"
	"sort"
)

type Service struct {

	ServiceName string
	ProtoPackage string
	Methods []Method
	Messages []Message
	//Enums []Enum
}

type Method struct {
	Name string
	RequestMessage string
	ResponseMessage string
}

type Message struct {
	Name string
	Fields []Field
}

type Field struct {
	Name string
	Type string
	IsNestedMessage bool
}

func (s Service) Method(
	methodName string,
) (Method, error) {

	for _, m := range s.Methods {
		if m.Name == methodName {
			return m, nil
		}
	}
	return Method{}, fmt.Errorf("No such method '%s'", methodName)
}

func (s Service) Message(
	messageName string,
) (Message, error) {

	for _, m := range s.Messages {
		if m.Name == messageName {
			return m, nil
		}
	}
	return Message{}, fmt.Errorf("No such message '%s'", messageName)
}

func (s Service) MethodRequestMessage(
	methodName string,
) (Message, error) {

	method, err := s.Method(methodName)
	if err != nil {
		return Message{}, err
	}

	message, err := s.Message(method.RequestMessage)
	if err != nil {
		return Message{}, fmt.Errorf("No such request message for method '%+v'", method)
	}

	return message, nil
}

func (s Service) MethodResponseMessage(
	methodName string,
) (Message, error) {

	method, err := s.Method(methodName)
	if err != nil {
		return Message{}, err
	}

	message, err := s.Message(method.ResponseMessage)
	if err != nil {
		return Message{}, fmt.Errorf("No such response message for method '%+v'", method)
	}

	return message, nil
}

func (s Service) MethodRequestFields(
	methodName string,
) ([]Field, error) {

	message, err := s.MethodRequestMessage(methodName)
	if err != nil {
		return nil, err
	}

	return message.Fields, nil
}

func (s Service) MethodResponseFields(
	methodName string,
) ([]Field, error) {

	message, err := s.MethodResponseMessage(methodName)
	if err != nil {
		return nil, err
	}

	return message.Fields, nil
}

func (s Service) MethodResponseTypes(
	methodName string,
) ([]string, error) {

	message, err := s.MethodResponseMessage(methodName)
	if err != nil {
		return nil, err
	}

	return message.FieldTypes(), nil
}

func (m Message) FieldNames() []string {

	types := make([]string, len(m.Fields))
	for i := range m.Fields {
		types[i] = m.Fields[i].Name
	}
	return types
}

func (m Message) FieldTypes() []string {

	types := make([]string, len(m.Fields))
	for i := range m.Fields {
		types[i] = m.Fields[i].Type
	}
	return types
}

// NestedMessages returns all of the messages that are
// nested within any other message in the proto interface
func (s Service) NestedMessages(
) []Message {

	nestedSet := make(map[string]Message)
	for _, m := range s.Messages {
		for _, f := range m.Fields {
			nestedMessage, err := s.Message(f.Type)
			if err != nil {
				continue // Not a nested type
			}
			nestedSet[nestedMessage.Name] = nestedMessage
		}
	}

	msgs := make([]Message, 0, len(nestedSet))
	for _, m := range nestedSet {
		msgs = append(msgs, m)
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Name < msgs[j].Name
	})

	return msgs
}

func (s Service) IsMessage(
	name string,
) bool {

	for _, m := range s.Messages {
		if m.Name == name {
			return true
		}
	}
	return false
}

func (s Service) MethodUsesNestedMessages(
	methodName string,
) (bool, error) {

	reqFields, err := s.MethodRequestFields(methodName)
	if err != nil {
		return false, err
	}

	for _, f := range reqFields {
		if f.IsNestedMessage {
			return true, nil
		}
	}

	resFields, err := s.MethodResponseFields(methodName)
	if err != nil {
		return false, err
	}

	for _, f := range resFields {
		if f.IsNestedMessage {
			return true, nil
		}
	}

	return false, nil
}
