package rarch_parser

func (n node32) LenChildren() uint {
	if n.up == nil {
		return 0
	}
	child := n.up
	var children uint = 0
	for ; child != nil; child = child.next {
		children++
	}
	return children
}

func (n *node32) Ident(buffer string) string {
	return buffer[n.token32.begin:n.token32.end]
}
