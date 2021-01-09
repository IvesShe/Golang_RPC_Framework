package rpc

import (
	"net"
	"reflect"
)

// 聲明客戶端
type Client struct {
	conn net.Conn
}

// 創建客戶端對象
func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

// 實現通用的RPC客戶端
// 綁定RPC訪問的方法
// 傳入訪問的函數名

// 函數具體實現在Server端，Client只有函數原型
// 使用MakeFunc()完成原型到函數的調用

// fPtr指向函數原型
// xxx.callRPC("queryUser",&query)
func (c *Client) callRPC(rpcName string, fPtr interface{}) {
	// 通過反射，獲取fPtr未初始化的函數原型
	fn := reflect.ValueOf(fPtr).Elem()

	// 另一個函數，作用是對第一個函數參數操作
	// 完成與Server的交互
	f := func(args []reflect.Value) []reflect.Value {
		// 處理輸入的參數
		inArgs := make([]interface{}, 0, len(args))
		for _, arg := range args {
			inArgs = append(inArgs, arg.Interface())
		}

		// 創建連接
		cliSession := NewSession(c.conn)

		// 編碼數據
		reqRPC := RPCData{Name: rpcName, Args: inArgs}
		b, err := encode(reqRPC)
		if err != nil {
			panic(err)
		}

		// 寫出數據
		err = cliSession.Write(b)
		if err != nil {
			panic(err)
		}

		// 讀取響應數據
		respBytes, err := cliSession.Read()
		if err != nil {
			panic(err)
		}

		// 解碼數據
		respRPC, err := decode(respBytes)
		if err != nil {
			panic(err)
		}

		// 處理服務端返回的數據
		outArgs := make([]reflect.Value, 0, len(respRPC.Args))
		for i, arg := range respRPC.Args {
			// 必須進行nil轉換
			if arg == nil {
				// 必須填充一個真正的類型，不能是nil
				outArgs = append(outArgs, reflect.Zero(fn.Type().Out(i)))
				continue
			}
			outArgs = append(outArgs, reflect.ValueOf(arg))
		}
		return outArgs
	}
	// 參數1: 一個未初始化函數的方法值，類型是reflect.Type
	// 參數2: 另一個函數，作用是對第一個函數參數操作
	// 返回 reflect.Value類型
	// MakeFunc 使用傳入的函數原型，創建一個綁定 參數2 的新函數
	v := reflect.MakeFunc(fn.Type(), f)
	// 為函數fPtr賦值
	fn.Set(v)
}
