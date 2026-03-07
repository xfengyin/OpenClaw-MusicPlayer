const { app, BrowserWindow, ipcMain } = require('electron')
const path = require('path')
const { spawn } = require('child_process')
const fs = require('fs')

let mainWindow
let serverProcess

function createWindow () {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js')
    },
    title: 'OpenClaw Music Player',
    icon: path.join(__dirname, '../build/icon.png'),
    show: false // 先不显示窗口，等加载完成后再显示
  })

  // 加载前端页面
  const isDev = !app.isPackaged
  if (isDev) {
    mainWindow.loadURL('http://localhost:5173')
    mainWindow.webContents.openDevTools()
    mainWindow.show()
  } else {
    // 生产环境：查找 index.html
    const indexPath = findIndexHtml()
    console.log('Loading index.html from:', indexPath)
    
    if (indexPath && fs.existsSync(indexPath)) {
      mainWindow.loadFile(indexPath)
    } else {
      console.error('index.html not found!')
      // 显示错误页面
      mainWindow.loadURL(`data:text/html,<h1>Error: index.html not found</h1><p>Path: ${indexPath}</p>`)
    }
    
    // 等页面加载完成后再显示窗口
    mainWindow.once('ready-to-show', () => {
      mainWindow.show()
    })
  }
  
  // 启动 Go 后端服务
  startServer()
}

// 查找 index.html 文件
function findIndexHtml() {
  const possiblePaths = [
    path.join(__dirname, '../renderer/index.html'),
    path.join(__dirname, '../../renderer/index.html'),
    path.join(process.resourcesPath, 'renderer/index.html'),
    path.join(app.getAppPath(), 'renderer/index.html'),
    path.join(__dirname, 'renderer/index.html'),
    // 兼容旧路径
    path.join(__dirname, '../dist/index.html'),
    path.join(__dirname, '../../dist/index.html'),
    path.join(process.resourcesPath, 'dist/index.html')
  ]
  
  console.log('Searching for index.html in:')
  for (const p of possiblePaths) {
    console.log('  -', p, fs.existsSync(p) ? '(exists)' : '(not found)')
    if (fs.existsSync(p)) {
      return p
    }
  }
  
  // 如果没找到，返回第一个路径（让错误信息显示正确的路径）
  return possiblePaths[0]
}

function startServer() {
  const isDev = !app.isPackaged
  let serverPath
  
  if (isDev) {
    serverPath = path.join(__dirname, '../../../server/server')
  } else {
    // 生产环境根据平台选择正确的二进制文件
    const platform = process.platform
    const exeName = platform === 'win32' ? 'server.exe' : 'server'
    serverPath = path.join(process.resourcesPath, 'server', exeName)
  }
  
  console.log('Starting server from:', serverPath)
  
  serverProcess = spawn(serverPath, [], {
    cwd: path.dirname(serverPath),
    env: { ...process.env, SERVER_PORT: '8080' }
  })
  
  serverProcess.stdout.on('data', (data) => {
    console.log(`Server: ${data}`)
  })
  
  serverProcess.stderr.on('data', (data) => {
    console.error(`Server Error: ${data}`)
  })
  
  serverProcess.on('close', (code) => {
    console.log(`Server exited with code ${code}`)
  })
}

app.whenReady().then(() => {
  createWindow()
  
  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

app.on('window-all-closed', () => {
  if (serverProcess) serverProcess.kill()
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

app.on('before-quit', () => {
  if (serverProcess) serverProcess.kill()
})
