package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env Environment) AnalyseMorphemes(ms *palisade.Morpheme) prism.Expression {
	mor := env.AnalyseMorpheme(ms)
	if vec, ok := mor.(prism.Vector); ok {
		if len(*vec.Body) == 1 {
			return (*vec.Body)[0]
		}
	}

	return mor
}

func (env Environment) AnalyseMorpheme(m *palisade.Morpheme) prism.Expression {
	switch {
	case m.Char != nil:
		vec := make([]prism.Expression, len(*m.Char))
		for i, c := range *m.Char {
			vec[i] = prism.Char{Value: string(c[0])}
		}

		return prism.Vector{
			ElementType: prism.VectorType{Type: prism.CharType},
			Body:        &vec,
		}
	case m.Int != nil:
		vec := make([]prism.Expression, len(*m.Int))
		for i, c := range *m.Int {
			vec[i] = prism.Int{Value: c}
		}

		return prism.Vector{
			ElementType: prism.VectorType{Type: prism.IntType},
			Body:        &vec,
		}
	case m.Real != nil:
		vec := make([]prism.Expression, len(*m.Real))
		for i, c := range *m.Real {
			vec[i] = prism.Real{Value: c}
		}

		return prism.Vector{
			ElementType: prism.VectorType{Type: prism.RealType},
			Body:        &vec,
		}
	case m.String != nil:
		vec := make([]prism.Expression, len(*m.String))
		for i, c := range *m.String {
			vec[i] = prism.String{Value: c}
		}

		return prism.Vector{
			ElementType: prism.VectorType{Type: prism.StringType},
			Body:        &vec,
		}
	case m.Alpha != nil:
		if len(*m.Alpha) == 1 {
			if f, ok := env.CurrentFunctionIR.(prism.DyadicFunction); ok {
				return prism.Alpha{
					TypeOf: f.AlphaType,
				}
			}
		}

		panic("Unreachable")

	case m.Omega != nil:
		if len(*m.Omega) == 1 {
			if f, ok := env.CurrentFunctionIR.(prism.DyadicFunction); ok {
				return prism.Omega{
					TypeOf: f.OmegaType,
				}
			} else if f, ok := env.CurrentFunctionIR.(prism.MonadicFunction); ok {
				return prism.Omega{
					TypeOf: f.OmegaType,
				}
			}
		} else {
			panic("Unreachable")

		}
	}

	panic("Other types not implemented")
}