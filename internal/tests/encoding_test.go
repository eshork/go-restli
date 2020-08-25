package tests

import (
	"encoding/json"
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/stretchr/testify/require"
)

func TestEncodePrimitives(t *testing.T) {
	expected := &Primitives{
		PrimitiveInteger: 1,
		PrimitiveLong:    23,
		PrimitiveFloat:   52.5,
		PrimitiveDouble:  66.5,
		PrimitiveBytes:   []byte("@ABC🕴" + string([]byte{1})),
		PrimitiveString:  `a string,()'`,
	}

	t.Run("json", func(t *testing.T) {
		testJsonEncoding(t, expected, new(Primitives), `{
  "primitiveInteger": 1,
  "primitiveLong": 23,
  "primitiveFloat": 52.5,
  "primitiveDouble": 66.5,
  "primitiveBytes": "@ABC🕴\u0001",
  "primitiveString": "a string,()'"
}`)
	})

	testRestliEncoding(t, expected, new(Primitives),
		`(primitiveBytes:@ABC🕴,primitiveDouble:66.5,primitiveFloat:52.5,primitiveInteger:1,primitiveLong:23,primitiveString:a string%2C%28%29%27)`,
	)
}

func TestEncodeComplexTypes(t *testing.T) {
	integer := int32(5)
	hello := "Hello"
	expected := &ComplexTypes{
		ArrayOfMaps: ArrayOfMaps{
			ArrayOfMaps: []map[string]int32{
				{
					"one": 1,
				},
				{
					"two": 2,
				},
			},
		},
		MapOfInts: MapOfInts{
			MapOfInts: map[string]int32{
				"one": 1,
			},
		},
		RecordWithProps: RecordWithProps{
			Integer: &integer,
		},
		UnionOfComplexTypes: UnionOfComplexTypes{
			ComplexTypeUnion: UnionOfComplexTypes_ComplexTypeUnion{
				Fruits: conflictresolution.Fruits_ORANGE.Pointer(),
			},
		},
		UnionOfPrimitives: UnionOfPrimitives{
			PrimitivesUnion: UnionOfPrimitives_PrimitivesUnion{
				Int: &integer,
			},
		},
		AnotherUnionOfComplexTypes: UnionOfComplexTypes{
			ComplexTypeUnion: UnionOfComplexTypes_ComplexTypeUnion{
				Fruits: conflictresolution.Fruits_APPLE.Pointer(),
			},
		},
		UnionOfSameTypes: UnionOfSameTypes{
			SameTypesUnion: UnionOfSameTypes_SameTypesUnion{
				Greeting: &hello,
			},
			UnionWithArrayMembers: UnionOfSameTypes_UnionWithArrayMembers{
				FruitArray: &[]conflictresolution.Fruits{
					conflictresolution.Fruits_ORANGE,
					conflictresolution.Fruits_APPLE,
				},
			},
			UnionWithMapMembers: UnionOfSameTypes_UnionWithMapMembers{
				IntMap: &map[string]int32{
					"one": 1,
				},
			},
		},
	}

	t.Run("json", func(t *testing.T) {
		testJsonEncoding(t, expected, new(ComplexTypes), `{
  "arrayOfMaps": {
    "arrayOfMaps": [
      {
        "one": 1
      },
      {
        "two": 2
      }
    ]
  },
  "mapOfInts": {
    "mapOfInts": {
      "one": 1
    }
  },
  "recordWithProps": {
    "integer": 5
  },
  "unionOfComplexTypes": {
    "complexTypeUnion": {
      "testsuite.Fruits": "ORANGE"
    }
  },
  "unionOfPrimitives": {
    "primitivesUnion": {
      "int": 5
    }
  },
  "anotherUnionOfComplexTypes": {
    "complexTypeUnion": {
      "testsuite.Fruits": "APPLE"
    }
  },
  "unionOfSameTypes": {
    "sameTypesUnion": {
      "greeting": "Hello"
    },
    "unionWithArrayMembers": {
      "fruitArray": [
        "ORANGE",
        "APPLE"
      ]
    },
    "unionWithMapMembers": {
      "intMap": {
        "one": 1
      }
    }
  }
}`)
	})

	testRestliEncoding(t, expected, new(ComplexTypes),
		`(anotherUnionOfComplexTypes:(complexTypeUnion:(testsuite.Fruits:APPLE)),`+
			`arrayOfMaps:(arrayOfMaps:List((one:1),(two:2))),`+
			`mapOfInts:(mapOfInts:(one:1)),`+
			`recordWithProps:(integer:5),`+
			`unionOfComplexTypes:(complexTypeUnion:(testsuite.Fruits:ORANGE)),`+
			`unionOfPrimitives:(primitivesUnion:(int:5)),`+
			`unionOfSameTypes:(sameTypesUnion:(greeting:Hello),unionWithArrayMembers:(fruitArray:List(ORANGE,APPLE)),unionWithMapMembers:(intMap:(one:1))))`,
	)
}

