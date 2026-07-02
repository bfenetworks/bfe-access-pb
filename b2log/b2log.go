// Copyright (c) 2026 The BFE Authors.
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

package b2log

import (
	"unsafe"
)

// magic number
// This can only work in little-endian machine (e.g., x86)
// Remember to make MAGIC_NUMBER and MAGIC_NUMBER_STR consistent
const (
	MAGIC_NUMBER   = 0xB0AEBEA7
	HEADER_VERSION = 1
)

var MAGIC_NUMBER_STR = []byte{0xA7, 0xBE, 0xAE, 0xB0}

// size of Header
var HEADER_SIZE = int(unsafe.Sizeof(demoHeader))
var demoHeader Header // this var is only for getting size of Header

const MAX_RECORD_LEN = 100 * 1024 // max length of single b2log record

// header for b2log record
type Header struct {
	MagicNumber   uint32 // magic number
	Version       uint32 // version
	UnCompressLen uint32 // length of upcompress log
	CompressLen   uint32 // length of compress log
	TimeStamp     uint64 // timestamp the log generated
}

// binary format of b2log record
type Record []byte
