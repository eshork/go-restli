package types

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type ComplexKey struct {
	NamedType
	Key    utils.Identifier
	Params utils.Identifier
}

func (ck *ComplexKey) InnerTypes() utils.IdentifierSet {
	return utils.NewIdentifierSet(ck.Key, ck.Params)
}

func (ck *ComplexKey) GenerateCode() *Statement {
	def := utils.AddWordWrappedComment(Empty(), ck.Doc).Line().
		Type().Id(ck.Name).
		StructFunc(func(def *Group) {
			def.Add(ck.Key.Qual())
			def.Id(ComplexKeyParams).Op("*").Add(ck.Params.Qual()).Tag(utils.JsonFieldTag("$params", false))
		}).Line().Line()

	record := &Record{
		NamedType: ck.NamedType,
		Fields:    utils.TypeRegistry.Resolve(ck.Key).(*Record).Fields,
	}

	AddMarshalRestLi(def, record.Receiver(), ck.Name, func(def *Group) {
		record.generateMarshaler(def, Id(record.Receiver()).Dot(ck.Key.Name))
	})

	AddUnmarshalRestli(def, record.Receiver(), ck.Name, func(def *Group) {
		record.generateUnmarshaler(def, Id(record.Receiver()).Dot(ck.Key.Name), &ck.Params)
	})

	return def
}