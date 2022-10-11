// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"github.com/ethereum/go-ethereum/common"
)

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 20
	// Storage value length
	StorageUnitLength = 32
)

/*
需要结合opSload和opSstore两个指令的具体实现进行设计：
opSload：
1.首先从栈中读取状态变量的存储位置-loc （loc是一个地址，或者变量引用，loc.Bytes32()是取出该地址处的32字节的数据，SetBytes()写数据）
2.计算loc对应的hash值
3.下一步val的加载，是直接从evm.StateDB.中去读取的对应地址的hash表达的位置的值，简化情况下不考虑系统缓存，即是直接访问的状态数据库
4.最后给内存数据模型中对应位置的合约状态变量写入读出的值
func opSload(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
	loc := scope.Stack.peek()
	hash := common.Hash(loc.Bytes32())
	val := interpreter.evm.StateDB.GetState(scope.Contract.Address(), hash)
	loc.SetBytes(val.Bytes())
	return nil, nil
}

func opSstore(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
	if interpreter.readOnly {
		return nil, ErrWriteProtection
	}
	loc := scope.Stack.pop()
	val := scope.Stack.pop()
	interpreter.evm.StateDB.SetState(scope.Contract.Address(),
		loc.Bytes32(), val.Bytes32())
	return nil, nil
}

*/

/*
假设两个状态变量分别为uint64 a = 5; uint64 b = 16;
他们聚合存储在一个slot中(假设该slot在storage中的索引为1——第二个slot)，用“|”划分每64位【0|0|a=5|b=16】
则分别对应两个这样的MemStorageCache对象：
msc_a.sKey = 1; msc_a.sOffset = 8;msc_a.sValue = 5;
msc_b.sKey = 1; msc_b.sOffset = 0;msc_b.sValue = 16;
下面要考虑是否应该把sKey和sOffset的类型从[]byte换成其他更合适的，
比如上面hash计算的时候，loc有一个转换成Bytes32的过程，是不是sKey可以直接定义为[32]byte类型
*/
// type MemStorageCache struct {
// 	// sKey         []byte	// the key index of key-value database (storage)
// 	// sOffset   	 []byte	// the byte offset in one slot (a pair of key-value data), byte
// 	// sValue       []byte	// the value (the data of state variable) of key-value database (storage)
// 	map data [int][]Bytes32	// one slot : key -> value
// 	// lastGasCost  uint64
// }

// // Key returns the asked key
// func (msc *MemStorageCache) Key() []byte {
// 	return msc.sKey
// }

// // Offset returns the asked key
// func (msc *MemStorageCache) Offset() []byte {
// 	return msc.sOffset
// }

// // Value returns the asked value
// func (msc *MemStorageCache) Value() []byte {
// 	return msc.sValue
// }

// Address represents the 20 byte address of an Ethereum account.
// type Address [AddressLength]byte
type StorageKey [StorageUnitLength]byte

// type StorageValue [StorageUnitLength]byte
// type StorageValue []byte
// type StorageValue common.Hash

// StorageSlot represents one store unit in storage, a key-value database
// type StorageSlot map[StorageKey]common.Hash

// type StorageCache map[Address]StorageSlot
// type StorageCache map[common.Address]StorageSlot
type StorageCache map[common.Address]map[[32]byte]common.Hash

// MemStorageCache implements a cache memory model for the ethereum storage.
type MemStorageCache struct {
	data StorageCache // one slot : key -> value
}

// NewMemStorageCache returns a new MemStorageCache model.
func NewMemStorageCache() *MemStorageCache {
	return &MemStorageCache{}
}

// getValue returns the asked value with specific Address "addr" and StorageKey "key"
// func (msc *MemStorageCache) getValue(addr common.Address, key StorageKey) []byte {
func (msc *MemStorageCache) getValue(addr common.Address, key StorageKey) common.Hash {
	_, ok := msc.data[addr][key]
	if ok {
		return msc.data[addr][key]
	}
	return common.Hash{}
}

// setValue write StorageValue "value" into specific Address "addr" and StorageKey "key"
// func (msc *MemStorageCache) setValue(addr common.Address, key StorageKey, value StorageValue) []byte {
func (msc *MemStorageCache) setValue(addr common.Address, key StorageKey, value common.Hash) []byte {
	msc.data[addr][key] = value
	return nil
}

// preload implementation
func (msc *MemStorageCache) preload() []byte {
	// pass
	return nil
}

// preload implementation
func (msc *MemStorageCache) persist(interpreter *EVMInterpreter) {
	// pass
	for addr, skey := range msc.data {
		for subkey, val := range skey {
			interpreter.evm.StateDB.SetState(addr, subkey, val)
		}
	}
}
