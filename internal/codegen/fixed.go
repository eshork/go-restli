package codegen

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

type Fixed struct {
	NamedType
	Size int
}

func (f *Fixed) InnerTypes() IdentifierSet {
	return nil
}

func (f *Fixed) GenerateCode() (def *Statement) {
	def = Empty()
	AddWordWrappedComment(def, f.Doc).Line()
	def.Type().Id(f.Name).Index(Lit(f.Size)).Byte().Line().Line()

	receiver := ReceiverName(f.Name)
	errorMsg := fmt.Sprintf("size of %s must be exactly %d bytes (was %%d)", f.Name, f.Size)

	AddUnmarshalJSON(def, receiver, f.Name, func(def *Group) {
		def.Id("bytes").Op(":=").Make(Bytes(), Lit(f.Size))
		def.Err().Op("=").Id("bytes").Dot(UnmarshalJSON).Call(Id("data"))
		IfErrReturn(def)
		def.If(Len(Id("bytes")).Op("!=").Lit(f.Size)).BlockFunc(func(def *Group) {
			def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(errorMsg), Len(Id("bytes")))
			def.Return()
		})
		def.Copy(Id(receiver).Index(Op(":")), Id("bytes").Index(Op(":").Lit(f.Size)))
		def.Return()
	}).Line().Line()

	AddRestLiEncode(def, receiver, f.Name, func(def *Group) {
		def.Add(Encoder).Dot("Bytes").Call(Id(receiver).Index(Op(":")))
		def.Return(Nil())
	}).Line().Line()
	AddRestLiDecode(def, receiver, f.Name, func(def *Group) {
		def.Id("bytes").Op(":=").Make(Bytes(), Lit(f.Size))
		def.Err().Op("=").Id(Codec).Dot("DecodeBytes").Call(Id("data"), Op("&").Id("bytes"))
		IfErrReturn(def)
		def.If(Len(Id("bytes")).Op("!=").Lit(f.Size)).BlockFunc(func(def *Group) {
			def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(errorMsg), Len(Id("bytes")))
			def.Return()
		})
		def.Copy(Id(receiver).Index(Op(":")), Id("bytes").Index(Op(":").Lit(f.Size)))
		def.Return()
	}).Line().Line()

	return def
}

func (f *Fixed) getLit(rawJson string) *Statement {
	var v protocol.Bytes
	if err := (&v).UnmarshalJSON([]byte(rawJson)); err != nil {
		Logger.Panicf("(%+v) Illegal primitive literal: \"%s\" (%s)", f, rawJson, err)
	}

	if f.Size != len(v) {
		Logger.Panicf("(%+v) Default value %q does not have %d bytes (got %d bytes)", f, rawJson, f.Size, len(v))
	}

	return f.Qual().ValuesFunc(func(def *Group) {
		for _, b := range v {
			def.LitByte(b)
		}
	})
}
