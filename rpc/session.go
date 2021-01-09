package rpc

import (
	"encoding/binary"
	"io"
	"net"
)

// 編寫會話中數據讀寫

// 會話連接的結構體
type Session struct {
	conn net.Conn
}

// 創建新連接
func NewSession(conn net.Conn) *Session {
	return &Session{conn: conn}
}

// 向連接中寫數據
func (s *Session) Write(data []byte) error {
	// 4字節頭 + 數據長度切片
	buf := make([]byte, 4+len(data))

	// 寫入頭部數據，記錄數據長度
	// binary只認固定長度的類似，所以使用uint32，而不是直接寫入
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))

	// 寫入數據
	copy(buf[4:], data)
	_, err := s.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// 向連接中讀數據
func (s *Session) Read() ([]byte, error) {
	// 讀取頭部長度
	header := make([]byte, 4)

	// 按頭部長度，讀取頭部數據
	_, err := io.ReadFull(s.conn, header)
	if err != nil {
		return nil, err
	}

	// 讀取數據長度
	dataLen := binary.BigEndian.Uint32(header)
	// 按照數據長度去讀取數據
	data := make([]byte, dataLen)
	_, err = io.ReadFull(s.conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
