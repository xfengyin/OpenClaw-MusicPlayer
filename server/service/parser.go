package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type MusicInfo struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Artist      string   `json:"artist"`
	Album       string   `json:"album"`
	URL         string   `json:"url"`
	Cover       string   `json:"cover"`
	Duration    int      `json:"duration"`
	Bitrate     int      `json:"bitrate"`
	Sources     []string `json:"sources"`
	Source      string   `json:"source"`
	Hash        string   `json:"hash"`
}

type SearchResult struct {
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
	Items     []MusicInfo `json:"items"`
}

type ParserService struct {
	cache      map[string]interface{}
	cacheLock  sync.RWMutex
	httpClient *http.Client
}

func NewParserService() *ParserService {
	return &ParserService{
		cache:      make(map[string]interface{}),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SearchMusic 搜索音乐
func (s *ParserService) SearchMusic(keyword string, page int, size int) (*SearchResult, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("search:%s:%d:%d", keyword, page, size)
	if info, ok := s.getFromCache(cacheKey); ok {
		return info, nil
	}

	// 调用网易云 API
	result := &SearchResult{}
	
	// 构建请求参数
	params := map[string]string{
		"keywords": keyword,
		"limit":    fmt.Sprintf("%d", size),
		"offset":   fmt.Sprintf("%d", (page-1)*size),
	}

	// 并发请求多个音源
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	// 请求网易云音乐
	wg.Add(1)
	go func() {
		defer wg.Done()
		neteaseResult, err := s.requestNetease(params)
		if err == nil && neteaseResult != nil {
			mu.Lock()
			result.Items = append(result.Items, neteaseResult...)
			result.Total = len(result.Items)
			mu.Unlock()
		}
	}()
	
	// 等待所有请求完成
	wg.Wait()

	// 添加到缓存
	s.addToCache(cacheKey, &SearchResult{
		Total:     result.Total,
		Page:      page,
		PageSize:  size,
		Items:     result.Items,
	})

	return result, nil
}

// GetMusicDetail 获取音乐详情
func (s *ParserService) GetMusicDetail(id string) (*MusicInfo, error) {
	// 请求详情
	detail, err := s.requestNeteaseDetail(id)
	if err != nil {
		return nil, err
	}

	// 计算歌曲 hash
	hash := s.calculateHash(detail.ID + detail.Title)
	detail.Hash = hash

	// 添加到缓存
	s.addToCache(id, detail)

	return detail, nil
}

// GetMusicUrl 获取音乐播放地址
func (s *ParserService) GetMusicUrl(id string, source string) (*MusicInfo, error) {
	detail, err := s.GetMusicDetail(id)
	if err != nil {
		return nil, err
	}

	// 解密音乐地址
	url, err := s.decryptMusicUrl(id, source)
	if err != nil {
		return nil, err
	}

	detail.URL = url
	return detail, nil
}

// decryptMusicUrl 解密音乐播放地址
func (s *ParserService) decryptMusicUrl(id string, source string) (string, error) {
	// 调用网易云解密接口
	params := map[string]string{
		"id": id,
	}

	// 这里实际实现解密逻辑
	// 参考: https://github.com/anonymous5l/ncm
	_, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	resp, err := s.httpClient.Post(
		"https://interface3.music.163.com/eapi/song/enhance/player/url",
		"application/json",
		nil,
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Code int    `json:"code"`
		Data []struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Code != 200 || len(result.Data) == 0 {
		return "", fmt.Errorf("获取音乐地址失败")
	}

	return result.Data[0].URL, nil
}

// calculateHash 计算音乐 hash
func (s *ParserService) calculateHash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}

// requestNetease 请求网易云音乐 API
func (s *ParserService) requestNetease(params map[string]string) ([]MusicInfo, error) {
	// 构建请求 URL
	url := "https://music.163.com/api/cloudsearch/pc"
	
	// 添加公共参数
	params["csrf_token"] = ""
	
	// 构建表单数据
	formData := ""
	for k, v := range params {
		formData += fmt.Sprintf("%s=%s&", k, v)
	}
	formData = strings.TrimSuffix(formData, "&")

	// 发送 POST 请求
	req, err := http.NewRequest("POST", url, strings.NewReader(formData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://music.163.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var result struct {
		Code int `json:"code"`
		Result struct {
			Songs []struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Ar       []struct {
					Name string `json:"name"`
				} `json:"ar"`
				Al struct {
					Name string `json:"name"`
					Pic  string `json:"picUrl"`
				} `json:"al"`
				Dt int `json:"dt"`
			} `json:"songs"`
			Total int `json:"total"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API 返回错误: %d", result.Code)
	}

	// 转换为 MusicInfo
	var items []MusicInfo
	for _, song := range result.Result.Songs {
		artist := ""
		if len(song.Ar) > 0 {
			artist = song.Ar[0].Name
		}

		items = append(items, MusicInfo{
			ID:      fmt.Sprintf("%d", song.ID),
			Title:   song.Name,
			Artist:  artist,
			Album:   song.Al.Name,
			Cover:   song.Al.Pic,
			Duration: song.Dt / 1000,
			Sources: []string{"netease"},
			Source:  "netease",
		})
	}

	return items, nil
}

// requestNeteaseDetail 请求歌曲详情
func (s *ParserService) requestNeteaseDetail(id string) (*MusicInfo, error) {
	url := fmt.Sprintf("https://music.163.com/api/song/detail?ids=[%s]", id)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Referer", "https://music.163.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code int `json:"code"`
		Songs []struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Ar     []struct {
				Name string `json:"name"`
			} `json:"ar"`
			Al struct {
				Name string `json:"name"`
				Pic  string `json:"picUrl"`
			} `json:"al"`
			Dt int `json:"dt"`
		} `json:"songs"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 200 || len(result.Songs) == 0 {
		return nil, fmt.Errorf("获取详情失败")
	}

	song := result.Songs[0]
	artist := ""
	if len(song.Ar) > 0 {
		artist = song.Ar[0].Name
	}

	return &MusicInfo{
		ID:       fmt.Sprintf("%d", song.ID),
		Title:    song.Name,
		Artist:   artist,
		Album:    song.Al.Name,
		Cover:    song.Al.Pic,
		Duration: song.Dt / 1000,
		Sources:  []string{"netease"},
		Source:   "netease",
	}, nil
}

// ParsePlaylist 解析歌单
func (s *ParserService) ParsePlaylist(url string) (map[string]interface{}, error) {
	// 实现歌单解析逻辑
	return map[string]interface{}{
		"url": url,
		"songs": []MusicInfo{},
	}, nil
}

// ParseSong 解析单曲
func (s *ParserService) ParseSong(url string) (*MusicInfo, error) {
	// 实现单曲解析逻辑
	return &MusicInfo{
		ID:     "0",
		Title:  "Unknown",
		Artist: "Unknown",
		Source: "unknown",
	}, nil
}

// getFromCache 从缓存获取
func (s *ParserService) getFromCache(key string) (*SearchResult, bool) {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()

	val, ok := s.cache[key]
	if !ok {
		return nil, false
	}

	result, ok := val.(*SearchResult)
	if !ok {
		return nil, false
	}
	return result, true
}

// addToCache 添加到缓存
func (s *ParserService) addToCache(key string, value interface{}) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	s.cache[key] = value
}