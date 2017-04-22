package zen

import (
	"encoding/binary"
	"errors"
	// "fmt"
	"unsafe"

	"github.com/ngaut/log"
)

type VM struct {
	globalValue []*Value
	code        []byte
}

type returnFrame struct {
	pc  int
	env *Value
}

var (
	errBadMagic      = errors.New("error bad magic")
	errTruncatedFile = errors.New("file seems to be truncated")
	errInvalidEXE    = errors.New("not a bytecode executable file")
)

func (vm *VM) Run() *Value {
	var (
		acc *Value
		env *Value
		pc  int

		stk [500]*Value
		asp int

		rstk [100]returnFrame
		rsp  int
	)
	code := vm.code
	globalValue := vm.globalValue
	stop := false

	log.Debug(code)
	// log.Debug(globalValue)

	for !stop {
		op := code[pc]
		pc++
		switch op {
		case opACCESS:
			acc = builtinNth(env, int(code[pc]))
			log.Debug("ACCESS", code[pc])
			pc++
		case opADDINT:
			asp--
			log.Debug("ADDINT", acc.Int(), stk[asp].Int())
			acc = NewIntegerValue(acc.Int() + stk[asp].Int())
		case opSUBINT:
			asp--
			log.Debug("SUBINT")
			acc = NewIntegerValue(acc.Int() - stk[asp].Int())
		case opMULINT:
			asp--
			acc = NewIntegerValue(acc.Int() * stk[asp].Int())
		case opADDFLOAT:
			asp--
			log.Debug("ADDFLOAT")
			acc = NewFloatValue(acc.Float() + stk[asp].Float())
		case opPUSHMARK:
			log.Debug("PUSHMARK")
			stk[asp] = nil
			asp++
		case opPUSH:
			log.Debug("PUSH")
			stk[asp] = acc
			asp++
		case opRETURN:
			log.Debug("RETURN")
			if stk[asp-1] == nil {
				asp--
				rsp--
				pc = rstk[rsp].pc
				env = rstk[rsp].env
			}
			// TODO
		case opAPPLY:
			pushRetFrame(rstk[:], &rsp, pc, env)
			clo := acc.Closure()
			pc = clo.pc
			asp--
			log.Debug("APPLY")
			env = builtinCons(stk[asp], clo.env)
		case opGRAB:
			if stk[asp-1] == nil {
				log.Debug("GRAB not enough")
				asp--
				rsp--
				pc = rstk[rsp].pc
				env = rstk[rsp].env
			} else {
				asp--
				log.Debug("GRAB", stk[asp])
				env = builtinCons(stk[asp], env)
			}
		case opCONSTINT8:
			acc = NewIntegerValue(int(code[pc]))
			log.Debug("CONSTINT8", code[pc])
			pc++
		case opSETGLOBAL:
			alignPC(code, &pc)
			v := binary.LittleEndian.Uint16(code[pc:])
			log.Debug("SETGLOBAL", v)
			globalValue[v] = acc
			pc += 2
		case opGETGLOBAL:
			alignPC(code, &pc)
			n := pu16(code, pc)
			acc = globalValue[n]
			log.Debug("GETGLOBAL", n)
			pc += 2
		case opNOP:
			log.Debug("NOP", pc)
		case opSTOP:
			log.Debug("STOP")
			stop = true
		case opCUR:
			alignPC(code, &pc)
			ofst := int(pi16(code, pc))
			log.Debug("CUR", pc+ofst)
			acc = NewClosureValue(pc+ofst, env)
			pc += 2
		case opBRANCH:
			alignPC(code, &pc)
			pc += int(pi16(code, pc))
			log.Debug("BRANCH", pc)
		case opBRANCHIFNOT:
			log.Debug("BRANCHIFNOT")
			alignPC(code, &pc)
			if acc == False {
				pc += int(pi16(code, pc))
			} else {
				pc += 2
			}
			// case opMAKEBLOCK:
			// 	NewBlockValue(n int)
		case opCCALL2:
			idx := code[pc]
			log.Debug("OPCALL2", idx)
			fn := cprims[idx]
			asp--
			acc = fn(acc, stk[asp])
			pc++
		case opEQ:
			asp--
			acc = equal(acc, stk[asp])
		case opTERMAPPLY:
			clo := acc.Closure()
			pc = clo.pc
			log.Debug("TERMAPPLY", pc)
			asp--
			env = builtinCons(stk[asp], clo.env)
		default:
			log.Fatal("unknown instruct", op)
		}
	}

	return acc
}

func pi16(code []byte, pc int) int16 {
	return int16(binary.LittleEndian.Uint16(code[pc:]))
}

func pu16(code []byte, pc int) uint16 {
	return binary.LittleEndian.Uint16(code[pc:])
}

func alignPC(code []byte, pc *int) {
	if uintptr(unsafe.Pointer(&code[*pc]))&1 != 0 {
		*pc = *pc + 1
	}
}

func pushRetFrame(stk []returnFrame, sp *int, pc int, env *Value) {
	stk[*sp] = returnFrame{
		pc:  pc,
		env: env,
	}
	*sp = *sp + 1
}
