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

// gen_testdata generates test data files with realistic BfeLog records
// in b2log binary format for testing bfePblogTool.
//
// Usage: go run ./cmd/gen_testdata [-o output_file] [-n count]

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"google.golang.org/protobuf/proto"

	bfe_access_pb "github.com/bfenetworks/bfe-access-pb/bfe_access_pb"
	"github.com/bfenetworks/bfe-access-pb/b2log"
)

func main() {
	outputFile := flag.String("o", "test_data/pb_access.log", "output file path")
	recordCount := flag.Int("n", 20, "number of records to generate")
	flag.Parse()

	records := generateRecords(*recordCount)

	if err := writeB2logFile(*outputFile, records); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %d records to %s\n", len(records), *outputFile)
}

func generateRecords(count int) []*bfe_access_pb.BfeLog {
	records := make([]*bfe_access_pb.BfeLog, 0, count)
	baseTime := time.Now().Add(-1 * time.Hour)

	products := []string{"BFE", "NEWS", "VIDEO", "MAP"}
	clusters := []string{"cluster_news", "cluster_video", "cluster_map"}
	uris := []string{"/", "/api/data", "/static/image.jpg", "/favicon.ico", "/video/stream"}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36",
		"curl/7.68.0",
	}

	for i := 0; i < count; i++ {
		ts := baseTime.Add(time.Duration(i) * time.Second)
		logid := uint64(ts.UnixNano()) + uint64(i*1000)

		product := products[i%len(products)]
		cluster := clusters[i%len(clusters)]
		uri := uris[i%len(uris)]
		ua := userAgents[i%len(userAgents)]

		logType := bfe_access_pb.BfeLogType_Request
		logTag := fmt.Sprintf("req_%s", product)

		addrInfo := &bfe_access_pb.ConnAddrInfo{
			BfeIp:       proto.Uint32(uint32(0x0A000101 + i%10)),
			SockSrcIp:   proto.Uint32(uint32(0xC0A80001 + i)),
			IsTrustSrcIp: proto.Bool(false),
		}

		reqLog := &bfe_access_pb.RequestLog{
			ErrCode:        proto.String(""),
			ErrMsg:         proto.String(""),
			ReqHeaderLen:   proto.Uint32(uint32(500 + i*10)),
			ReqBodyLen:     proto.Uint32(0),
			SessionId:      proto.Uint64(uint64(1000000 + i)),
			AddrInfo:       addrInfo,
			ClientIp:       proto.Uint32(uint32(0xC0A80001 + i)), // 192.168.0.1+
			ReqNum:         proto.Uint32(1),
			Proto:          proto.String("http/1.1"),
			HeaderHost:     proto.String(fmt.Sprintf("%s.example.com", product)),
			OriginUri:      proto.String(uri),
			UserAgent:      proto.String(ua),
			Product:        proto.String(product),
			Cluster:        proto.String(cluster),
			SubCluster:     proto.String(fmt.Sprintf("%s-%d", cluster, i%3)),
			ResStatusCode:  proto.Uint32(200),
			ResHeaderLen:   proto.Uint32(uint32(200 + i*5)),
			ResBodyLen:     proto.Uint32(uint32(1000 + i*100)),
			ResContentType: proto.String("text/html"),
			AllTime:        proto.Uint32(uint32(10 + i*2)),
			ReadClientTime: proto.Uint32(0),
			ClusterServeTime:   proto.Uint32(uint32(8 + i*2)),
			BackendServeTime:   proto.Uint32(uint32(5 + i)),
			WriteClientTime:    proto.Uint32(1),
			SessionOffsetTime:  proto.Uint32(uint32(100 + i*10)),
		}

		record := &bfe_access_pb.BfeLog{
			Product:    bfe_access_pb.ProductID_BFE.Enum(),
			Timestamp:  proto.Uint64(uint64(ts.Unix())),
			Logid:      proto.Uint64(logid),
			LogTag:     proto.String(logTag),
			LogType:    &logType,
			RequestLog: reqLog,
		}

		records = append(records, record)
	}

	return records
}

func writeB2logFile(filename string, records []*bfe_access_pb.BfeLog) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	for _, record := range records {
		data, err := proto.Marshal(record)
		if err != nil {
			return fmt.Errorf("marshal protobuf: %w", err)
		}

		buf := make([]byte, b2log.HEADER_SIZE+len(data))
		if err := b2log.HeaderWrite(buf, len(data)); err != nil {
			return fmt.Errorf("write header: %w", err)
		}
		copy(buf[b2log.HEADER_SIZE:], data)

		if _, err := f.Write(buf); err != nil {
			return fmt.Errorf("write to file: %w", err)
		}
	}

	return nil
}
