package parser

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// NetEaseParser 网易云音乐解析器
type NetEaseParser struct {
	client    *http.Client
	baseURL   string
	csrfToken string
}

// NewNetEaseParser 创建网易云解析器
func NewNetEaseParser() *NetEaseParser {
	return &NetEaseParser{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://music.163.com",
	}
}

// SongInfo 歌曲信息
type SongInfo struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Artists  []Artist `json:"artists"`
	Album    Album    `json:"album"`
	Duration int      `json:"duration"`
	URL      string   `json:"url"`
}

// Artist 艺术家信息
type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Album 专辑信息
type Album struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	PicURL string `json:"picUrl"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Songs []SongInfo `json:"songs"`
	Total int        `json:"total"`
}

// Search 搜索音乐
func (p *NetEaseParser) Search(keyword string, limit, offset int) (*SearchResult, error) {
	apiURL := fmt.Sprintf("%s/weapi/cloudsearch/get/web", p.baseURL)

	params := map[string]interface{}{
		"s":      keyword,
		"type":   1,
		"limit":  limit,
		"offset": offset,
	}

	data, err := p.encryptParams(params)
	if err != nil {
		return nil, err
	}

	resp, err := p.post(apiURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code   int `json:"code"`
		Result struct {
			Songs []struct {
				ID       int64  `json:"id"`
				Name     string `json:"name"`
				Ar       []struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
				} `json:"ar"`
				Al struct {
					ID     int64  `json:"id"`
					Name   string `json:"name"`
					PicURL string `json:"picUrl"`
				} `json:"al"`
				Dt int `json:"dt"`
			} `json:"songs"`
			SongCount int `json:"songCount"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API error: code %d", result.Code)
	}

	songs := make([]SongInfo, 0, len(result.Result.Songs))
	for _, s := range result.Result.Songs {
		artists := make([]Artist, 0, len(s.Ar))
		for _, a := range s.Ar {
			artists = append(artists, Artist{
				ID:   strconv.FormatInt(a.ID, 10),
				Name: a.Name,
			})
		}

		songs = append(songs, SongInfo{
			ID:   strconv.FormatInt(s.ID, 10),
			Name: s.Name,
			Artists: artists,
			Album: Album{
				ID:     strconv.FormatInt(s.Al.ID, 10),
				Name:   s.Al.Name,
				PicURL: s.Al.PicURL,
			},
			Duration: s.Dt / 1000,
		})
	}

	return &SearchResult{
		Songs: songs,
		Total: result.Result.SongCount,
	}, nil
}

// GetSongURL 获取歌曲播放地址
func (p *NetEaseParser) GetSongURL(id string, quality int) (string, error) {
	apiURL := fmt.Sprintf("%s/weapi/song/enhance/player/url", p.baseURL)

	ids := []string{id}
	params := map[string]interface{}{
		"ids": ids,
		"br":  quality,
	}

	data, err := p.encryptParams(params)
	if err != nil {
		return "", err
	}

	resp, err := p.post(apiURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Code int `json:"code"`
		Data []struct {
			ID  int64  `json:"id"`
			URL string `json:"url"`
			BR  int    `json:"br"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Code != 200 || len(result.Data) == 0 {
		return "", fmt.Errorf("failed to get song URL")
	}

	return result.Data[0].URL, nil
}

// GetLyrics 获取歌词
func (p *NetEaseParser) GetLyrics(id string) (string, string, error) {
	apiURL := fmt.Sprintf("%s/weapi/song/lyric", p.baseURL)

	params := map[string]interface{}{
		"id": id,
		"lv": -1,
		"tv": -1,
	}

	data, err := p.encryptParams(params)
	if err != nil {
		return "", "", err
	}

	resp, err := p.post(apiURL, data)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var result struct {
		Code int `json:"code"`
		Lrc  struct {
			Lyric string `json:"lyric"`
		} `json:"lrc"`
		Tlyric struct {
			Lyric string `json:"lyric"`
		} `json:"tlyric"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", err
	}

	return result.Lrc.Lyric, result.Tlyric.Lyric, nil
}

// 加密参数
func (p *NetEaseParser) encryptParams(params map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	// 简化的加密实现（实际应该使用网易云的真实加密算法）
	secretKey := "7246674226682325323F5E6544673A51"
	encrypted := p.aesEncrypt(string(jsonData), secretKey)

	data := url.Values{}
	data.Set("params", encrypted)
	data.Set("encSecKey", p.rsaEncrypt(secretKey))

	return data.Encode(), nil
}

// AES加密（简化实现）
func (p *NetEaseParser) aesEncrypt(text, key string) string {
	// 实际应该实现AES-128-CBC加密
	// 这里使用MD5作为占位符
	h := md5.New()
	h.Write([]byte(text + key))
	return hex.EncodeToString(h.Sum(nil))
}

// RSA加密（简化实现）
func (p *NetEaseParser) rsaEncrypt(text string) string {
	// 实际应该实现RSA加密
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// POST请求
func (p *NetEaseParser) post(url, data string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://music.163.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.0")

	return p.client.Do(req)
}
