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
	"testing"
)

// test of HeaderWrite()
func Test_HeaderWrite(t *testing.T) {
	// write b2log msg to buffer
	payload := []byte("this is a test")
	buff := make([]byte, HEADER_SIZE+len(payload))

	// write header
	err := HeaderWrite(buff, len(payload))
	if err != nil {
		t.Errorf("HeaderWrite():%s", err.Error())
		return
	}

	// write payload
	copy(buff[HEADER_SIZE:], payload)

	// try to read b2log msg from buffer
	records, buff := BuffParse(buff)
	if len(records) != 1 {
		t.Errorf("len(records) should be 1, now it's %d", len(records))
		return
	}
	if bytes.Compare(records[0], payload) != 0 {
		t.Errorf("records[0] = %s", records[0])
		return
	}
	if len(buff) != 0 {
		t.Errorf("len(buff) should be 0, now it's %d", len(buff))
	}
}
