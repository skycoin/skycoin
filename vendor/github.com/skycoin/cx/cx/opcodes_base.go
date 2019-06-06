// +build base extra full

package base

// op codes
const (
	// time
	OP_TIME_SLEEP = iota + END_OF_BARE_OPS
	OP_TIME_UNIX
	OP_TIME_UNIX_MILLI
	OP_TIME_UNIX_NANO

	// serialize
	OP_SERIAL_PROGRAM

	// os
	OP_OS_GET_WORKING_DIRECTORY
	OP_OS_OPEN
	OP_OS_CLOSE
	OP_OS_RUN
	OP_OS_EXIT

	// http
	OP_HTTP_GET

	// object explorer
	OP_OBJ_QUERY

	END_OF_BASE_OPS
)

var execNativeBase func(*CXProgram)

func init() {
	// time
	AddOpCode(OP_TIME_SLEEP, "time.Sleep",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_TIME_UNIX_MILLI, "time.UnixMilli",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_TIME_UNIX_NANO, "time.UnixNano",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_I64, false)})

	// http
	AddOpCode(OP_HTTP_GET, "http.Get",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})

	// os
	AddOpCode(OP_OS_GET_WORKING_DIRECTORY, "os.GetWorkingDirectory",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_OS_OPEN, "os.Open",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_OS_CLOSE, "os.Close",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_OS_RUN, "os.Run",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_STR, false)})
	AddOpCode(OP_OS_EXIT, "os.Exit",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})

	// exec
	execNativeBase = func(prgrm *CXProgram) {
		defer RuntimeError()
		call := &prgrm.CallStack[prgrm.CallCounter]
		expr := call.Operator.Expressions[call.Line]
		opCode := expr.Operator.OpCode
		fp := call.FramePointer

		if opCode < END_OF_BARE_OPS {
			execNativeBare(prgrm)
		} else {
			switch opCode {
			// time
			case OP_TIME_SLEEP:
				op_time_Sleep(expr, fp)
			case OP_TIME_UNIX:
			case OP_TIME_UNIX_MILLI:
				op_time_UnixMilli(expr, fp)
			case OP_TIME_UNIX_NANO:
				op_time_UnixNano(expr, fp)

			// http
			case OP_HTTP_GET:
				op_http_get(expr, fp)

			// os
			case OP_OS_GET_WORKING_DIRECTORY:
				op_os_GetWorkingDirectory(expr, fp)
			case OP_OS_OPEN:
				op_os_Open(expr, fp)
			case OP_OS_CLOSE:
				op_os_Close(expr, fp)
			case OP_OS_RUN:
				op_os_Run(expr, fp)
			case OP_OS_EXIT:
				op_os_Exit(expr, fp)

			default:
				// DumpOpCodes(opCode)
				panic("invalid base opcode")
			}
		}
	}

	execNative = execNativeBase
}
