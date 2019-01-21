package schema

import (
	. "github.com/dave/jennifer/jen"
	. "go-restli/codegen"
)

func (r *Resource) GenerateCode(packagePrefix string) *CodeFile {
	c := &CodeFile{
		PackagePath: r.PackagePath(packagePrefix),
		Filename:    ExportedIdentifier(r.Name),
		Code:        Empty(),
	}

	// WIP
	//r.generateClient(packagePrefix, c.Code)

	for _, s := range r.generateAllActionStructs(packagePrefix) {
		c.Code.Add(s).Line()
	}

	return c
}
