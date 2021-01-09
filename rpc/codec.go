package rpc

import (
	"bytes"
	"encoding/gob"
)

// 定義數據格式和編解碼

// 定義RPC交互的數據格式
type RPCData struct {
	// 訪問的函數
	Name string

	// 訪問時候的參數
	Args []interface{}
}

// 編碼
func encode(data RPCData) ([]byte, error) {
	var buf bytes.Buffer

	// 得到字節數組的編碼器
	bufEnc := gob.NewEncoder(&buf)

	// 對數據編碼
	if err := bufEnc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 解碼
func decode(b []byte) (RPCData, error) {
	buf := bytes.NewBuffer(b)

	// 返回字節數組解碼器
	bufDec := gob.NewDecoder(buf)
	var data RPCData

	// 對數據解碼
	if err := bufDec.Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}
