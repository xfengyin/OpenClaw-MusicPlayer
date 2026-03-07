const { app, BrowserWindow, ipcMain } = require('electron')
const path = require('path')
const { spawn } = require('child_process')

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
    icon: path.join(__dirname, '../build/icon.png')
  })

  // 加载前端页面
  const isDev = !app.isPackaged
  if (isDev) {
    mainWindow.loadURL('http://localhost:5173')
    mainWindow.webContents.openDevTools()
  } else {
    mainWindow.loadFile(path.join(__dirname, '../dist/index.html'))
  }
  
  // 启动 Go 后端服务
  startServer()
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
