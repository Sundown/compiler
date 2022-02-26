package apotheosis

import (
	"github.com/sundown/solution/prism"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (env Environment) castInt(from Value) value.Value {
	switch from.Type.Kind() {
	case prism.TypeInt:
		return from.Value
	case prism.TypeReal:
		return env.Block.NewFPToSI(from.Value, types.I64)
	case prism.TypeBool:
		return env.Block.NewSExt(from.Value, types.I64)
	}

	panic("Unreachable")
}

func (env Environment) castReal(from Value) value.Value {
	switch from.Type.Kind() {
	case prism.TypeInt:
		return env.Block.NewSIToFP(from.Value, types.Double)
	case prism.TypeReal:
		return from.Value
	case prism.TypeBool:
		return env.Block.NewSIToFP(from.Value, types.Double)
	}

	panic("Unreachable")
}

func (env Environment) compileCast(cast prism.Cast) value.Value {
	val := Value{Value: env.compileExpression(&cast.Value), Type: cast.Value.Type()}
	var castfn MCallable
	var from prism.Type
	pred := false
	if _, ok := cast.Value.Type().(prism.VectorType); ok {
		from = cast.ToType.(prism.VectorType).Type
		pred = true
	} else {
		from = cast.ToType
	}

	switch from.Kind() {
	case prism.TypeInt:
		castfn = env.castInt
	case prism.TypeReal:
		castfn = env.castReal
	}

	if pred {
		return env.vectorCast(castfn, val, cast.ToType.(prism.VectorType).Type)
	}

	panic("Unreachable")
}

func (env *Environment) vectorCast(caster MCallable, vec Value, to prism.Type) value.Value {
	elmType := vec.Type.(prism.VectorType).Type.Realise()
	irToHeadType := prism.VectorType{Type: to}
	toHeadType := irToHeadType.Realise()
	toElmType := to.Realise()
	leng := env.ReadVectorLength(vec)

	var head *ir.InstAlloca
	var body *ir.InstBitCast

	cap := env.ReadVectorCapacity(vec)
	head = env.Block.NewAlloca(toHeadType)

	env.WriteLLVectorLength(Value{head, irToHeadType}, leng)
	env.WriteLLVectorCapacity(Value{head, irToHeadType}, cap)

	// Allocate a body of capacity * element width, and cast to element type
	body = env.Block.NewBitCast(
		env.Block.NewCall(env.GetCalloc(),
			I32(to.Width()), // Byte size of elements
			cap),            // How much memory to alloc
		types.NewPointer(toElmType)) // Cast alloc'd memory to typ

	// --- Loop body ---
	vecBody := env.Block.NewLoad(
		types.NewPointer(elmType),
		env.Block.NewGetElementPtr(vec.Type.Realise(), vec.Value, I32(0), vectorBodyOffset))

	counter := env.Block.NewAlloca(types.I32)
	env.Block.NewStore(I32(0), counter)

	// Get elem, add to accum, increment counter, conditional jump to body
	loopblock := env.CurrentFunction.NewBlock("")
	env.Block.NewBr(loopblock)
	env.Block = loopblock
	// Add to accum
	curCounter := loopblock.NewLoad(types.I32, counter)

	var curElm value.Value = loopblock.NewGetElementPtr(elmType, vecBody, curCounter)

	if _, ok := vec.Type.(prism.VectorType).Type.(prism.AtomicType); ok {
		curElm = loopblock.NewLoad(elmType, curElm)
	}

	loopblock.NewStore(
		caster(Value{
			curElm,
			vec.Type.(prism.VectorType).Type}),
		loopblock.NewGetElementPtr(toElmType, body, curCounter))

	incr := loopblock.NewAdd(curCounter, I32(1))

	loopblock.NewStore(incr, counter)

	exitblock := env.CurrentFunction.NewBlock("")

	loopblock.NewCondBr(loopblock.NewICmp(enum.IPredSLT, incr, leng), loopblock, exitblock)

	env.Block = exitblock

	env.WriteVectorPointer(head, body, toHeadType)

	return head
}