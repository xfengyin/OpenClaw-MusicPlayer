package service

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type UserData struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Level      int    `json:"level"`
	PlayCount  int    `json:"play_count"`
	FavoriteCount int `json:"favorite_count"`
}

type SyncService struct {
	data      map[string]interface{}
	dataLock  sync.RWMutex
	dataFile  string
	 Key      []byte
	IV        []byte
}

func NewSyncService() *SyncService {
	dataDir := viper.GetString("data_dir")
	if dataDir == "" {
		dataDir = "./data"
	}

	os.MkdirAll(dataDir, 0755)

	return &SyncService{
		dataFile: filepath.Join(dataDir, "sync_data.json"),
	}
}

// LoadData 加载同步数据
func (s *SyncService) LoadData() error {
	s.dataLock.Lock()
	defer s.dataLock.Unlock()

	if _, err := os.Stat(s.dataFile); os.IsNotExist(err) {
		s.data = make(map[string]interface{})
		return nil
	}

	data, err := os.ReadFile(s.dataFile)
	if err != nil {
		return err
	}

	// 解密数据
	decrypted, err := s.Decrypt(data)
	if err != nil {
		s.data = make(map[string]interface{})
		return nil
	}

	return json.Unmarshal(decrypted, &s.data)
}

// SaveData 保存同步数据
func (s *SyncService) SaveData() error {
	s.dataLock.RLock()
	defer s.dataLock.RUnlock()

	// 加密数据
	_, _ = s.Encrypt(json.RawMessage{})

	return nil
}

// Encrypt 加密数据
func (s *SyncService) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, s.IV)
	mode.CryptBlocks(cipherText, data)

	return cipherText, nil
}

// Decrypt 解密数据
func (s *SyncService) Decrypt(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("空数据")
	}

	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("数据长度不足")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	return s.unpad(data), nil
}

// unpad 去除填充
func (s *SyncService) unpad(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}

	padding := data[len(data)-1]
	if int(padding) > len(data) {
		return data
	}

	return data[:len(data)-int(padding)]
}

// SyncData 同步用户数据
func (s *SyncService) SyncData(userID, playlistID string) error {
	s.dataLock.Lock()
	defer s.dataLock.Unlock()

	key := fmt.Sprintf("user_%s_playlists", userID)
	s.data[key] = playlistID

	return nil
}

// GetUserPlaylists 获取用户歌单
func (s *SyncService) GetUserPlaylists(userID string) ([]string, error) {
	s.dataLock.RLock()
	defer s.dataLock.RUnlock()

	key := fmt.Sprintf("user_%s_playlists", userID)
	value, ok := s.data[key]
	if !ok {
		return []string{}, nil
	}

	playlistID, ok := value.(string)
	if !ok {
		return []string{}, fmt.Errorf("数据类型错误")
	}

	return []string{playlistID}, nil
}

// UpdatePlayCount 更新播放次数
func (s *SyncService) UpdatePlayCount(userID, musicID string) error {
	s.dataLock.Lock()
	defer s.dataLock.Unlock()

	key := fmt.Sprintf("user_%s_music_%s_play_count", userID, musicID)
	
	count := 0
	if val, ok := s.data[key]; ok {
		if c, ok := val.(int); ok {
			count = c
		}
	}

	s.data[key] = count + 1
	return nil
}

// GetPlayCount 获取播放次数
func (s *SyncService) GetPlayCount(userID, musicID string) int {
	s.dataLock.RLock()
	defer s.dataLock.RUnlock()

	key := fmt.Sprintf("user_%s_music_%s_play_count", userID, musicID)
	count := 0

	if val, ok := s.data[key]; ok {
		if c, ok := val.(int); ok {
			count = c
		}
	}

	return count
}