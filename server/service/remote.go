package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RemoteConnection struct {
	ID         string    `json:"id"`
	IP         string    `json:"ip"`
	AuthKey    string    `json:"auth_key"`
	CreatedAt  time.Time `json:"created_at"`
	LastActive time.Time `json:"last_active"`
	WebSocket  *websocket.Conn
}

type RemoteControlService struct {
	connections   map[string]*RemoteConnection
	connectionsLock sync.RWMutex
	upgrader      websocket.Upgrader
	onControl     func(command string, data json.RawMessage) error
}

func NewRemoteService() *RemoteControlService {
	return &RemoteControlService{
		connections: make(map[string]*RemoteConnection),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// 允许跨域
				return true
			},
		},
	}
}

// SetControlHandler 设置控制回调
func (s *RemoteControlService) SetControlHandler(handler func(string, json.RawMessage) error) {
	s.onControl = handler
}

// HandleWS 处理 WebSocket 连接
func (s *RemoteControlService) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	ip := s.getClientIP(r)
	authKey := r.URL.Query().Get("auth")

	connection := &RemoteConnection{
		ID:        fmt.Sprintf("remote_%d", time.Now().UnixNano()),
		IP:        ip,
		AuthKey:   authKey,
		CreatedAt: time.Now(),
		LastActive: time.Now(),
		WebSocket: conn,
	}

	s.connectionsLock.Lock()
	s.connections[connection.ID] = connection
	s.connectionsLock.Unlock()

	go s.readLoop(connection)
}

// readLoop 读取循环
func (s *RemoteControlService) readLoop(conn *RemoteConnection) {
	defer func() {
		s.connectionsLock.Lock()
		delete(s.connections, conn.ID)
		s.connectionsLock.Unlock()
		conn.WebSocket.Close()
	}()

	for {
		_, message, err := conn.WebSocket.ReadMessage()
		if err != nil {
			return
		}

		conn.LastActive = time.Now()

		// 解析命令
		var cmd struct {
			Command string          `json:"command"`
			Data    json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(message, &cmd); err != nil {
			continue
		}

		// 调用处理函数
		if s.onControl != nil {
			if err := s.onControl(cmd.Command, cmd.Data); err != nil {
				// 发送错误响应
				conn.WebSocket.WriteJSON(map[string]string{
					"type": "error",
					"msg":  err.Error(),
				})
			}
		}
	}
}

// getClientIP 获取客户端 IP
func (s *RemoteControlService) getClientIP(r *http.Request) string {
	// 检查 X-Forwarded-For
	if xfwd := r.Header.Get("X-Forwarded-For"); xfwd != "" {
		ips := strings.Split(xfwd, ",")
		return strings.TrimSpace(ips[0])
	}

	// 检查 X-Real-IP
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	// 使用 RemoteAddr
	return r.RemoteAddr
}

// Broadcast 广播消息
func (s *RemoteControlService) Broadcast(message interface{}) error {
	_, err := json.Marshal(message)
	if err != nil {
		return err
	}

	s.connectionsLock.RLock()
	defer s.connectionsLock.RUnlock()

	for _, conn := range s.connections {
		if conn.WebSocket != nil {
			conn.WebSocket.WriteJSON(message)
		}
	}

	return nil
}

// GetConnections 获取连接列表
func (s *RemoteControlService) GetConnections() []*RemoteConnection {
	s.connectionsLock.RLock()
	defer s.connectionsLock.RUnlock()

	connections := make([]*RemoteConnection, 0, len(s.connections))
	for _, conn := range s.connections {
		connections = append(connections, conn)
	}

	return connections
}

// RemoteControl 处理远程控制
func (s *RemoteControlService) RemoteControl(command string, data json.RawMessage) error {
	// 广播命令
	return s.Broadcast(map[string]interface{}{
		"type":    "command",
		"command": command,
		"data":    data,
	})
}

// StartHTTPServer 启动 HTTP 服务
func (s *RemoteControlService) StartHTTPServer(port string) error {
	http.HandleFunc("/remote/ws", s.HandleWS)
	http.HandleFunc("/remote/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	address := fmt.Sprintf(":%s", port)
	return http.ListenAndServe(address, nil)
}