const { app, BrowserWindow, ipcMain } = require('electron')
const path = require('path')
const { spawn } = require('child_process')
const fs = require('fs')

let mainWindow
let serverProcess

function createWindow () {
  console.log('Creating window...')
  console.log('__dirname:', __dirname)
  console.log('app.getAppPath():', app.getAppPath())
  console.log('process.resourcesPath:', process.resourcesPath)
  
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js')
    },
    title: 'OpenClaw Music Player',
    show: true // 立即显示窗口
  })

  // 加载前端页面
  const isDev = !app.isPackaged
  console.log('isDev:', isDev)
  
  if (isDev) {
    mainWindow.loadURL('http://localhost:5173')
    mainWindow.webContents.openDevTools()
  } else {
    // 生产环境：查找 index.html
    const indexPath = findIndexHtml()
    console.log('Loading index.html from:', indexPath)
    
    if (indexPath && fs.existsSync(indexPath)) {
      console.log('File exists, loading...')
      mainWindow.loadFile(indexPath)
    } else {
      console.error('index.html not found at:', indexPath)
      // 显示错误页面
      mainWindow.loadURL(`data:text/html,<html><body style="padding:20px;font-family:sans-serif"><h1>Error: index.html not found</h1><p>Expected path: ${indexPath}</p><p>__dirname: ${__dirname}</p><p>appPath: ${app.getAppPath()}</p><p>resourcesPath: ${process.resourcesPath}</p></body></html>`)
    }
    
    // 打开开发者工具以便调试
    mainWindow.webContents.openDevTools()
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
    path.join(__dirname, '../dist/index.html'),
    path.join(__dirname, '../../dist/index.html'),
    path.join(process.resourcesPath, 'dist/index.html'),
    path.join(app.getAppPath(), 'dist/index.html')
  ]
  
  console.log('Searching for index.html in', possiblePaths.length, 'locations:')
  for (const p of possiblePaths) {
    const exists = fs.existsSync(p)
    console.log('  -', p, exists ? '(EXISTS)' : '(not found)')
    if (exists) {
      return p
    }
  }
  
  // 如果没找到，列出所有目录内容帮助调试
  console.log('\nDirectory contents for debugging:')
  try {
    console.log('__dirname:', __dirname)
    console.log('Contents:', fs.readdirSync(__dirname))
    
    const parentDir = path.join(__dirname, '..')
    console.log('\nParent dir:', parentDir)
    if (fs.existsSync(parentDir)) {
      console.log('Contents:', fs.readdirSync(parentDir))
    }
    
    const resourcesDir = process.resourcesPath
    console.log('\nResources dir:', resourcesDir)
    if (fs.existsSync(resourcesDir)) {
      console.log('Contents:', fs.readdirSync(resourcesDir))
    }
  } catch (e) {
    console.error('Error listing directories:', e)
  }
  
  // 返回第一个路径（让错误信息显示正确的路径）
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
  
  // 检查服务器文件是否存在
  if (!fs.existsSync(serverPath)) {
    console.error('Server binary not found at:', serverPath)
    return
  }
  
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
  
  serverProcess.on('error', (err) => {
    console.error('Failed to start server:', err)
  })
}

app.whenReady().then(() => {
  console.log('App is ready')
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
