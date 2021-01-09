package rpc

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

// 寫一個test程序

// 測試讀寫
func TestSession_ReadWrite(t *testing.T) {
	// 定義監聽ip和端口
	addr := "127.0.0.1:8000"

	// 定義傳輸的數據
	myData := "hello rpc"

	// 等待組
	wg := sync.WaitGroup{}

	// 協程 1個讀、1個寫
	wg.Add(2)

	// 寫數據協程
	go func() {
		defer wg.Done()

		// 創建tcp連接
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		conn, _ := lis.Accept()
		s := Session{conn: conn}

		// 寫數據
		err = s.Write([]byte(myData))
		if err != nil {
			t.Fatal(err)
		}
	}()

	// 讀數據協程
	go func() {
		defer wg.Done()

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}

		s := Session{conn: conn}

		// 讀數據
		data, err := s.Read()
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != myData {
			t.Fatal(err)
		}
		fmt.Println(string(data))
	}()

	wg.Wait()
}
