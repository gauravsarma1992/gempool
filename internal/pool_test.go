package internal

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	NodeOne struct {
		// Total bytes = 2
		AttrOne uint8
		AttrTwo uint8
	}
	NodeTwo struct {
		// Total bytes = 2 * 10 + 2 +
		AttrOne   [100]uint16
		AttrTwo   *NodeOne
		AttrThree [100]*NodeOne
	}
)

func TestPoolAllocWithStruct(t *testing.T) {
	blk, err := NewBlock()
	assert.Nil(t, err)
	bArr, err := blk.Alloc(reflect.TypeOf(NodeOne{}))
	assert.Nil(t, err)
	nodeOne := (*NodeOne)(bArr.ptr)
	nodeOne.AttrOne = 1
	nodeOne.AttrTwo = 2

	assert.Equal(t, uint8(1), nodeOne.AttrOne)
	assert.Equal(t, uint8(2), nodeOne.AttrTwo)
}

func TestPoolAllocWithArrayOfStruct(t *testing.T) {
	blk, err := NewBlock()
	assert.Nil(t, err)
	bArr, err := blk.Alloc(reflect.TypeOf([100]NodeOne{}))
	assert.Nil(t, err)
	nodes := (*[100]NodeOne)(bArr.ptr)
	for _, node := range nodes {
		tmp := node
		tmp.AttrOne = 1
		tmp.AttrTwo = 2

		assert.Equal(t, uint8(1), tmp.AttrOne)
		assert.Equal(t, uint8(2), tmp.AttrTwo)
	}
}

func TestPoolAllocWithArrayOfStructWithOverflow(t *testing.T) {
	blk, err := NewBlock()
	assert.Nil(t, err)

	bArr, err := blk.Alloc(reflect.TypeOf([10000]NodeOne{}))
	if err != nil {
		assert.NotNil(t, err)
		assert.Nil(t, bArr.ptr)
	}

}

func TestPoolAllocWithStructMultipleItems(t *testing.T) {
	blk, err := NewBlock()
	assert.Nil(t, err)

	bArr, err := blk.Alloc(reflect.TypeOf(NodeOne{}))
	assert.Nil(t, err)
	nodeOne := (*NodeOne)(bArr.ptr)
	nodeOne.AttrOne = 1
	nodeOne.AttrTwo = 2

	assert.Equal(t, uint8(1), nodeOne.AttrOne)
	assert.Equal(t, uint8(2), nodeOne.AttrTwo)

	bArr, err = blk.Alloc(reflect.TypeOf(NodeTwo{}))
	assert.Nil(t, err)
	nodeTwo := (*NodeTwo)(bArr.ptr)
	nodeTwo.AttrOne = [100]uint16{}
	nodeTwo.AttrTwo = nodeOne
	//nodeTwo.AttrThree = [100]*NodeOne{}

	assert.Equal(t, [100]uint16{}, nodeTwo.AttrOne)
	assert.Equal(t, uint8(1), nodeTwo.AttrTwo.AttrOne)
}

func BenchmarkPoolAlloc(b *testing.B) {
	blk, _ := NewBlock()
	nodeOneTyp := NodeOne{}
	nodeTwoTyp := NodeTwo{}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		blk.Reset()
		b.StartTimer()

		bArrOne, _ := blk.Alloc(reflect.TypeOf(nodeOneTyp))
		bArrTwo, _ := blk.Alloc(reflect.TypeOf(nodeTwoTyp))
		nodeOne := (*NodeOne)(bArrOne.ptr)
		nodeOne.AttrOne = 1
		nodeOne.AttrTwo = 2

		nodeTwo := (*NodeTwo)(bArrTwo.ptr)
		nodeTwo.AttrTwo = nodeOne
		for idx := 0; idx < 100; idx++ {
			nodeTwo.AttrOne[idx] = uint16(idx)
			nodeTwo.AttrThree[idx] = &NodeOne{}
		}
	}
}
