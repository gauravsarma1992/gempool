package internal

import (
	"fmt"
	"reflect"
	"unsafe"

	"math/rand"
)

const (
	BufferSize = 8096
)

type (
	Pool struct {
		blks []*Block
	}

	Block struct {
		id uint64

		addrs     [BufferSize]*Addr
		addrCount uint64

		layout  [BufferSize]byte
		currPos uint64

		bufSz uint64
	}

	Addr struct {
		startAddr uint64
		stopAddr  uint64

		ptr unsafe.Pointer
	}
)

func NewBlock() (blk *Block, err error) {
	// Check if memory is available here
	blk = &Block{
		id:    uint64(rand.Uint64()),
		bufSz: BufferSize,
	}
	blk.Reset()
	return
}

func (blk *Block) Alloc(typ reflect.Type) (addr *Addr, err error) {
	sz := uint64(typ.Size())
	addr = blk.addrs[blk.addrCount]
	blk.addrCount += 1

	addr.startAddr = blk.currPos + 1
	addr.stopAddr = blk.currPos + 1 + sz

	if addr.stopAddr > uint64(len(blk.layout)) {
		err = fmt.Errorf("no size left for allocation on block %d for elem of sz %d", blk.id, sz)
		return
	}
	//log.Printf("Total bytes allocated from the block %d - %d. Bytes left - %d\n", blk.id, sz, blk.bufSz-blk.currPos)
	addr.ptr = unsafe.Pointer(&blk.layout[addr.startAddr:addr.stopAddr][0])
	blk.currPos = addr.stopAddr
	return
}

func (blk *Block) Free(addr *Addr) (err error) {
	return
}

func (blk *Block) Reset() (err error) {
	blk.currPos = 0
	blk.addrCount = 0
	for idx := 0; idx < BufferSize; idx++ {
		blk.addrs[idx] = &Addr{}
		blk.layout[idx] = byte(0)
	}
	return
}
