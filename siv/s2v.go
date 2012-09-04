// Copyright 2012 Aaron Jacobs. All Rights Reserved.
// Author: aaronjjacobs@gmail.com (Aaron Jacobs)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package siv

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"github.com/jacobsa/aes/cmac"
	"github.com/jacobsa/aes/common"
)

var s2vZero []byte

func init() {
	s2vZero = bytes.Repeat([]byte{0x00}, aes.BlockSize)
}

// Run the S2V "string to vector" function of RFC 5297 using the input key and
// string vector, which must be non-empty. (RFC 5297 defines S2V to handle the
// empty vector case, but it is never used that way by higher-level functions.)
func s2v(key []byte, strings [][]byte) []byte {
	numStrings := len(strings)
	if numStrings == 0 {
		panic("strings vector must be non-empty.")
	}

	// Create a CMAC hash.
	h, err := cmac.New(key)
	if err != nil {
		panic(fmt.Sprintf("cmac.New: %v", err))
	}

	// Initialize.
	h.Write(s2vZero)
	d := h.Sum([]byte{})
	h.Reset()

	// Handle all strings but the last.
	for i := 0; i < numStrings-1; i++ {
		h.Write(strings[i])
		d = common.Xor(dbl(d), h.Sum([]byte{}))
		h.Reset()
	}

	// Handle the last string.
	lastString := strings[numStrings-1]
	var t []byte
	if len(lastString) >= aes.BlockSize {
		t = xorend(lastString, d)
	} else {
		t = common.Xor(d, common.PadBlock(lastString))
	}

	h.Write(t)
	return h.Sum([]byte{})
}