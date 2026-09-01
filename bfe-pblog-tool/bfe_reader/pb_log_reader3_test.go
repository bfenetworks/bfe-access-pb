// Copyright(c) 2024 Beijing Yingfei Networks Technology Co.Ltd. All rights reserved.
package bfe_reader

import "testing"

func TestNewPbLogReader(t *testing.T) {
	pbr := NewPbLogReader("test_data/pb_access3.log", nil, "test1")
	err := pbr.logFileOpen()
	defer pbr.closeFileAndInit()
	if err != nil {
		t.Fatal(err)
	}
	_, err = pbr.logFd.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	buf, err := pbr.fileRead(0)
	if err != nil {
		t.Fatal(err)
	}
	pbr.dataBuffer = buf
	records := pbr.dataBufferParse()
	t.Log(len(records), len(buf))

	// Characterization assertions: lock wire-format compatibility across migration
	if len(records) != 9 {
		t.Fatalf("expected 9 records, got %d", len(records))
	}
	if records[0].GetLogid() == 0 {
		t.Fatal("expected records[0].Logid to be non-zero")
	}
	if records[0].GetLogTag() == "" {
		t.Fatal("expected records[0].LogTag to be non-empty")
	}

	for i, r := range records {
		t.Log(i, r)
	}
}