func TestUnknownFieldReads(t *testing.T) {
	id := int64(1)
	expected := conflictresolution.Message{
		Id:      &id,
		Message: "test",
	}

	tests := []struct {
		Name   string
		Json   string
		RestLi string
		Actual conflictresolution.Message
	}{
		{
			Name:   "Extra field before",
			Json:   `{"foo":false,"id":1,"message":"test"}`,
			RestLi: `(foo:false,id:1,message:test)`,
		},
		{
			Name:   "Extra field in the middle",
			Json:   `{"id":1,"foo":false,"message":"test"}`,
			RestLi: `(id:1,foo:false,message:test)`,
		},
		{
			Name:   "Extra field at the end",
			Json:   `{"id":1,"message":"test","foo":false}`,
			RestLi: `(id:1,message:test,foo:false)`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				reader := restlicodec.NewJsonReader([]byte(test.Json))
				require.NoError(t, test.Actual.UnmarshalRestLi(reader))
				require.Equal(t, expected, test.Actual)
			})

			t.Run("url", func(t *testing.T) {
				reader := restlicodec.NewHeaderReader(test.RestLi)
				require.NoError(t, test.Actual.UnmarshalRestLi(reader))
				require.Equal(t, expected, test.Actual)
			})
		})
	}
}

func TestReadIncorrectType(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		var actual conflictresolution.Message
		reader := restlicodec.NewJsonReader([]byte(`{"message": [1]}`))
		require.Error(t, actual.UnmarshalRestLi(reader))
	})

	t.Run("url", func(t *testing.T) {
		var actual conflictresolution.Message
		reader := restlicodec.NewHeaderReader(`(message:List(1))`)
		require.Error(t, actual.UnmarshalRestLi(reader))
	})
}

func TestMissingRequiredFields(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		var actual conflictresolution.Message
		reader := restlicodec.NewJsonReader([]byte(`{"id":1}`))
		require.Error(t, actual.UnmarshalRestLi(reader))
	})

	t.Run("url", func(t *testing.T) {
		var actual conflictresolution.Message
		reader := restlicodec.NewHeaderReader(`(id:1)`)
		require.Error(t, actual.UnmarshalRestLi(reader))
	})
}

func TestMapEncoding(t *testing.T) {
	expected := &Optionals{
		OptionalMap: &map[string]int32{
			"one": 1,
			"two": 2,
		},
	}

	t.Run("multipleElements", func(t *testing.T) {
		writer := restlicodec.NewHeaderWriter()
		require.NoError(t, expected.MarshalRestLi(writer))

		serialized := writer.Finalize()
		if serialized != `(optionalMap:(one:1,two:2))` && serialized != `(optionalMap:(two:2,one:1))` {
			t.Fail()
		}
	})

	expected = &Optionals{OptionalMap: &map[string]int32{}}
	writer := restlicodec.NewHeaderWriter()
	require.NoError(t, expected.MarshalRestLi(writer))
	require.Equal(t, `(optionalMap:())`, writer.Finalize())
}

func TestArrayEncoding(t *testing.T) {
	expected := &Optionals{OptionalArray: &[]int32{1, 2}}

	testRestliEncoding(t, expected, new(Optionals), `(optionalArray:List(1,2))`)

	expected = &Optionals{OptionalArray: &[]int32{}}
	writer := restlicodec.NewHeaderWriter()
	require.NoError(t, expected.MarshalRestLi(writer))
	require.Equal(t, `(optionalArray:List())`, writer.Finalize())
}

func TestEmptyStringAndBytes(t *testing.T) {
	expected := &Optionals{
		OptionalBytes:  new([]byte),
		OptionalString: new(string),
	}

	testRestliEncoding(t, expected, new(Optionals), `(optionalBytes:'',optionalString:'')`)
}

type restliObject interface {
	restlicodec.Marshaler
	restlicodec.Unmarshaler
}

func testJsonEncoding(t *testing.T, expected, actual restliObject, expectedRawJson string) {
	t.Run("encode", func(t *testing.T) {
		writer := restlicodec.NewCompactJsonWriter()
		require.NoError(t, expected.MarshalRestLi(writer))

		var expectedRaw map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(expectedRawJson), &expectedRaw))
		var raw map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(writer.Finalize()), &raw))
		require.Equal(t, expectedRaw, raw)
	})

	t.Run("decode", func(t *testing.T) {
		decoder := restlicodec.NewJsonReader([]byte(expectedRawJson))
		require.NoError(t, actual.UnmarshalRestLi(decoder))
		require.Equal(t, expected, actual)
	})
}

func testRestliEncoding(t *testing.T, expected, actual restliObject, expectedRawEncoded string) {
	t.Run("encode", func(t *testing.T) {
		writer := restlicodec.NewHeaderWriter()
		require.NoError(t, expected.MarshalRestLi(writer))
		require.Equal(t, expectedRawEncoded, writer.Finalize())
	})

	t.Run("decode", func(t *testing.T) {
		reader := restlicodec.NewHeaderReader(expectedRawEncoded)
		require.NoError(t, actual.UnmarshalRestLi(reader))
		require.Equal(t, expected, actual)
	})
}
