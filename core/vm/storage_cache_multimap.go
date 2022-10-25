// in-memory storage cache implementation
// data structure of this custom designed storage cache
package vm

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// the mixed key which is built by  jointing the address and storage key of one state variable value
type AddrKey struct {
	addr common.Address
	key  [32]byte
}

type MemStorageCache struct {
	data map[AddrKey]common.Hash
}

// NewMemStorageCache returns a new MemStorageCache model.
func NewMemStorageCache() *MemStorageCache {
	return &MemStorageCache{}
}

// getValue returns the asked value with specific Address "addr" and StorageKey "key"
func (msc *MemStorageCache) getValue(addr common.Address, key [32]byte) common.Hash {
	_, ok := msc.data[AddrKey{addr, key}]
	if ok {
		return msc.data[AddrKey{addr, key}]
	}
	return common.Hash{}
}

// setValue write StorageValue "value" into specific Address "addr" and StorageKey "key"
func (msc *MemStorageCache) setValue(addr common.Address, key [32]byte, value common.Hash) []byte {
	if msc.data == nil {
		msc.data = make(map[AddrKey]common.Hash)
	}
	// _, ok := msc.data[AddrKey{addr, key}]
	// if ok {
	// 	msc.data[AddrKey{addr, key}] = value
	// 	return nil
	// }
	msc.data[AddrKey{addr, key}] = value
	return nil
}

// preload implementation
func (msc *MemStorageCache) preload(scope *ScopeContext) []byte {
	log.Info("Function preload is executting!")
	msc.data = make(map[AddrKey]common.Hash)
	// pass
	return nil
}

// persist implementation
func (msc *MemStorageCache) persist(interpreter *EVMInterpreter) {
	// pass
	m_log, _err := os.OpenFile("./dversion0/1025/persist_MemStorageCache_2.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if _err != nil {
		fmt.Println(_err.Error())
	}
	defer m_log.Close()

	message := "Current executing persist process!"
	m_log.WriteString(message + "\n")
	for addrk, val := range msc.data {
		addr, subkey := addrk.addr, addrk.key
		interpreter.evm.StateDB.SetState(addr, subkey, val)
		// 打印输出，查看MemStorageCache的状态
		message = fmt.Sprintf("Address: %v	\nloc_key= %v  \nloc_value= %v \n\n\n", addr.Bytes(), subkey, val.Bytes())
		m_log.WriteString(message)
	}
}
