package base

// CorePackages ...
var CorePackages = []string{
	// temporary solution until we can implement these packages in pure CX I guess
	"gl", "glfw", "time", "http", "os", "explorer", "aff", "gltext", "cx",
}

// op codes
// nolint golint
const (
	OP_IDENTITY = iota + 1
	OP_JMP
	OP_DEBUG

	OP_SERIALIZE
	OP_DESERIALIZE

	START_UND_OPS

	OP_UND_EQUAL
	OP_UND_UNEQUAL
	OP_UND_BITAND
	OP_UND_BITXOR
	OP_UND_BITOR
	OP_UND_BITCLEAR
	OP_UND_MUL
	OP_UND_DIV
	OP_UND_MOD
	OP_UND_ADD
	OP_UND_SUB
	OP_UND_BITSHL
	OP_UND_BITSHR
	OP_UND_LT
	OP_UND_GT
	OP_UND_LTEQ
	OP_UND_GTEQ
	OP_UND_LEN
	OP_UND_PRINTF
	OP_UND_SPRINTF
	OP_UND_READ

	END_UND_OPS

	// parse ops

	START_PARSE_OPS

	OP_BYTE_BYTE
	OP_BYTE_STR
	OP_BYTE_I32
	OP_BYTE_I64
	OP_BYTE_F32
	OP_BYTE_F64

	OP_I32_BYTE
	OP_I32_STR
	OP_I32_I32
	OP_I32_I64
	OP_I32_F32
	OP_I32_F64

	OP_I64_BYTE
	OP_I64_STR
	OP_I64_I32
	OP_I64_I64
	OP_I64_F32
	OP_I64_F64

	OP_F32_BYTE
	OP_F32_STR
	OP_F32_I32
	OP_F32_I64
	OP_F32_F32
	OP_F32_F64

	OP_F64_BYTE
	OP_F64_STR
	OP_F64_I32
	OP_F64_I64
	OP_F64_F32
	OP_F64_F64

	OP_STR_BYTE
	OP_STR_STR
	OP_STR_I32
	OP_STR_I64
	OP_STR_F32
	OP_STR_F64

	END_PARSE_OPS

	OP_BOOL_PRINT

	OP_BOOL_EQUAL
	OP_BOOL_UNEQUAL
	OP_BOOL_NOT
	OP_BOOL_OR
	OP_BOOL_AND

	OP_BYTE_PRINT

	OP_I32_PRINT
	OP_I32_ADD
	OP_I32_SUB
	OP_I32_MUL
	OP_I32_DIV
	OP_I32_ABS
	OP_I32_POW
	OP_I32_GT
	OP_I32_GTEQ
	OP_I32_LT
	OP_I32_LTEQ
	OP_I32_EQ
	OP_I32_UNEQ
	OP_I32_MOD
	OP_I32_RAND
	OP_I32_BITAND
	OP_I32_BITOR
	OP_I32_BITXOR
	OP_I32_BITCLEAR
	OP_I32_BITSHL
	OP_I32_BITSHR
	OP_I32_SQRT
	OP_I32_LOG
	OP_I32_LOG2
	OP_I32_LOG10
	OP_I32_MAX
	OP_I32_MIN

	OP_I64_PRINT
	OP_I64_ADD
	OP_I64_SUB
	OP_I64_MUL
	OP_I64_DIV
	OP_I64_ABS
	OP_I64_POW
	OP_I64_GT
	OP_I64_GTEQ
	OP_I64_LT
	OP_I64_LTEQ
	OP_I64_EQ
	OP_I64_UNEQ
	OP_I64_MOD
	OP_I64_RAND
	OP_I64_BITAND
	OP_I64_BITOR
	OP_I64_BITXOR
	OP_I64_BITCLEAR
	OP_I64_BITSHL
	OP_I64_BITSHR
	OP_I64_SQRT
	OP_I64_LOG
	OP_I64_LOG10
	OP_I64_LOG2
	OP_I64_MAX
	OP_I64_MIN

	OP_F32_IS_NAN

	OP_F32_PRINT
	OP_F32_ADD
	OP_F32_SUB
	OP_F32_MUL
	OP_F32_DIV
	OP_F32_ABS
	OP_F32_POW
	OP_F32_GT
	OP_F32_GTEQ
	OP_F32_LT
	OP_F32_LTEQ
	OP_F32_EQ
	OP_F32_UNEQ
	OP_F32_COS
	OP_F32_SIN
	OP_F32_SQRT
	OP_F32_LOG
	OP_F32_LOG2
	OP_F32_LOG10
	OP_F32_MAX
	OP_F32_MIN

	OP_F64_PRINT
	OP_F64_ADD
	OP_F64_SUB
	OP_F64_MUL
	OP_F64_DIV
	OP_F64_ABS
	OP_F64_POW
	OP_F64_GT
	OP_F64_GTEQ
	OP_F64_LT
	OP_F64_LTEQ
	OP_F64_EQ
	OP_F64_UNEQ
	OP_F64_COS
	OP_F64_SIN

	OP_F64_SQRT
	OP_F64_LOG
	OP_F64_LOG2
	OP_F64_LOG10
	OP_F64_MAX
	OP_F64_MIN

	OP_STR_PRINT
	OP_STR_CONCAT
	OP_STR_SUBSTR
	OP_STR_INDEX
	OP_STR_TRIM_SPACE
	OP_STR_EQ

	OP_APPEND
	OP_RESIZE
	OP_INSERT
	OP_REMOVE
	OP_COPY

	OP_MAKE
	OP_READ
	OP_WRITE
	OP_LEN
	OP_CONCAT
	OP_CAST
	OP_EQ
	OP_UNEQ
	OP_RAND
	OP_AND
	OP_OR
	OP_NOT
	OP_SLEEP
	OP_HALT
	OP_GOTO
	OP_REMCX
	OP_ADDCX
	OP_QUERY
	OP_EXECUTE
	OP_INDEX
	OP_NAME
	OP_EVOLVE

	OP_ASSERT
	OP_TEST
	OP_PANIC
	OP_STRERROR

	// affordances
	OP_AFF_PRINT
	OP_AFF_QUERY
	OP_AFF_ON
	OP_AFF_OF
	OP_AFF_INFORM
	OP_AFF_REQUEST

	OP_UND_NEG
	OP_I32_NEG
	OP_I64_NEG
	OP_F32_NEG
	OP_F64_NEG
	END_OF_BARE_OPS
)

