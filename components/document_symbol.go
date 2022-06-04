package components

import (
	"context"
	"pls/proto/view"

	"github.com/TobiasYin/go-lsp/logs"
	"github.com/TobiasYin/go-lsp/lsp/defines"
)

type sessionKeyType struct{}

var sessionKey = sessionKeyType{}

func ProvideDocumentSymbol(ctx context.Context, req *defines.DocumentSymbolParams) (result *[]defines.DocumentSymbol, err error) {
	file, err := view.ViewManager.GetFile(req.TextDocument.Uri)
	res := []defines.DocumentSymbol{}
	if err != nil {
		logs.Printf("GetFile err: %v", err)
		return &res, nil
	}
	for _, pack := range file.Proto().Packages() {
		res = append(res, defines.DocumentSymbol{
			Name: pack.ProtoPackage.Name,
			Kind: defines.SymbolKindPackage,
			SelectionRange: defines.Range{
				Start: defines.Position{Line: uint(pack.ProtoPackage.Position.Line - 1)},
				End:   defines.Position{Line: uint(pack.ProtoPackage.Position.Line - 1)},
			},
		})
	}
	for _, enums := range file.Proto().Enums() {
		res = append(res, defines.DocumentSymbol{
			Name: enums.Protobuf().Name,
			Kind: defines.SymbolKindEnum,
			SelectionRange: defines.Range{
				Start: defines.Position{Line: uint(enums.Protobuf().Position.Line - 1)},
				End:   defines.Position{Line: uint(enums.Protobuf().Position.Line - 1)},
			},
		})
	}
	for _, message := range file.Proto().Messages() {
		message_proto := message.Protobuf()
		res = append(res, defines.DocumentSymbol{
			Name: message_proto.Name,
			Kind: defines.SymbolKindClass,
			SelectionRange: defines.Range{
				Start: defines.Position{Line: uint(message_proto.Position.Line - 1)},
				End:   defines.Position{Line: uint(message_proto.Position.Line - 1)},
			},
		})
	}
	for _, service := range file.Proto().Services() {
		service_sym := defines.DocumentSymbol{
			Name: service.Protobuf().Name,
			Kind: defines.SymbolKindNamespace,
			SelectionRange: defines.Range{
				Start: defines.Position{Line: uint(service.Protobuf().Position.Line - 1)},
				End:   defines.Position{Line: uint(service.Protobuf().Position.Line - 1)},
			},
			Children: &[]defines.DocumentSymbol{},
		}
		child := []defines.DocumentSymbol{}
		for _, rpc := range service.RPCs() {
			rpc := defines.DocumentSymbol{
				Name: rpc.ProtoRPC.Name,
				Kind: defines.SymbolKindMethod,
				SelectionRange: defines.Range{
					Start: defines.Position{Line: uint(rpc.ProtoRPC.Position.Line - 1)},
					End:   defines.Position{Line: uint(rpc.ProtoRPC.Position.Line - 1)},
				},
			}
			child = append(child, rpc)
		}
		service_sym.Children = &child
		res = append(res, service_sym)
	}
	return &res, nil
}
