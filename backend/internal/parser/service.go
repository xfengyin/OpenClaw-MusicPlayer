package parser

import (
	"fmt"
	"sync"
	"time"
)

// ParserService 音乐解析服务
type ParserService struct {
	parsers map[string]MusicParser
	cache   *Cache
	mu      sync.RWMutex
}

// MusicParser 音乐解析器接口
type MusicParser interface {
	Search(keyword string, limit, offset int) (*SearchResult, error)
	GetSongURL(id string, quality int) (string, error)
	GetLyrics(id string) (string, string, error)
}

// Cache 缓存
type Cache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
	ttl  time.Duration
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewParserService 创建解析服务
func NewParserService() *ParserService {
	s := &ParserService{
		parsers: make(map[string]MusicParser),
		cache: &Cache{
			data: make(map[string]cacheItem),
			ttl:  30 * time.Minute,
		},
	}

	// 注册网易云解析器
	s.RegisterParser("netease", NewNetEaseParser())

	return s
}

// RegisterParser 注册解析器
func (s *ParserService) RegisterParser(name string, parser MusicParser) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.parsers[name] = parser
}

// GetParser 获取解析器
func (s *ParserService) GetParser(name string) (MusicParser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parser, ok := s.parsers[name]
	if !ok {
		return nil, fmt.Errorf("parser not found: %s", name)
	}

	return parser, nil
}

// Search 搜索音乐（聚合多个源）
func (s *ParserService) Search(keyword string, limit, offset int) (*SearchResult, error) {
	// 检查缓存
	cacheKey := fmt.Sprintf("search:%s:%d:%d", keyword, limit, offset)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*SearchResult), nil
	}

	// 使用网易云搜索
	parser, err := s.GetParser("netease")
	if err != nil {
		return nil, err
	}

	result, err := parser.Search(keyword, limit, offset)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	s.cache.Set(cacheKey, result)

	return result, nil
}

// GetSongURL 获取歌曲URL
func (s *ParserService) GetSongURL(source, id string, quality int) (string, error) {
	parser, err := s.GetParser(source)
	if err != nil {
		return "", err
	}

	return parser.GetSongURL(id, quality)
}

// GetLyrics 获取歌词
func (s *ParserService) GetLyrics(source, id string) (string, string, error) {
	// 检查缓存
	cacheKey := fmt.Sprintf("lyrics:%s:%s", source, id)
	if cached, ok := s.cache.Get(cacheKey); ok {
		item := cached.([2]string)
		return item[0], item[1], nil
	}

	parser, err := s.GetParser(source)
	if err != nil {
		return "", "", err
	}

	lyrics, translated, err := parser.GetLyrics(id)
	if err != nil {
		return "", "", err
	}

	// 缓存结果
	s.cache.Set(cacheKey, [2]string{lyrics, translated})

	return lyrics, translated, nil
}

// Cache methods
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]cacheItem)
}