// For the parser. These shouldn't be used in the runtime for performance reasons
var (
	OpNames        = map[int]string{}
	OpCodes        = map[string]int{}
	Natives        = map[int]*CXFunction{}
	execNativeBare func(*CXProgram)
	execNative     func(*CXProgram)
)

// AddOpCode ...
func AddOpCode(code int, name string, inputs []*CXArgument, outputs []*CXArgument) {
	OpNames[code] = name
	OpCodes[name] = code
	Natives[code] = MakeNativeFunction(code, inputs, outputs)
}

/*
// debug helper
func DumpOpCodes(opCode int) () {
	var keys []int
	for k := range OpNames {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Printf("%5d : %s\n", k, OpNames[k])
	}

	fmt.Printf("opCode : %d\n", opCode)
}*/

// Helper function for creating parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
func newOpPar(typCode int, isSlice bool) *CXArgument {
	arg := MakeArgument("", "", -1).AddType(TypeNames[typCode])
	if isSlice {
		arg.IsSlice = true
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_SLICE)
	}
	return arg
}

func init() {
	AddOpCode(OP_IDENTITY, "identity",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_JMP, "jmp",
		[]*CXArgument{newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{}) // newOpPar(TYPE_UNDEFINED, false) to allow 0 inputs (goto)
	AddOpCode(OP_DEBUG, "debug",
		[]*CXArgument{},
		[]*CXArgument{})

	AddOpCode(OP_SERIALIZE, "serialize",
		[]*CXArgument{newOpPar(TYPE_AFF, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_DESERIALIZE, "deserialize",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{})

	AddOpCode(OP_UND_EQUAL, "eq",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_UNEQUAL, "uneq",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_BITAND, "bitand",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_BITXOR, "bitxor",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_BITOR, "bitor",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_BITCLEAR, "bitclear",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_MUL, "mul",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_DIV, "div",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_MOD, "mod",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_ADD, "add",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_SUB, "sub",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_NEG, "neg",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_BITSHL, "bitshl",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_BITSHR, "bitshr",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_LT, "lt",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_GT, "gt",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_LTEQ, "lteq",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_GTEQ, "gteq",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_LEN, "len",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_UND_PRINTF, "printf",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{})
	AddOpCode(OP_UND_SPRINTF, "sprintf",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_UND_READ, "read",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_STR, false)})

	AddOpCode(OP_BYTE_BYTE, "byte.byte",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_BYTE_STR, "byte.str",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_BYTE_I32, "byte.i32",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_BYTE_I64, "byte.i64",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_BYTE_F32, "byte.f32",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_BYTE_F64, "byte.f64",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})

	AddOpCode(OP_BYTE_PRINT, "byte.print",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{})

	AddOpCode(OP_BOOL_PRINT, "bool.print",
		[]*CXArgument{newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{})
	AddOpCode(OP_BOOL_EQUAL, "bool.eq",
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_BOOL_UNEQUAL, "bool.uneq",
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_BOOL_NOT, "bool.not",
		[]*CXArgument{newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_BOOL_OR, "bool.or",
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_BOOL_AND, "bool.and",
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})

	AddOpCode(OP_I32_BYTE, "i32.byte",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_I32_STR, "i32.str",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_I32_I32, "i32.i32",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_I64, "i32.i64",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I32_F32, "i32.f32",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_I32_F64, "i32.f64",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})

	AddOpCode(OP_I32_PRINT, "i32.print",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_I32_ADD, "i32.add",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_SUB, "i32.sub",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_NEG, "i32.neg",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_MUL, "i32.mul",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_DIV, "i32.div",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_ABS, "i32.abs",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_POW, "i32.pow",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_GT, "i32.gt",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_GTEQ, "i32.gteq",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_LT, "i32.lt",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_LTEQ, "i32.lteq",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_EQ, "i32.eq",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_UNEQ, "i32.uneq",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_MOD, "i32.mod",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_RAND, "i32.rand",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITAND, "i32.bitand",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITOR, "i32.bitor",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITXOR, "i32.bitxor",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITCLEAR, "i32.bitclear",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITSHL, "i32.bitshl",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITSHR, "i32.bitshr",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_SQRT, "i32.sqrt",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_LOG, "i32.log",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_LOG2, "i32.log2",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_LOG10, "i32.log10",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_MAX, "i32.max",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_MIN, "i32.min",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})

	AddOpCode(OP_I64_BYTE, "i64.byte",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_I64_STR, "i64.str",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_I64_I32, "i64.i32",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I64_I64, "i64.i64",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_F32, "i64.f32",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_I64_F64, "i64.f64",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})

	AddOpCode(OP_I64_PRINT, "i64.print",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{})
	AddOpCode(OP_I64_ADD, "i64.add",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_SUB, "i64.sub",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_NEG, "i64.neg",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_MUL, "i64.mul",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_DIV, "i64.div",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_ABS, "i64.abs",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_POW, "i64.pow",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_GT, "i64.gt",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_GTEQ, "i64.gteq",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_LT, "i64.lt",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_LTEQ, "i64.lteq",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_EQ, "i64.eq",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_UNEQ, "i64.uneq",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_MOD, "i64.mod",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_RAND, "i64.rand",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITAND, "i64.bitand",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITOR, "i64.bitor",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITXOR, "i64.bitxor",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITCLEAR, "i64.bitclear",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITSHL, "i64.bitshl",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITSHR, "i64.bitshr",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_SQRT, "i64.sqrt",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_LOG, "i64.log",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_LOG2, "i64.log2",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_LOG10, "i64.log10",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_MAX, "i64.max",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_MIN, "i64.min",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})

	AddOpCode(OP_F32_IS_NAN, "f32.isnan",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_BYTE, "f32.byte",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_F32_STR, "f32.str",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_F32_I32, "f32.i32",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_F32_I64, "f32.i64",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_F32_F32, "f32.f32",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_F64, "f32.f64",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F32_PRINT, "f32.print",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_F32_ADD, "f32.add",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_SUB, "f32.sub",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_NEG, "f32.neg",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_MUL, "f32.mul",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_DIV, "f32.div",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_ABS, "f32.abs",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_POW, "f32.pow",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_GT, "f32.gt",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_GTEQ, "f32.gteq",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_LT, "f32.lt",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_LTEQ, "f32.lteq",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_EQ, "f32.eq",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_UNEQ, "f32.uneq",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_COS, "f32.cos",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_SIN, "f32.sin",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_SQRT, "f32.sqrt",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_LOG, "f32.log",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_LOG2, "f32.log2",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_LOG10, "f32.log10",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_MAX, "f32.max",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_MIN, "f32.min",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})

	AddOpCode(OP_F64_BYTE, "f64.byte",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_F64_STR, "f64.str",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_F64_I32, "f64.i32",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_F64_I64, "f64.i64",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_F64_F32, "f64.f32",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F64_F64, "f64.f64",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})

	AddOpCode(OP_F64_PRINT, "f64.print",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{})
	AddOpCode(OP_F64_ADD, "f64.add",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_SUB, "f64.sub",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_NEG, "f64.neg",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_MUL, "f64.mul",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_DIV, "f64.div",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_ABS, "f64.abs",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_POW, "f64.pow",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_GT, "f64.gt",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_GTEQ, "f64.gteq",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_LT, "f64.lt",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_LTEQ, "f64.lteq",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_EQ, "f64.eq",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_UNEQ, "f64.uneq",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_COS, "f64.cos",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_SIN, "f64.sin",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_SQRT, "f64.sqrt",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_LOG, "f64.log",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_LOG2, "f64.log2",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_LOG10, "f64.log10",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_MAX, "f64.max",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_MIN, "f64.min",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})

	AddOpCode(OP_STR_PRINT, "str.print",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_STR_CONCAT, "str.concat",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_STR_SUBSTR, "str.substr",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_STR_INDEX, "str.index",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_STR_TRIM_SPACE, "str.trimspace",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_STR_EQ, "str.eq",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})

	AddOpCode(OP_STR_BYTE, "str.byte",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_STR_STR, "str.str",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_STR_I32, "str.i32",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_STR_I64, "str.i64",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_STR_F32, "str.f32",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_STR_F64, "str.f64",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})

	AddOpCode(OP_APPEND, "append",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true), newOpPar(TYPE_UNDEFINED, true)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true)})
	AddOpCode(OP_RESIZE, "resize",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true)})
	AddOpCode(OP_INSERT, "insert",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true), newOpPar(TYPE_UNDEFINED, true)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true)})
	AddOpCode(OP_REMOVE, "remove",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true)})
	AddOpCode(OP_COPY, "copy",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true), newOpPar(TYPE_UNDEFINED, true)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})

	AddOpCode(OP_ASSERT, "assert",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_TEST, "test",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_PANIC, "panic",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_STRERROR, "strerror",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	// affordances
	AddOpCode(OP_AFF_PRINT, "aff.print",
		[]*CXArgument{newOpPar(TYPE_AFF, false)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_QUERY, "aff.query",
		[]*CXArgument{newOpPar(TYPE_AFF, false)},
		[]*CXArgument{newOpPar(TYPE_AFF, false)})
	AddOpCode(OP_AFF_ON, "aff.on",
		[]*CXArgument{newOpPar(TYPE_AFF, false), newOpPar(TYPE_AFF, false)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_OF, "aff.of",
		[]*CXArgument{newOpPar(TYPE_AFF, false), newOpPar(TYPE_AFF, false)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_INFORM, "aff.inform",
		[]*CXArgument{newOpPar(TYPE_AFF, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_AFF, false)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_REQUEST, "aff.request",
		[]*CXArgument{newOpPar(TYPE_AFF, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_AFF, false)},
		[]*CXArgument{})

	// exec
	execNativeBare = func(prgrm *CXProgram) {
		defer RuntimeError()
		call := &prgrm.CallStack[prgrm.CallCounter]
		expr := call.Operator.Expressions[call.Line]
		opCode := expr.Operator.OpCode
		fp := call.FramePointer

		switch opCode {
		case OP_IDENTITY:
			opIdentity(expr, fp)
		case OP_JMP:
			opJmp(expr, fp, call)
		case OP_DEBUG:
			prgrm.PrintStack()

		case OP_SERIALIZE:
			opSerialize(expr, fp)
		case OP_DESERIALIZE:
			opDeserialize(expr, fp)

		case OP_UND_EQUAL:
			opEqual(expr, fp)
		case OP_UND_UNEQUAL:
			opUnequal(expr, fp)
		case OP_UND_BITAND:
			opBitand(expr, fp)
		case OP_UND_BITXOR:
			opBitxor(expr, fp)
		case OP_UND_BITOR:
			opBitor(expr, fp)
		case OP_UND_BITCLEAR:
			opBitclear(expr, fp)
		case OP_UND_MUL:
			opMul(expr, fp)
		case OP_UND_DIV:
			opDiv(expr, fp)
		case OP_UND_MOD:
			opMod(expr, fp)
		case OP_UND_ADD:
			opAdd(expr, fp)
		case OP_UND_SUB, OP_UND_NEG:
			opSub(expr, fp)
		case OP_UND_BITSHL:
			opBitshl(expr, fp)
		case OP_UND_BITSHR:
			opBitshr(expr, fp)
		case OP_UND_LT:
			opLt(expr, fp)
		case OP_UND_GT:
			opGt(expr, fp)
		case OP_UND_LTEQ:
			opLteq(expr, fp)
		case OP_UND_GTEQ:
			opGteq(expr, fp)
		case OP_UND_LEN:
			opLen(expr, fp)
		case OP_UND_PRINTF:
			opPrintf(expr, fp)
		case OP_UND_SPRINTF:
			opSprintf(expr, fp)
		case OP_UND_READ:
			opRead(expr, fp)

		case OP_BYTE_BYTE:
			opByteByte(expr, fp)
		case OP_BYTE_STR:
			opByteByte(expr, fp)
		case OP_BYTE_I32:
			opByteByte(expr, fp)
		case OP_BYTE_I64:
			opByteByte(expr, fp)
		case OP_BYTE_F32:
			opByteByte(expr, fp)
		case OP_BYTE_F64:
			opByteByte(expr, fp)

		case OP_BYTE_PRINT:
			opBytePrint(expr, fp)

		case OP_BOOL_PRINT:
			opBoolPrint(expr, fp)
		case OP_BOOL_EQUAL:
			opBoolEqual(expr, fp)
		case OP_BOOL_UNEQUAL:
			opBoolUnequal(expr, fp)
		case OP_BOOL_NOT:
			opBoolNot(expr, fp)
		case OP_BOOL_OR:
			opBoolOr(expr, fp)
		case OP_BOOL_AND:
			opBoolAnd(expr, fp)

		case OP_I32_BYTE:
			opI32I32(expr, fp)
		case OP_I32_STR:
			opI32I32(expr, fp)
		case OP_I32_I32:
			opI32I32(expr, fp)
		case OP_I32_I64:
			opI32I32(expr, fp)
		case OP_I32_F32:
			opI32I32(expr, fp)
		case OP_I32_F64:
			opI32I32(expr, fp)

		case OP_I32_PRINT:
			opI32Print(expr, fp)
		case OP_I32_ADD:
			opI32Add(expr, fp)
		case OP_I32_SUB, OP_I32_NEG:
			opI32Sub(expr, fp)
		case OP_I32_MUL:
			opI32Mul(expr, fp)
		case OP_I32_DIV:
			opI32Div(expr, fp)
		case OP_I32_ABS:
			opI32Abs(expr, fp)
		case OP_I32_POW:
			opI32Pow(expr, fp)
		case OP_I32_GT:
			opI32Gt(expr, fp)
		case OP_I32_GTEQ:
			opI32Gteq(expr, fp)
		case OP_I32_LT:
			opI32Lt(expr, fp)
		case OP_I32_LTEQ:
			opI32Lteq(expr, fp)
		case OP_I32_EQ:
			opI32Eq(expr, fp)
		case OP_I32_UNEQ:
			opI32Uneq(expr, fp)
		case OP_I32_MOD:
			opI32Mod(expr, fp)
		case OP_I32_RAND:
			opI32Rand(expr, fp)
		case OP_I32_BITAND:
			opI32Bitand(expr, fp)
		case OP_I32_BITOR:
			opI32Bitor(expr, fp)
		case OP_I32_BITXOR:
			opI32Bitxor(expr, fp)
		case OP_I32_BITCLEAR:
			opI32Bitclear(expr, fp)
		case OP_I32_BITSHL:
			opI32Bitshl(expr, fp)
		case OP_I32_BITSHR:
			opI32Bitshr(expr, fp)
		case OP_I32_SQRT:
			opI32Sqrt(expr, fp)
		case OP_I32_LOG:
			opI32Log(expr, fp)
		case OP_I32_LOG2:
			opI32Log2(expr, fp)
		case OP_I32_LOG10:
			opI32Log10(expr, fp)

		case OP_I32_MAX:
			opI32Max(expr, fp)
		case OP_I32_MIN:
			opI32Min(expr, fp)

		case OP_I64_BYTE:
			opI64I64(expr, fp)
		case OP_I64_STR:
			opI64I64(expr, fp)
		case OP_I64_I32:
			opI64I64(expr, fp)
		case OP_I64_I64:
			opI64I64(expr, fp)
		case OP_I64_F32:
			opI64I64(expr, fp)
		case OP_I64_F64:
			opI64I64(expr, fp)

		case OP_I64_PRINT:
			opI64Print(expr, fp)
		case OP_I64_ADD:
			opI64Add(expr, fp)
		case OP_I64_SUB, OP_I64_NEG:
			opI64Sub(expr, fp)
		case OP_I64_MUL:
			opI64Mul(expr, fp)
		case OP_I64_DIV:
			opI64Div(expr, fp)
		case OP_I64_ABS:
			opI64Abs(expr, fp)
		case OP_I64_POW:
			opI64Pow(expr, fp)
		case OP_I64_GT:
			opI64Gt(expr, fp)
		case OP_I64_GTEQ:
			opI64Gteq(expr, fp)
		case OP_I64_LT:
			opI64Lt(expr, fp)
		case OP_I64_LTEQ:
			opI64Lteq(expr, fp)
		case OP_I64_EQ:
			opI64Eq(expr, fp)
		case OP_I64_UNEQ:
			opI64Uneq(expr, fp)
		case OP_I64_MOD:
			opI64Mod(expr, fp)
		case OP_I64_RAND:
			opI64Rand(expr, fp)
		case OP_I64_BITAND:
			opI64Bitand(expr, fp)
		case OP_I64_BITOR:
			opI64Bitor(expr, fp)
		case OP_I64_BITXOR:
			opI64Bitxor(expr, fp)
		case OP_I64_BITCLEAR:
			opI64Bitclear(expr, fp)
		case OP_I64_BITSHL:
			opI64Bitshl(expr, fp)
		case OP_I64_BITSHR:
			opI64Bitshr(expr, fp)
		case OP_I64_SQRT:
			opI64Sqrt(expr, fp)
		case OP_I64_LOG:
			opI64Log(expr, fp)
		case OP_I64_LOG2:
			opI64Log2(expr, fp)
		case OP_I64_LOG10:
			opI64Log10(expr, fp)
		case OP_I64_MAX:
			opI64Max(expr, fp)
		case OP_I64_MIN:
			opI64Min(expr, fp)

		case OP_F32_IS_NAN:
			opF32Isnan(expr, fp)
		case OP_F32_BYTE:
			opF32F32(expr, fp)
		case OP_F32_STR:
			opF32F32(expr, fp)
		case OP_F32_I32:
			opF32F32(expr, fp)
		case OP_F32_I64:
			opF32F32(expr, fp)
		case OP_F32_F32:
			opF32F32(expr, fp)
		case OP_F32_F64:
			opF32F32(expr, fp)
		case OP_F32_PRINT:
			opF32Print(expr, fp)
		case OP_F32_ADD:
			opF32Add(expr, fp)
		case OP_F32_SUB, OP_F32_NEG:
			opF32Sub(expr, fp)
		case OP_F32_MUL:
			opF32Mul(expr, fp)
		case OP_F32_DIV:
			opF32Div(expr, fp)
		case OP_F32_ABS:
			opF32Abs(expr, fp)
		case OP_F32_POW:
			opF32Pow(expr, fp)
		case OP_F32_GT:
			opF32Gt(expr, fp)
		case OP_F32_GTEQ:
			opF32Gteq(expr, fp)
		case OP_F32_LT:
			opF32Lt(expr, fp)
		case OP_F32_LTEQ:
			opF32Lteq(expr, fp)
		case OP_F32_EQ:
			opF32Eq(expr, fp)
		case OP_F32_UNEQ:
			opF32Uneq(expr, fp)
		case OP_F32_COS:
			opF32Cos(expr, fp)
		case OP_F32_SIN:
			opF32Sin(expr, fp)
		case OP_F32_SQRT:
			opF32Sqrt(expr, fp)
		case OP_F32_LOG:
			opF32Log(expr, fp)
		case OP_F32_LOG2:
			opF32Log2(expr, fp)
		case OP_F32_LOG10:
			opF32Log10(expr, fp)
		case OP_F32_MAX:
			opF32Max(expr, fp)
		case OP_F32_MIN:
			opF32Min(expr, fp)

		case OP_F64_BYTE:
			opF64F64(expr, fp)
		case OP_F64_STR:
			opF64F64(expr, fp)
		case OP_F64_I32:
			opF64F64(expr, fp)
		case OP_F64_I64:
			opF64F64(expr, fp)
		case OP_F64_F32:
			opF64F64(expr, fp)
		case OP_F64_F64:
			opF64F64(expr, fp)

		case OP_F64_PRINT:
			opF64Print(expr, fp)
		case OP_F64_ADD:
			opF64Add(expr, fp)
		case OP_F64_SUB, OP_F64_NEG:
			opF64Sub(expr, fp)
		case OP_F64_MUL:
			opF64Mul(expr, fp)
		case OP_F64_DIV:
			opF64Div(expr, fp)
		case OP_F64_ABS:
			opF64Abs(expr, fp)
		case OP_F64_POW:
			opF64Pow(expr, fp)
		case OP_F64_GT:
			opF64Gt(expr, fp)
		case OP_F64_GTEQ:
			opF64Gteq(expr, fp)
		case OP_F64_LT:
			opF64Lt(expr, fp)
		case OP_F64_LTEQ:
			opF64Lteq(expr, fp)
		case OP_F64_EQ:
			opF64Eq(expr, fp)
		case OP_F64_UNEQ:
			opF64Uneq(expr, fp)
		case OP_F64_COS:
			opF64Cos(expr, fp)
		case OP_F64_SIN:
			opF64Sin(expr, fp)
		case OP_F64_SQRT:
			opF64Sqrt(expr, fp)
		case OP_F64_LOG:
			opF64Log(expr, fp)
		case OP_F64_LOG2:
			opF64Log2(expr, fp)
		case OP_F64_LOG10:
			opF64Log10(expr, fp)
		case OP_F64_MAX:
			opF64Max(expr, fp)
		case OP_F64_MIN:
			opF64Min(expr, fp)

		case OP_STR_PRINT:
			opStrPrint(expr, fp)
		case OP_STR_EQ:
			opStrEq(expr, fp)
		case OP_STR_CONCAT:
			opStrConcat(expr, fp)
		case OP_STR_SUBSTR:
			opStrSubstr(expr, fp)
		case OP_STR_INDEX:
			opStrIndex(expr, fp)
		case OP_STR_TRIM_SPACE:
			opStrTrimSpace(expr, fp)

		case OP_STR_BYTE:
			opStrStr(expr, fp)
		case OP_STR_STR:
			opStrStr(expr, fp)
		case OP_STR_I32:
			opStrStr(expr, fp)
		case OP_STR_I64:
			opStrStr(expr, fp)
		case OP_STR_F32:
			opStrStr(expr, fp)
		case OP_STR_F64:
			opStrStr(expr, fp)

		case OP_APPEND:
			opAppend(expr, fp)
		case OP_RESIZE:
			opResize(expr, fp)
		case OP_INSERT:
			opInsert(expr, fp)
		case OP_REMOVE:
			opRemove(expr, fp)
		case OP_COPY:
			opCopy(expr, fp)

		case OP_MAKE:
		case OP_READ:
		case OP_WRITE:
		case OP_LEN:
		case OP_CONCAT:
		case OP_CAST:
		case OP_EQ:
		case OP_UNEQ:
		case OP_AND:
		case OP_OR:
		case OP_NOT:
		case OP_SLEEP:
		case OP_HALT:
		case OP_GOTO:
		case OP_REMCX:
		case OP_ADDCX:
		case OP_QUERY:
		case OP_EXECUTE:
		case OP_INDEX:
		case OP_NAME:
		case OP_EVOLVE:
		case OP_ASSERT:
			opAssertValue(expr, fp)
		case OP_TEST:
			opTest(expr, fp)
		case OP_PANIC:
			opPanic(expr, fp)
		case OP_STRERROR:
			opStrError(expr, fp)

		// affordances
		case OP_AFF_PRINT:
			opAffPrint(expr, fp)
		case OP_AFF_QUERY:
			opAffQuery(expr, fp)
		case OP_AFF_ON:
			opAffOn(expr, fp)
		case OP_AFF_OF:
			opAffOf(expr, fp)
		case OP_AFF_INFORM:
			opAffInform(expr, fp)
		case OP_AFF_REQUEST:
			opAffRequest(expr, fp)
		default:
			// DumpOpCodes(opCode) // debug helper
			panic("invalid bare opcode")
		}
	}

	execNative = execNativeBare
}
