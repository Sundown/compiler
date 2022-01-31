package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func (env Environment) AnalyseMonadic(m *palisade.Monadic) (app prism.MApplication) {
	op := env.FetchVerb(m.Verb)
	if _, ok := op.(prism.MonadicFunction); !ok {
		panic("Verb is not a monadic function")
	}

	fn := op.(prism.MonadicFunction)
	expr := env.AnalyseExpression(m.Expression)

	if !prism.PrimativeTypeEq(expr.Type(), fn.OmegaType) {
		if derived := prism.DeriveSemiDeterminedType(fn.OmegaType, expr.Type()); derived != nil {
			integrated_omega := prism.IntegrateSemiDeterminedType(derived, fn.OmegaType)

			fn.OmegaType = integrated_omega

			if prism.PredicateSemiDeterminedType(fn.Returns) {
				integrated_return := prism.IntegrateSemiDeterminedType(derived, fn.Returns)

				fn.Returns = integrated_return
			}
		} else {
			panic("Omega type mismatch between" + fn.OmegaType.String() + " and " + expr.Type().String())
		}
	}

	if fn.Name.Package == "_" && fn.Name.Name == "Return" {
		if !prism.PrimativeTypeEq(env.CurrentFunctionIR.Type(), fn.Returns) {
			if !prism.PredicateSemiDeterminedType(env.CurrentFunctionIR.Type()) {
				panic("Return recieves " + fn.Returns.String() + " which does not match determined-function's type " + env.CurrentFunctionIR.Type().String())
			} else {
				panic("Not implemented, pain")
			}
		}
	}

	return prism.MApplication{
		Operator: fn,
		Operand:  expr,
	}
}
func (env Environment) AnalyseDyadicOperator(d *palisade.Monadic) prism.DyadicOperator {
	dop := prism.DyadicOperator{}
	var lexpr prism.Expression
	if d.Verb != nil {
		lexpr = env.FetchVerb(d.Verb)
	} else {
		panic("Dyadic expression has no left operand")
	}

	rexpr := env.AnalyseExpression(d.Expression.Monadic.Expression)

	switch *d.Expression.Monadic.Verb.Ident {
	case "Map":
		if _, ok := lexpr.(prism.Function); !ok {
			panic("Left operand is not a function")
		}
		if _, ok := rexpr.(prism.Vector); !ok {
			panic("Right operand is not a vector")
		}

		dop = prism.DyadicOperator{
			Operator: prism.KindMapOperator,
			Left:     lexpr.(prism.Function),
			Right:    rexpr.(prism.Vector),
		}
	case "/":
		if _, ok := lexpr.(prism.Function); !ok {
			panic("Left operand is not a function")
		}
		if _, ok := rexpr.(prism.Vector); !ok {
			panic("Right operand is not a vector")
		}

		dop = prism.DyadicOperator{
			Operator: prism.KindFoldlOperator,
			Left:     lexpr.(prism.Function),
			Right:    rexpr.(prism.Vector),
		}
	}

	return dop
}

func (env Environment) AnalyseDyadic(d *palisade.Dyadic) prism.DApplication {
	op := env.FetchVerb(d.Verb)
	if _, ok := op.(prism.DyadicFunction); !ok {
		panic("Verb is not a dyadic function")
	}
	var left prism.Expression
	if d.Monadic != nil {
		left = env.AnalyseMonadic(d.Monadic)
	} else if d.Morphemes != nil {
		left = env.AnalyseMorphemes(d.Morphemes)
	} else {
		panic("Dyadic expression has no left operand")
	}

	right := env.AnalyseExpression(d.Expression)

	fn := op.(prism.DyadicFunction)

	tmp := right.Type()
	resolved_right, err := prism.Delegate(&fn.OmegaType, &tmp)
	if err != nil {
		prism.Panic(*err)
	}
	tmp = left.Type()
	resolved_left, err := prism.Delegate(&fn.AlphaType, &tmp)
	if err != nil {
		prism.Panic(*err)
	}

	if _, err := prism.Delegate(resolved_left, resolved_right); err != nil {
		prism.Panic(*err)
	}

	if prism.PredicateSemiDeterminedType(fn.Returns) {
		fn.Returns = prism.IntegrateSemiDeterminedType(*resolved_left, fn.Returns)
	}

	if fn.Name.Package == "_" && fn.Name.Name == "Return" {
		if !prism.PrimativeTypeEq(env.CurrentFunctionIR.Type(), fn.Returns) {
			if !prism.PredicateSemiDeterminedType(env.CurrentFunctionIR.Type()) {
				panic("Return recieves type which does not match determined-function's type")
			} else {
				panic("Not implemented, pain")
			}
		}
	}

	return prism.DApplication{
		Operator: fn,
		Left:     left,
		Right:    right,
	}
}
