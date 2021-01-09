package rpc

import (
	"fmt"
	"net"
	"reflect"
)

// 聲明服務端
type Server struct {
	// 地址
	addr string

	// 服務端維護的函數名到函數反射值的map
	funcs map[string]reflect.Value
}

// 創建服務端對象
func NewServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

// 服務端綁定註冊方法
// 將函數名與函數真正對應起來
// 第一個參數為函數名，第二個傳入真正函數
func (s *Server) Register(rpcName string, f interface{}) {
	if _, ok := s.funcs[rpcName]; ok {
		return
	}
	// map中沒有值，剛好映射添加進map，便於調用
	fVal := reflect.ValueOf(f)
	s.funcs[rpcName] = fVal
}

// 服務端等待調用
func (s *Server) Run() {
	// 監聽
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Printf("監聽 %s err: %v", s.addr, err)
		return
	}

	for {
		// 拿到連接
		conn, err := lis.Accept()
		if err != nil {
			fmt.Printf("accept err: %v", err)
			return
		}

		// 創建會話
		srvSession := NewSession(conn)

		// RPC讀取數據
		b, err := srvSession.Read()
		if err != nil {
			fmt.Printf("read err: %v", err)
			return
		}

		// 對數據解碼
		rpcData, err := decode(b)
		if err != nil {
			fmt.Printf("decode err: %v", err)
			return
		}

		// 根據讀取到的數據的Name，得到調用的函數名
		f, ok := s.funcs[rpcData.Name]
		if !ok {
			fmt.Printf("函數 %s 不存在", rpcData.Name)
			return
		}

		// 解析遍歷客戶端出來的參數，放到一個數組中
		inArgs := make([]reflect.Value, 0, len(rpcData.Args))
		for _, arg := range rpcData.Args {
			inArgs = append(inArgs, reflect.ValueOf(arg))
		}

		// 反射調用方法，傳入參數
		out := f.Call(inArgs)

		// 解析遍歷執行結果，放到一個數組中
		outArgs := make([]interface{}, 0, len(out))
		for _, o := range out {
			outArgs = append(outArgs, o.Interface())
		}

		// 包裝數據，返回給客戶端
		respRPCData := RPCData{rpcData.Name, outArgs}

		// 編碼
		respBytes, err := encode(respRPCData)
		if err != nil {
			fmt.Printf("encode err: %v", err)
			return
		}

		// 使用RPC寫出數據
		err = srvSession.Write(respBytes)
		if err != nil {
			fmt.Printf("session write err: %v", err)
			return
		}
	}
}
