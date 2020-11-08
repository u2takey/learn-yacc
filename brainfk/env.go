package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type env struct {
	codes Codes
}

type Codes []Code

func (c Codes) String() string {
	b := strings.Builder{}
	for _, a := range c {
		b.WriteString(a.opType.String())
		if a.argument > 1 && a.opType != JUMP_IF_DATA_ZERO && a.opType != JUMP_IF_DATA_NOT_ZERO {
			b.WriteString(strconv.Itoa(a.argument))
		}
	}
	return b.String()
}

func newEnv(codes []Code) *env {
	return &env{
		codes: codes,
	}
}

func (e *env) transform() []Code {
	// for [ ]
	var ops Codes
	pc, pSize := 0, len(e.codes)
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
			openBracketStack = openBracketStack[:len(openBracketStack)-1]

			ops[openBracketOffset].argument = len(ops)
			toOptimize := ops[openBracketOffset:]
			ops = ops[:openBracketOffset]
			toOptimize = append(toOptimize, Code{opType: JUMP_IF_DATA_NOT_ZERO, argument: openBracketOffset})
			optimized := e.optimizeLoop(toOptimize)
			ops = append(ops, optimized...)

			pc += 1
		} else {
			// 重复操作优化
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

func (e *env) optimizeLoop(loop []Code) []Code {
	// 对循环进行优化，比如 `[+]` 或者 `[-]`, 表示持续加 1 直到变为0 可以优化成 `(setzero)`； `[>]` 或者 `[<]`,
	// 表示向左/右移动指针直到指针下面的值为 0，可以优化成 `(moveptr)`； `[-<+>]` 或者 `[->+<]`
	// 表示向指针下/上n个位置移动当前指针下的值，可以优化成 `(movedata, n)`
	switch len(loop) {
	case 3:
		switch loop[1].opType {
		case INC_DATA, DEC_DATA:
			return []Code{{opType: LOOP_SET_TO_ZERO}}
		case INC_PTR:
			return []Code{{opType: LOOP_MOVE_PTR, argument: loop[1].argument}}
		case DEC_PTR:
			return []Code{{opType: LOOP_MOVE_PTR, argument: -loop[1].argument}}
		}
	case 6:
		if loop[1].opType == DEC_DATA &&
			loop[3].opType == INC_DATA &&
			loop[1].argument == loop[3].argument &&
			loop[1].argument == 1 &&
			loop[2].argument == loop[4].argument {
			if loop[2].opType == DEC_PTR && loop[4].opType == INC_PTR {
				return []Code{{LOOP_MOVE_DATA, -loop[2].argument}}
			} else if loop[4].opType == DEC_PTR && loop[2].opType == INC_PTR {
				return []Code{{LOOP_MOVE_DATA, loop[2].argument}}
			}
		}
	}
	return loop
}

func (e *env) Run() {
	log.Println("before transform", e.codes)
	e.codes = e.transform()
	log.Println("after transform", e.codes)

	memory, dataPtr, pc := make([]byte, 1024*1024*20), 0, 0
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
