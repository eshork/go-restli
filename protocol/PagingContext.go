/*
DO NOT EDIT

Code automatically generated by github.com/PapaCharlie/go-restli
Source file: https://github.com/PapaCharlie/go-restli/blob/master/internal/codegen/resources/pagingcontext.go
*/

package protocol

import (
	fnv1a "github.com/PapaCharlie/go-restli/fnv1a"
	restlicodec "github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

type PagingContext struct {
	Start *int32 `json:"start,omitempty"`

	Count *int32 `json:"count,omitempty"`
}

func (p *PagingContext) EqualsInterface(otherInterface interface{}) bool {
	other, ok := otherInterface.(*PagingContext)
	if !ok {
		return false
	}

	return p.Equals(other)
}

func (p *PagingContext) Equals(other *PagingContext) bool {
	if p == nil || other == nil {
		return false
	}

	if p.Count != other.Count {
		if p.Count == nil || other.Count == nil {
			return false
		}
		if *p.Count != *other.Count {
			return false
		}
	}

	if p.Start != other.Start {
		if p.Start == nil || other.Start == nil {
			return false
		}
		if *p.Start != *other.Start {
			return false
		}
	}

	return true
}

func (p *PagingContext) ComputeHash() (hash fnv1a.Hash) {
	if p == nil {
		return hash
	}
	hash = fnv1a.NewHash()

	if p.Count != nil {
		hash.AddInt32(*p.Count)
	}

	if p.Start != nil {
		hash.AddInt32(*p.Start)
	}

	return hash
}

func (p *PagingContext) EncodeQueryParams() (rawQuery string, err error) {
	writer := restlicodec.NewRestLiQueryParamsWriter()
	err = writer.WriteParams(func(paramNameWriter func(string) restlicodec.Writer) (err error) {
		if p.Count != nil {
			paramNameWriter("count").WriteInt32(*p.Count)
		}

		if p.Start != nil {
			paramNameWriter("start").WriteInt32(*p.Start)
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	return writer.Finalize(), nil
}
