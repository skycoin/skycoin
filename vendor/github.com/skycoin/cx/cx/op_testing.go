package base

import (
	"fmt"
	"os"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
)

var assertSuccess = true

// AssertFailed ...
func AssertFailed() bool {
	return assertSuccess == false
}

func assert(expr *CXExpression, fp int) (same bool) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	var byts1, byts2 []byte

	if inp1.Type == TYPE_STR {
		byts1 = []byte(ReadStr(fp, inp1))
		byts2 = []byte(ReadStr(fp, inp2))
	} else {
		byts1 = ReadMemory(GetFinalOffset(fp, inp1), inp1)
		byts2 = ReadMemory(GetFinalOffset(fp, inp2), inp2)
	}

	same = true

	if len(byts1) != len(byts2) {
		same = false
		fmt.Println("byts1", byts1)
		fmt.Println("byts2", byts2)
	}

	if same {
		for i, byt := range byts1 {
			if byt != byts2[i] {
				same = false
				fmt.Println("byts1", byts1)
				fmt.Println("byts2", byts2)
				break
			}
		}
	}

	message := ReadStr(fp, inp3)

	if !same {
		if message != "" {
			fmt.Printf("%s: %d: result was not equal to the expected value; %s\n", expr.FileName, expr.FileLine, message)
		} else {
			fmt.Printf("%s: %d: result was not equal to the expected value\n", expr.FileName, expr.FileLine)
		}
	}

	assertSuccess = assertSuccess && same
	return same
}

func opAssertValue(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	same := assert(expr, fp)
	WriteMemory(GetFinalOffset(fp, out1), FromBool(same))
}

func opTest(expr *CXExpression, fp int) {
	assert(expr, fp)
}

func opPanic(expr *CXExpression, fp int) {
	if assert(expr, fp) == false {
		os.Exit(CX_ASSERT)
	}
}

func opStrError(expr *CXExpression, fp int) {
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(ErrorString(int(ReadI32(fp, expr.Inputs[0])))))
}
