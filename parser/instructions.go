package rarch_parser

import (
	"log"
	"strconv"
	"unicode"

	rarch_encoding "github.com/jmbert/rarch-encoding"
)

func GenRegister(reg string) (byte, uint8) {
	var r rarch_encoding.Register

	if reg[0] != '$' {
		log.Fatalf("Invalid register %s\n", reg)
	}

	reg = reg[1:]

	switch reg[0] {
	case 'r':
		r.Register_type = rarch_encoding.RegisterGP
	case 'i':
		r.Register_type = rarch_encoding.RegisterIndex
	case 's':
		if reg[1] == 'p' {
			r.Register_type = rarch_encoding.RegisterGP
			break
		}
		r.Register_type = rarch_encoding.RegisterMMove
	case 'd':
		r.Register_type = rarch_encoding.RegisterMMove
	case 'c':
		r.Register_type = rarch_encoding.RegisterControl
	case 'b':
		r.Register_type = rarch_encoding.RegisterGP
	default:
		log.Fatalf("Invalid register %s\n", reg)
	}

	if unicode.IsDigit(rune(reg[1])) { // Numbered register
		var regNum uint64
		var err error
		if r.Register_type != rarch_encoding.RegisterGP || unicode.IsDigit(rune(reg[len(reg)-1])) { // We don't need to account for size-specifiers
			regNum, err = strconv.ParseUint(reg[1:], 10, 8)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			regNum, err = strconv.ParseUint(reg[1:len(reg)-1], 10, 8)
			if err != nil {
				log.Fatal(err)
			}
		}

		r.Register = byte(regNum)
	} else { // Stack or Base pointer
		if reg == "sp" {
			r.Register = rarch_encoding.REG_SP
		} else if reg == "bp" {
			r.Register = rarch_encoding.REG_BP
		}
	}

	if r.Register_type == rarch_encoding.RegisterGP && !unicode.IsDigit(rune(reg[len(reg)-1])) { // Deal with size specifier
		switch reg[len(reg)-1] {
		case 'b':
			r.Register_size = rarch_encoding.RegisterLen8
		case 's':
			r.Register_size = rarch_encoding.RegisterLen16
		case 'd':
			r.Register_size = rarch_encoding.RegisterLen32
		default:
			log.Fatalf("Unkown ending prefix %c\n", reg[len(reg)-1])
		}
	} else {
		r.Register_size = rarch_encoding.RegisterLen64
	}

	return r.Encode(), r.Register_size
}

func GenFormatA(p *Parser, n *node32, opcode, prefix byte) rarch_encoding.FormatA {
	return rarch_encoding.FormatA{Prefix: prefix, Opcode: opcode}
}
func GenFormatB(p *Parser, n *node32, opcode, prefix byte) rarch_encoding.FormatB {
	var err error
	var instr rarch_encoding.FormatB
	instr.Prefix = prefix
	instr.Opcode = opcode
	if n.up.next.up == nil {
		log.Fatalf("Bad AST structure")
	}
	instr.Register, _ = GenRegister(n.up.next.up.Ident(p.Buffer))
	if err != nil {
		log.Fatal(err)
	}

	return instr
}
func GenFormatC(p *Parser, n *node32, opcode, prefix byte, length_override uint8) rarch_encoding.FormatC {
	var err error
	var instr rarch_encoding.FormatC
	instr.Prefix = prefix
	instr.Opcode = opcode
	if n.up.next.up == nil {
		log.Fatalf("Bad AST structure")
	}
	if n.up.next.up.next == nil {
		log.Fatalf("Bad AST structure")
	}
	instr.Register, instr.Immediate.Length = GenRegister(n.up.next.up.Ident(p.Buffer))
	instr.Immediate.Value, err = strconv.ParseUint(n.up.next.up.next.Ident(p.Buffer), 0, 64)
	if err != nil {
		log.Fatal(err)
	}

	if length_override != 0 {
		instr.Immediate.Length = length_override - 1
	}

	return instr
}
func GenFormatD(p *Parser, n *node32, opcode, prefix byte) rarch_encoding.FormatD {
	// TODO
	return rarch_encoding.FormatD{}
}

func GenInstruction(p *Parser, n *node32) {
	var instr rarch_encoding.Instruction

	if n.up == nil {
		log.Fatalf("Bad AST structure")
	}
	switch n.up.Ident(p.Buffer) {
	case "ld":
		instr = GenFormatC(p, n, rarch_encoding.OP_LD, rarch_encoding.PREF_ABSOLUTE, rarch_encoding.ImmediateLen64+1)
	case "ldlit":
		instr = GenFormatC(p, n, rarch_encoding.OP_LD, rarch_encoding.PREF_LITERAL, 0)
	case "ldrel":
		instr = GenFormatC(p, n, rarch_encoding.OP_LD, rarch_encoding.PREF_PCREL, rarch_encoding.ImmediateLen64+1)
	case "ldind":
		instr = GenFormatC(p, n, rarch_encoding.OP_LD, rarch_encoding.PREF_INDREL, rarch_encoding.ImmediateLen64+1)
	case "hlt":
		instr = GenFormatA(p, n, rarch_encoding.OP_HLT, 00)
	case "nop":
		instr = GenFormatA(p, n, rarch_encoding.OP_NOP, 00)
	case "sind":
		instr = GenFormatB(p, n, rarch_encoding.OP_SIND, 00)
	case "soff":
		instr = GenFormatB(p, n, rarch_encoding.OP_SOFF, 00)
	case "jmp":
		instr = GenFormatB(p, n, rarch_encoding.OP_JMP, 00)
	default:
		log.Fatalf("Unknown instruction \"%s\"\n", n.up.Ident(p.Buffer))
	}
	current_reference_address += uint(len(instr.Encode()))
	program_bytes = append(program_bytes, instr.Encode()...)
}
