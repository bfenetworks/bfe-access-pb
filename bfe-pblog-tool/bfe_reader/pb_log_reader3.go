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

// pb_log_reader3.go - Read access log record from pb log file

package bfe_reader

import (
	"fmt"
	"os"
	"sync"

	bfe_access_pb "github.com/bfenetworks/bfe-access-pb/bfe_access_pb"
	"github.com/bfenetworks/go-lib/web-monitor/module_state2"
)

type PbLogReader struct {
	*LogFileReader
}

func NewPbLogReader(logPath string, state *module_state2.State, clusterName string) *PbLogReader {
	pbr := &PbLogReader{newLogFileReader(logPath, state, clusterName)}
	pbr.dataBuffer = make([]byte, 0)
	return pbr
}

// parse records from the data buffer
func (lr *PbLogReader) dataBufferParse() []*bfe_access_pb.BfeLog {
	var records []*bfe_access_pb.BfeLog

	// parse pb record from buffer
	records, lr.dataBuffer = pbBuffParse(lr.dataBuffer, lr.state)

	lr.state.Inc("SUM_PB_RECORD", len(records))

	return records
}

/*
logRead - Read data from bp log file

Returns:

	(pbRecords, error)
*/
func (lr *PbLogReader) logRead() ([]*bfe_access_pb.BfeLog, error) {
	var records []*bfe_access_pb.BfeLog
	var data []byte
	var err error
	var hasNewLog bool

	// check whether file is open
	if lr.logFd == nil {
		// check whether file exist
		_, err = os.Stat(lr.logPath)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not exit: %s", lr.logPath)
		}

		// if not open, or sth goes wrong
		lr.dataBuffer = nil
		lr.logRelocate()

		if lr.logFd == nil {
			// Error happens
			return nil, fmt.Errorf("logRelocate() fail")
		}
	}

	// read data from opened log file
	data, err = lr.fileRead(MAX_BUFF_SIZE)
	lr.state.Inc("SUM_READ_DATA", 1)
	if err != nil {
		// Error happens
		return nil, fmt.Errorf("fileRead() fail:%s", err.Error())
	}

	if len(data) == 0 {
		lr.state.Inc("SUM_READ_DATA_EMPTY", 1)
		records = make([]*bfe_access_pb.BfeLog, 0)

		// End Of File
		hasNewLog, data, err = lr.eofHandler()
		if err == nil && len(data) != 0 {
			lr.dataBuffer = append(lr.dataBuffer, data...)

			// parse records from the data buffer
			records = lr.dataBufferParse()
		}

		// clear the read buffer, if there is new log file
		if hasNewLog {
			lr.dataBuffer = nil
		}
	} else {
		lr.dataBuffer = append(lr.dataBuffer, data...)

		// parse records from the data buffer
		records = lr.dataBufferParse()
	}

	return records, nil
}

func PblogCat(fp string, callback func(result []string)) error {
	pbr := NewPbLogReader(fp, nil, "")
	err := pbr.logFileOpen()
	defer pbr.closeFileAndInit()
	if err != nil {
		return err
	}
	_, err = pbr.logFd.Seek(0, 0)
	if err != nil {
		return err
	}

	bufChan := make(chan []byte, 1024)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			buf, err := pbr.fileRead(8192)
			if err != nil || len(buf) == 0 {
				close(bufChan)
				return
			}
			bufChan <- buf
		}
	}()

	go func() {
		defer wg.Done()
		for buf := range bufChan {
			pbr.dataBuffer = append(pbr.dataBuffer, buf...)
			records := pbr.dataBufferParse()
			var res = make([]string, 0)
			for _, r := range records {
				line := r.String()
				res = append(res, line)
			}
			callback(res)
		}
	}()

	wg.Wait()
	return nil
}
