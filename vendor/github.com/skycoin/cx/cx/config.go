package base

import (
	"os"
)

var InREPL bool = false
var FoundCompileErrors bool

const DBG_GOLANG_STACK_TRACE = true

// global reference to our program
var PROGRAM *CXProgram

var CXPATH string = os.Getenv("CXPATH") + "/"
var BINPATH string = CXPATH + "bin/"
var PKGPATH string = CXPATH + "pkg/"
var SRCPATH string = CXPATH + "src/"
var COREPATH string

const STACK_OVERFLOW_ERROR = "stack overflow"
const HEAP_EXHAUSTED_ERROR = "heap exhausted"
const MAIN_FUNC = "main"
const SYS_INIT_FUNC = "*init"
const MAIN_PKG = "main"
const OS_PKG = "os"
const OS_ARGS = "Args"

const NON_ASSIGN_PREFIX = "nonAssign"
const LOCAL_PREFIX = "*lcl"
const LABEL_PREFIX = "*lbl"
const ID_FN = "identity"
const INIT_FN = "initDef"

const I32_SIZE = 4
const I64_SIZE = 8
const STR_SIZE = 4

const MARK_SIZE = 1
const OBJECT_HEADER_SIZE = 9
const OBJECT_GC_HEADER_SIZE = 5
const FORWARDING_ADDRESS_SIZE = 4
const OBJECT_SIZE = 4

const CALLSTACK_SIZE = 1000

var STACK_SIZE = 1048576     // 1 Mb
var INIT_HEAP_SIZE = 2097152 // 2 Mb
var MAX_HEAP_SIZE = 67108864 // 64 Mb
const MIN_HEAP_FREE_RATIO = 40
const MAX_HEAP_FREE_RATIO = 70

var NULL_ADDRESS = STACK_SIZE

const NULL_HEAP_ADDRESS_OFFSET = 4
const NULL_HEAP_ADDRESS = 0
const STR_HEADER_SIZE = 4
const TYPE_POINTER_SIZE = 4
const SLICE_HEADER_SIZE = 8

var MEMORY_SIZE = STACK_SIZE + INIT_HEAP_SIZE + TYPE_POINTER_SIZE

const MAX_UINT32 = ^uint32(0)
const MIN_UINT32 = 0
const MAX_INT32 = int(MAX_UINT32 >> 1)
const MIN_INT32 = -MAX_INT32 - 1

var BASIC_TYPES []string = []string{
	"bool", "str", "byte", "i32", "i64", "f32", "f64",
	"[]bool", "[]str", "[]byte", "[]i32", "[]i64", "[]f32", "[]f64",
}

const (
	CX_SUCCESS = iota
	CX_COMPILATION_ERROR
	CX_PANIC // 2
	CX_INTERNAL_ERROR
	CX_ASSERT
	CX_RUNTIME_ERROR
	CX_RUNTIME_STACK_OVERFLOW_ERROR
	CX_RUNTIME_HEAP_EXHAUSTED_ERROR
	CX_RUNTIME_INVALID_ARGUMENT
	CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE
)

var ErrorStrings map[int]string = map[int]string{
	CX_SUCCESS:                          "CX_SUCCESS",
	CX_COMPILATION_ERROR:                "CX_COMPILATION_ERROR",
	CX_PANIC:                            "CX_PANIC",
	CX_INTERNAL_ERROR:                   "CX_INTERNAL_ERROR",
	CX_ASSERT:                           "CX_ASSERT",
	CX_RUNTIME_ERROR:                    "CX_RUNTIME_ERROR",
	CX_RUNTIME_STACK_OVERFLOW_ERROR:     "CX_RUNTIME_STACK_OVERFLOW_ERROR",
	CX_RUNTIME_HEAP_EXHAUSTED_ERROR:     "CX_RUNTIME_HEAP_EXHAUSTED_ERROR",
	CX_RUNTIME_INVALID_ARGUMENT:         "CX_RUNTIME_INVALID_ARGUMENT",
	CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE: "CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE",
}

const (
	DECL_POINTER  = iota // 0
	DECL_DEREF           // 1
	DECL_ARRAY           // 2
	DECL_SLICE           // 3
	DECL_STRUCT          // 4
	DECL_INDEXING        // 5
	DECL_BASIC           // 6
)

// create a new scope or return to the previous scope
const (
	SCOPE_NEW = iota + 1 // 1
	SCOPE_REM            // 2
)

// what to write
const (
	PASSBY_VALUE = iota
	PASSBY_REFERENCE
)

const (
	DEREF_ARRAY   = iota // 0
	DEREF_FIELD          // 1
	DEREF_POINTER        // 2
	DEREF_DEREF          // 3
	DEREF_SLICE          // 4
)

const (
	TYPE_UNDEFINED = iota
	TYPE_AFF
	TYPE_BOOL
	TYPE_BYTE
	TYPE_STR
	TYPE_F32
	TYPE_F64
	TYPE_I8
	TYPE_I16
	TYPE_I32
	TYPE_I64
	TYPE_UI8
	TYPE_UI16
	TYPE_UI32
	TYPE_UI64

	TYPE_CUSTOM
	TYPE_POINTER
	TYPE_ARRAY
	TYPE_SLICE
	TYPE_IDENTIFIER
)

var TypeCounter int
var TypeCodes map[string]int = map[string]int{
	"ident": TYPE_IDENTIFIER,
	"aff":   TYPE_AFF,
	"bool":  TYPE_BOOL,
	"byte":  TYPE_BYTE,
	"str":   TYPE_STR,
	"f32":   TYPE_F32,
	"f64":   TYPE_F64,
	"i8":    TYPE_I8,
	"i16":   TYPE_I16,
	"i32":   TYPE_I32,
	"i64":   TYPE_I64,
	"ui8":   TYPE_UI8,
	"ui16":  TYPE_UI16,
	"ui32":  TYPE_UI32,
	"ui64":  TYPE_UI64,
	"und":   TYPE_UNDEFINED,
}

var TypeNames map[int]string = map[int]string{
	TYPE_IDENTIFIER: "ident",
	TYPE_AFF:        "aff",
	TYPE_BOOL:       "bool",
	TYPE_BYTE:       "byte",
	TYPE_STR:        "str",
	TYPE_F32:        "f32",
	TYPE_F64:        "f64",
	TYPE_I8:         "i8",
	TYPE_I16:        "i16",
	TYPE_I32:        "i32",
	TYPE_I64:        "i64",
	TYPE_UI8:        "ui8",
	TYPE_UI16:       "ui16",
	TYPE_UI32:       "ui32",
	TYPE_UI64:       "ui64",
	TYPE_UNDEFINED:  "und",
}

// memory locations
const (
	MEM_STACK = iota
	MEM_HEAP
	MEM_DATA
)
