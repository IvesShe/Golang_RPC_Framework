package rpc

import (
	"encoding/gob"
	"fmt"
	"net"
	"testing"
)

// 用戶查詢

// 用於測試的結構體
// 字段首字母必須大寫
type User struct {
	Name string
	Age  int
}

// 用於測試的查詢用戶的方法
func queryUser(uid int) (User, error) {
	user := make(map[int]User)
	user[0] = User{"ives", 20}
	user[1] = User{"Tom", 18}
	user[2] = User{"Jack", 30}

	// 模擬查詢用戶
	if u, ok := user[uid]; ok {
		return u, nil
	}

	return User{}, fmt.Errorf("id %d not in user db", uid)
}

// 測試
func TestRPC(t *testing.T) {
	// 需要對interface{}可能產生的類型進行註冊
	gob.Register(User{})

	addr := "127.0.0.1:8080"

	// 創建服務端
	srv := NewServer(addr)

	// 將方法註冊到服務端
	srv.Register("queryUser", queryUser)

	// 服務端等待調用
	go srv.Run()

	// 客戶端獲取連接
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
	}

	// 創建客戶端
	cli := NewClient(conn)

	// 聲明函數原型
	var query func(int) (User, error)
	cli.callRPC("queryUser", &query)

	// 得到查詢結果
	u, err := query(2)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(u)
}
