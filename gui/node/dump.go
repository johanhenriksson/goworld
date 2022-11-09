package node

import "fmt"

func Dump(node T) {
	dump(node, "")
}

func dump(node T, indent string) {
	fmt.Printf("%s%s\n", indent, node.Key())
	indent = indent + "  "
	for _, child := range node.Children() {
		dump(child, indent)
	}
}
