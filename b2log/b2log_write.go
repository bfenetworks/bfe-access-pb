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
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// generate timestamp
func timestampGen() uint64 {
	t := time.Now()
	sec := t.Unix()
	usec := t.Nanosecond() / 1000
	ts := uint64(sec*1000 + int64(usec)/1000)

	return ts
}

/*
HeaderWrite - write b2log header to given buffer

Params:
  - buffer: []byte to write header to
  - payloadLen: length of payload
*/
func HeaderWrite(buffer []byte, payloadLen int) error {
	// prepare header
	header := Header{MAGIC_NUMBER, HEADER_VERSION, uint32(payloadLen), 0, 0}
	header.TimeStamp = timestampGen()

	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.LittleEndian, header)
	if err != nil {
		return fmt.Errorf("binary.Write():%s", err.Error())
	}

	// write header to buffer
	copy(buffer, buff.Bytes())

	return nil
}
