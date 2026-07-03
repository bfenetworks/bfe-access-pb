# 项目名称
b2log(BFE access to log), read and write b2log record

## 快速开始
### read record
	//....	
    // read testing data from file
    data, err := ioutil.ReadFile("pb_access.log")
    if err != nil {
        t.Error("fail to open file for testing data")
        return
    }
    
    // parse b2log record from data
    records, buffer := b2log.BuffParse(data)

### write record
	//...

	// write b2log msg to buffer
	payload := []byte("this is a test")
	buff := make([]byte, b2log.HEADER_SIZE + len(payload))
	
	// write header
	err := b2log.HeaderWrite(buff, len(payload))
	if err != nil {
		t.Errorf("b2log.HeaderWrite():%s", err.Error())
		return
	}
	
	// write payload
	copy(buff[b2log.HEADER_SIZE:], payload)
	
	// write buffer to file, omitted
	
