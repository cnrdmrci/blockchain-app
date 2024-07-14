package merkle

import (
	"blockchain-app/handlers"
	"crypto/sha256"
	"errors"
)

type Tree struct {
	RootNode *Node
}

type Node struct {
	Left  *Node
	Right *Node
	Data  []byte
}

func CreateBaseMerkleNode(data []byte) *Node {
	node := Node{}
	hash := sha256.Sum256(data)
	node.Data = hash[:]

	return &node
}

func CreateMerkleNode(left, right *Node) *Node {
	node := Node{}

	if left == nil || right == nil {
		handlers.HandleErrors(errors.New("merkle node is nil"))
	}

	prevHashes := append(left.Data, right.Data...)
	hash := sha256.Sum256(prevHashes)
	node.Data = hash[:]
	node.Left = left
	node.Right = right

	return &node
}

func CreateMerkleTree(dataArr [][]byte) *Tree {

	if len(dataArr) == 0 {
		handlers.HandleErrors(errors.New("there is no data while creating a new merkle tree"))
	}

	var nodes []Node
	for _, data := range dataArr {
		node := CreateBaseMerkleNode(data)
		nodes = append(nodes, *node)
	}

	for len(nodes) > 1 {
		var level []Node
		for i := 1; i < len(nodes); i += 2 {
			node := CreateMerkleNode(&nodes[i-1], &nodes[i])
			level = append(level, *node)
		}

		if len(nodes)%2 != 0 {
			level = append(level, nodes[len(nodes)-1])
		}

		nodes = level
	}

	return &Tree{&nodes[0]}
}
