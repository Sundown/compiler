package subtle

import (
	"sundown/solution/oversight"

	"github.com/llir/llvm/ir/types"
)

func (t *Type) AsLLType() types.Type {
	if t.Atomic != nil {
		// Type already calculated
		return t.LLType
	} else if t.Vector != nil {
		// Recurse until atomic type(s) found
		// Vectors are always of the form <length | capacity | *data>
		return types.NewStruct(
			types.I64,                             // length
			types.I64,                             // capacity
			types.NewPointer(t.Vector.AsLLType())) // *data
	} else if t.Tuple != nil {
		// Recurse each item in tuple
		var lltypes []types.Type
		for _, t := range t.Tuple {
			lltypes = append(lltypes, t.AsLLType())
		}

		return types.NewStruct(lltypes...)
	} else {
		panic("Type is empty")
	}
}

// Used for calloc'ing vectors
func (t *Type) WidthInBytes() int64 {
	if t.Atomic != nil {
		return t.Width
	} else if t.Vector != nil {
		return 24
	} else if t.Tuple != nil {
		var sum int64
		for _, t := range t.Tuple {
			sum += t.WidthInBytes()
		}
		return sum
	} else {
		oversight.Warn("Using 32 bytes for unknown type " + t.String())
		return 8
	}
}

func (e *Expression) Type() *Type {
	if e.Morpheme != nil {
		return e.Morpheme.TypeOf
	} else if a := e.Application; a != nil {
		// Implement T -> T transform
		// ([T], Int) -> T i.e. ([Char], Int) -> Char
		/*if a.Argument.TypeOf.Atomic != nil && *a.Argument.TypeOf.Atomic == "T" {
			// Poly
		} else {

		}*/

		//lmfao this is completely broken
	}

	return nil
}