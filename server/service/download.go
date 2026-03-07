package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DownloadTask struct {
	ID         string    `json:"id"`
	ProcessID  string    `json:"process_id"`
	Title      string    `json:"title"`
	Size       int64     `json:"size"`
	Downloaded int64     `json:"downloaded"`
	Speed      float64   `json:"speed"`
	Status     string    `json:"status"` // queuing, downloading, completed, failed
	Progress   float64   `json:"progress"`
	CreateTime time.Time `json:"create_time"`
	FilePath   string    `json:"file_path"`
}

type DownloadService struct {
	tasks         map[string]*DownloadTask
	tasksLock     sync.RWMutex
	maxConcurrent int
}

func NewDownloadService() *DownloadService {
	return &DownloadService{
		tasks:         make(map[string]*DownloadTask),
		maxConcurrent: 5,
	}
}

// StartDownload 开始下载
func (s *DownloadService) StartDownload(musicID, title, url string) (*DownloadTask, error) {
	task := &DownloadTask{
		ID:          musicID,
		Title:       title,
		ProcessID:   fmt.Sprintf("%s_%d", musicID, time.Now().UnixNano()),
		CreateTime:  time.Now(),
		Status:      "queuing",
		FilePath:    filepath.Join("downloads", fmt.Sprintf("%s.mp3", musicID)),
	}

	// 检查是否已存在
	s.tasksLock.RLock()
	_, exists := s.tasks[task.ID]
	s.tasksLock.RUnlock()

	if exists {
		return task, fmt.Errorf("任务已存在")
	}

	// 添加到任务队列
	s.tasksLock.Lock()
	s.tasks[task.ID] = task
	s.tasksLock.Unlock()

	// 异步开始下载
	go s.download(task, url)

	return task, nil
}

// download 下载音乐
func (s *DownloadService) download(task *DownloadTask, url string) {
	// 等待队列
	s.tasksLock.Lock()
	task.Status = "downloading"
	s.tasksLock.Unlock()

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(task.FilePath), 0755); err != nil {
		task.Status = "failed"
		return
	}

	//发送 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		task.Status = "failed"
		return
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", task.Downloaded))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		task.Status = "failed"
		return
	}
	defer resp.Body.Close()

	//追加到文件
	file, err := os.OpenFile(task.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		task.Status = "failed"
		return
	}
	defer file.Close()

	startTime := time.Now()
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := file.Write(buf[:n]); writeErr != nil {
				task.Status = "failed"
				return
			}

			task.Downloaded += int64(n)
			task.Size = resp.ContentLength

			// 计算速度
			elapsed := time.Since(startTime).Seconds()
			if elapsed > 0 {
				task.Speed = float64(task.Downloaded) / elapsed / 1024 // KB/s
			}

			// 计算进度
			if task.Size > 0 {
				task.Progress = float64(task.Downloaded) / float64(task.Size) * 100
			}

			// 检查是否完成
			if task.Size > 0 && task.Downloaded >= task.Size {
				task.Status = "completed"
				return
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				task.Status = "completed"
				return
			}
			task.Status = "failed"
			return
		}
	}
}

// GetDownloadStatus 获取下载状态
func (s *DownloadService) GetDownloadStatus(id string) (*DownloadTask, error) {
	s.tasksLock.RLock()
	defer s.tasksLock.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("任务不存在")
	}

	return task, nil
}

// ListDownloads 列出所有下载
func (s *DownloadService) ListDownloads() []*DownloadTask {
	s.tasksLock.RLock()
	defer s.tasksLock.RUnlock()

	tasks := make([]*DownloadTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// CancelDownload 取消下载
func (s *DownloadService) CancelDownload(id string) error {
	s.tasksLock.Lock()
	defer s.tasksLock.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("任务不存在")
	}

	// 移除任务
	delete(s.tasks, id)

	// 删除未完成的文件
	if task.Status != "completed" && task.FilePath != "" {
		os.Remove(task.FilePath)
	}

	return nil
}