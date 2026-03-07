package service

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
)

type LyricLine struct {
	Timestamp string `json:"timestamp"`
	Seconds   float64 `json:"seconds"`
	Text      string `json:"text"`
}

type Lyrics struct {
	SongID     string       `json:"song_id"`
	Artist     string       `json:"artist"`
	Title      string       `json:"title"`
	Lines      []LyricLine  `json:"lines"`
	Translation []LyricLine `json:"translation"`
}

type LyricsService struct {
	krcRegex *regexp.Regexp
	lrcRegex *regexp.Regexp
}

func NewLyricsService() *LyricsService {
	return &LyricsService{
		krcRegex: regexp.MustCompile(`\[\d+\.\d+\]`),
		lrcRegex: regexp.MustCompile(`\[(\d+):(\d+\.\d+)\](.*)`),
	}
}

// ParseLyrics 解析歌词
func (s *LyricsService) ParseLyrics(krcData, lrcData string) (*Lyrics, error) {
	// 解析 LRC 歌词
	var lines []LyricLine
	scanner := bufio.NewScanner(strings.NewReader(lrcData))
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// 匹配歌词行
		matches := s.lrcRegex.FindStringSubmatch(line)
		if len(matches) == 4 {
			mins, _ := strconv.Atoi(matches[1])
			secs, _ := strconv.ParseFloat(matches[2])
			text := strings.TrimSpace(matches[3])
			
			seconds := float64(mins*60) + secs
			lines = append(lines, LyricLine{
				Timestamp: matches[0],
				Seconds:   seconds,
				Text:      text,
			})
		}
	}

	return &Lyrics{
		Lines: lines,
	}, nil
}

// ConvertToSimplified 转换为简体中文
func (s *LyricsService) ConvertToSimplified(text string) string {
	// Go 实现的繁简转换
	// 这里使用简单的映射，实际项目中应使用 opencc
	simplifiedMap := map[string]string{
		"開": "开",
		"發": "发",
		"後": "后",
		"體": "体",
		"裏": "里",
		"說": "说",
		"個": "个",
		"黨": "党",
		"學": "学",
		"體": "体",
	}

	result := text
	for traditional, simplified := range simplifiedMap {
		result = strings.ReplaceAll(result, traditional, simplified)
	}

	return result
}

// SearchLyrics 搜索歌词
func (s *LyricsService) SearchLyrics(songID string) (*Lyrics, error) {
	// 从网易云音乐 API 获取歌词
	return nil, nil
}

// GetLyricsForSong 获取歌曲歌词
func (s *LyricsService) GetLyricsForSong(songID string) (*Lyrics, error) {
	// 获取歌词
	lyrics, err := s.SearchLyrics(songID)
	if err != nil {
		return nil, err
	}

	// 转换为简体中文
	for i, line := range lyrics.Lines {
		lyrics.Lines[i].Text = s.ConvertToSimplified(line.Text)
	}

	return lyrics, nil
}