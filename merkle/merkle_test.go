package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMerkleNode(t *testing.T) {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
		[]byte("node4"),
		[]byte("node5"),
		[]byte("node6"),
		[]byte("node7"),
	}

	// level 1
	mn1 := CreateBaseMerkleNode(data[0])
	mn2 := CreateBaseMerkleNode(data[1])
	mn3 := CreateBaseMerkleNode(data[2])
	mn4 := CreateBaseMerkleNode(data[3])
	mn5 := CreateBaseMerkleNode(data[4])
	mn6 := CreateBaseMerkleNode(data[5])

	// level 2
	mn7 := CreateMerkleNode(mn1, mn2)
	mn8 := CreateMerkleNode(mn3, mn4)
	mn9 := CreateMerkleNode(mn5, mn6)
	mn10 := CreateBaseMerkleNode(data[6])

	//level 3
	mn11 := CreateMerkleNode(mn7, mn8)
	mn12 := CreateMerkleNode(mn9, mn10)

	//level 4
	mn15 := CreateMerkleNode(mn11, mn12)

	//Create Merkle
	merkleTree := CreateMerkleTree(data)

	assert.Equal(t, mn15.Data, merkleTree.RootNode.Data)
}
