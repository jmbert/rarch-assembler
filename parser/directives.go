package rarch_parser

import (
	"log"
	"strconv"
)

var dirmap map[string]func(p *Parser, n *node32) = map[string]func(p *Parser, n *node32){
	"org":    Org,
	"fillto": Fillto,
	"byte":   Byte,
	"short":  Short,
	"double": Double,
	"quad":   Quad,
}

func Byte(p *Parser, args *node32) {
	if args == nil {
		log.Fatalf("Incorrect number of arguments")
	}
	if args.LenChildren() != 1 {
		log.Fatalf("Incorrect number of arguments")
	}

	var b uint8
	switch args.up.up.pegRule {
	case ruleImmediate:
		shortim, err := strconv.ParseUint(string(p.buffer[args.up.up.token32.begin:args.up.up.token32.end]), 0, 8)
		if err != nil {
			log.Fatal(err)
		}
		b = uint8(shortim)
	default:
		log.Fatalln("Couldn't parse argument")
	}

	program_bytes = append(program_bytes, byte(b))
	current_reference_address += 1
}

func Short(p *Parser, args *node32) {
	if args == nil {
		log.Fatalf("Incorrect number of arguments")
	}
	if args.LenChildren() != 1 {
		log.Fatalf("Incorrect number of arguments")
	}

	var short uint16
	switch args.up.up.pegRule {
	case ruleImmediate:
		shortim, err := strconv.ParseUint(string(p.buffer[args.up.up.token32.begin:args.up.up.token32.end]), 0, 16)
		if err != nil {
			log.Fatal(err)
		}
		short = uint16(shortim)
	default:
		log.Fatalln("Couldn't parse argument")
	}

	program_bytes = append(program_bytes, byte(short>>8), byte(short))
	current_reference_address += 2
}

func Double(p *Parser, args *node32) {
	if args == nil {
		log.Fatalf("Incorrect number of arguments")
	}
	if args.LenChildren() != 1 {
		log.Fatalf("Incorrect number of arguments")
	}

	var double uint32
	switch args.up.up.pegRule {
	case ruleImmediate:
		doubleim, err := strconv.ParseUint(string(p.buffer[args.up.up.token32.begin:args.up.up.token32.end]), 0, 32)
		if err != nil {
			log.Fatal(err)
		}
		double = uint32(doubleim)
	default:
		log.Fatalln("Couldn't parse argument")
	}

	program_bytes = append(program_bytes, byte(double>>24), byte(double>>16), byte(double>>8), byte(double))
	current_reference_address += 4
}

func Quad(p *Parser, args *node32) {
	if args == nil {
		log.Fatalf("Incorrect number of arguments")
	}
	if args.LenChildren() != 1 {
		log.Fatalf("Incorrect number of arguments")
	}

	var double uint64
	switch args.up.up.pegRule {
	case ruleImmediate:
		doubleim, err := strconv.ParseUint(string(p.buffer[args.up.up.token32.begin:args.up.up.token32.end]), 0, 64)
		if err != nil {
			log.Fatal(err)
		}
		double = uint64(doubleim)
	default:
		log.Fatalln("Couldn't parse argument")
	}

	program_bytes = append(program_bytes, byte(double>>56), byte(double>>48), byte(double>>40), byte(double>>32), byte(double>>24), byte(double>>16), byte(double>>8), byte(double))
	current_reference_address += 8
}

func Org(p *Parser, args *node32) {
	if args == nil {
		log.Fatalf("Incorrect number of arguments")
	}
	if args.LenChildren() != 1 {
		log.Fatalf("Incorrect number of arguments")
	}
	var addr uint
	switch args.up.up.pegRule {
	case ruleImmediate:
		addrim, err := strconv.ParseUint(string(p.buffer[args.up.up.token32.begin:args.up.up.token32.end]), 0, 64)
		if err != nil {
			log.Fatal(err)
		}
		addr = uint(addrim)
	default:
		log.Fatalln("Couldn't parse argument")
	}

	current_reference_address = addr
}

func Fillto(p *Parser, args *node32) {
	if args == nil {
		log.Fatalf("Incorrect number of arguments")
	}
	if args.LenChildren() != 2 {
		log.Fatalf("Incorrect number of arguments")
	}
	var toAddr uint
	var fillByte uint8
	switch args.up.up.pegRule {
	case ruleImmediate:
		addrim, err := strconv.ParseUint(string(p.buffer[args.up.up.token32.begin:args.up.up.token32.end]), 0, 64)
		if err != nil {
			log.Fatal(err)
		}
		toAddr = uint(addrim)
	default:
		log.Fatalln("Couldn't parse argument")
	}
	switch args.up.next.up.pegRule {
	case ruleImmediate:
		fillim, err := strconv.ParseUint(string(p.buffer[args.up.next.up.token32.begin:args.up.next.up.token32.end]), 0, 64)
		if err != nil {
			log.Fatal(err)
		}
		fillByte = uint8(fillim)
	default:
		log.Fatalln("Couldn't parse argument")
	}

	for ; current_reference_address < toAddr; current_reference_address++ {
		program_bytes = append(program_bytes, fillByte)
	}

}
