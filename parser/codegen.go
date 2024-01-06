package rarch_parser

import (
	"log"
)

type Label struct {
	address    uint64
	defined    bool
	references []uint
}

var current_reference_address uint = 0

var labels map[string]*Label = make(map[string]*Label)

var program_bytes []byte = make([]byte, 0)

var genmap map[pegRule]func(p *Parser, n *node32) = map[pegRule]func(p *Parser, n *node32){
	ruleDirective:   RunDirective,
	ruleLabel:       GenLabel,
	ruleInstruction: GenInstruction,
}

func GenLabel(p *Parser, n *node32) {
	ident := string(p.buffer[n.next.token32.begin:n.next.token32.end])
	l, present := labels[ident]
	if !present {
		l = &Label{uint64(current_reference_address), true, make([]uint, 0)}
		labels[ident] = l
	} else {
		l.address = uint64(current_reference_address)
		l.defined = true
		for _, ref := range l.references {
			// Resolve reference
			ref = ref
		}
	}
}

func RunDirective(p *Parser, n *node32) {
	directive := string(p.buffer[n.up.token32.begin:n.up.token32.end])
	args := n.up.next
	fn, ok := dirmap[directive]
	if !ok {
		log.Fatalf("No such directive: %s", directive)
	}
	fn(p, args)
}

func (p *Parser) gen(node *node32) {
	for ; node != nil; node = node.next {
		fn, ok := genmap[node.pegRule]
		if ok {
			fn(p, node)
		}
		if node.up != nil {
			p.gen(node.up)
		}
	}
}

func (p *Parser) Codegen() []byte {
	p.gen(p.AST())
	return program_bytes
}
