package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type env struct {
	codes []Code
}

func newEnv(codes []Code) *env {
	return &env{
		codes: codes,
	}
}

func (e *env) transform() []Code {
	// for [ ]
	pc, pSize, ops := 0, len(e.codes), []Code{}
	var openBracketStack []int
	for pc < pSize {
		code := e.codes[pc]
		if code.opType == JUMP_IF_DATA_ZERO {
			openBracketStack = append(openBracketStack, len(ops))
			ops = append(ops, Code{opType: JUMP_IF_DATA_ZERO})
			pc += 1
		} else if code.opType == JUMP_IF_DATA_NOT_ZERO {
			if len(openBracketStack) == 0 {
				log.Panic("unmatched closing ']' at pc=", pc)
			}
			openBracketOffset := openBracketStack[len(openBracketStack)-1]
			ops[openBracketOffset].argument = len(ops)
			ops = append(ops, Code{opType: JUMP_IF_DATA_NOT_ZERO, argument: openBracketOffset})
			openBracketStack = openBracketStack[:len(openBracketStack)-1]
			pc += 1
		} else {

			start := pc
			pc += 1
			for pc < pSize && e.codes[pc].opType == code.opType {
				pc += 1
			}
			numRepeats := pc - start
			ops = append(ops, Code{opType: code.opType, argument: numRepeats})
		}
	}

	return ops
}

func (e *env) Run() {
	e.codes = e.transform()
	// log.Println("after transform", e.codes)

	memory, dataPtr, pc := make([]byte, 2000), 0, 0
	reader := bufio.NewReader(os.Stdin)
	for pc < len(e.codes) {
		code := e.codes[pc]
		// log.Println(code, dataPtr, pc)
		switch code.opType {
		case INC_PTR:
			dataPtr += code.argument
		case DEC_PTR:
			dataPtr -= code.argument
		case INC_DATA:
			memory[dataPtr] += byte(code.argument)
		case DEC_DATA:
			memory[dataPtr] -= byte(code.argument)
		case READ_STDIN:
			b, err := reader.ReadByte()
			if err != nil {
				panic(err)
			}
			memory[dataPtr] += b
		case WRITE_STDOUT:
			fmt.Printf("%s", strings.Repeat(string(memory[dataPtr]), code.argument))
		case LOOP_SET_TO_ZERO:
			memory[dataPtr] = 0
		case LOOP_MOVE_PTR:
			for memory[dataPtr] != 0 {
				dataPtr += code.argument
			}
		case LOOP_MOVE_DATA:
			if memory[dataPtr] != 0 {
				memory[dataPtr+code.argument] += memory[dataPtr]
				memory[dataPtr] = 0
			}
		case JUMP_IF_DATA_ZERO:
			if memory[dataPtr] == 0 {
				pc = code.argument
			}
		case JUMP_IF_DATA_NOT_ZERO:
			if memory[dataPtr] != 0 {
				pc = code.argument
			}
		default:
			log.Panic("INVALID_OP encountered on pc=", pc)
		}
		pc += 1
	}

}
