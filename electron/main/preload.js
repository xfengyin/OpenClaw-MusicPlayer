const { contextBridge, ipcRenderer } = require('electron')

// 安全地暴露 API 给渲染进程
contextBridge.exposeInMainWorld('electronAPI', {
  // 可以在这里添加 IPC 通信方法
  platform: process.platform
})
