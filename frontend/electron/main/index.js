const { app, BrowserWindow, ipcMain, Tray, Menu, nativeImage } = require('electron')
const path = require('path')
const { spawn } = require('child_process')
const fs = require('fs')

let mainWindow
let tray
let serverProcess

// 是否为开发模式
const isDev = process.argv.includes('--dev')

// 创建主窗口
function createWindow() {
  console.log('Creating main window...')
  console.log('isDev:', isDev)
  console.log('__dirname:', __dirname)
  console.log('app.getAppPath():', app.getAppPath())

  mainWindow = new BrowserWindow({
    width: 1400,
    height: 900,
    minWidth: 1000,
    minHeight: 600,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, '../preload/preload.js')
    },
    title: 'OpenClaw Music Player',
    show: false,
    icon: path.join(__dirname, '../build/icon.png')
  })

  // 加载页面
  if (isDev) {
    // 开发模式：加载 Vite 开发服务器
    mainWindow.loadURL('http://localhost:5173')
    mainWindow.webContents.openDevTools()
  } else {
    // 生产模式：加载打包后的文件
    const indexPath = findIndexHtml()
    console.log('Loading index.html from:', indexPath)
    
    if (indexPath && fs.existsSync(indexPath)) {
      mainWindow.loadFile(indexPath)
    } else {
      console.error('index.html not found!')
      mainWindow.loadURL(`data:text/html,<html><body style="padding:20px;font-family:sans-serif"><h1>Error: index.html not found</h1><p>Path: ${indexPath}</p></body></html>`)
    }
  }

  // 窗口准备好后显示
  mainWindow.once('ready-to-show', () => {
    mainWindow.show()
    
    if (isDev) {
      mainWindow.webContents.openDevTools()
    }
  })

  // 窗口关闭时隐藏到托盘
  mainWindow.on('close', (event) => {
    if (!app.isQuiting) {
      event.preventDefault()
      mainWindow.hide()
    }
  })

  // 启动后端服务
  startServer()
}

// 查找 index.html
function findIndexHtml() {
  const possiblePaths = [
    path.join(__dirname, '../renderer/index.html'),
    path.join(__dirname, '../../web/dist/index.html'),
    path.join(process.resourcesPath, 'renderer/index.html'),
    path.join(app.getAppPath(), 'renderer/index.html'),
    path.join(__dirname, 'renderer/index.html')
  ]

  console.log('Searching for index.html...')
  for (const p of possiblePaths) {
    const exists = fs.existsSync(p)
    console.log('  -', p, exists ? '(EXISTS)' : '(not found)')
    if (exists) {
      return p
    }
  }

  return possiblePaths[0]
}

// 启动后端服务
function startServer() {
  if (isDev) {
    console.log('Dev mode: skipping server start')
    return
  }

  const platform = process.platform
  const exeName = platform === 'win32' ? 'server.exe' : 'server'
  const serverPath = path.join(process.resourcesPath, 'server', exeName)

  console.log('Starting server from:', serverPath)

  if (!fs.existsSync(serverPath)) {
    console.error('Server binary not found:', serverPath)
    return
  }

  serverProcess = spawn(serverPath, [], {
    cwd: path.dirname(serverPath),
    env: { ...process.env, SERVER_PORT: '8080' }
  })

  serverProcess.stdout.on('data', (data) => {
    console.log(`[Server] ${data}`)
  })

  serverProcess.stderr.on('data', (data) => {
    console.error(`[Server Error] ${data}`)
  })

  serverProcess.on('close', (code) => {
    console.log(`Server exited with code ${code}`)
  })

  serverProcess.on('error', (err) => {
    console.error('Failed to start server:', err)
  })
}

// 创建系统托盘
function createTray() {
  const iconPath = path.join(__dirname, '../build/tray-icon.png')
  const trayIcon = nativeImage.createFromPath(iconPath)
  
  tray = new Tray(trayIcon.resize({ width: 16, height: 16 }))
  
  const contextMenu = Menu.buildFromTemplate([
    {
      label: '显示主窗口',
      click: () => {
        mainWindow.show()
      }
    },
    {
      label: '播放/暂停',
      click: () => {
        mainWindow.webContents.send('player-toggle')
      }
    },
    { type: 'separator' },
    {
      label: '退出',
      click: () => {
        app.isQuiting = true
        app.quit()
      }
    }
  ])

  tray.setToolTip('OpenClaw Music Player')
  tray.setContextMenu(contextMenu)
  
  tray.on('click', () => {
    mainWindow.show()
  })
}

// 应用就绪
app.whenReady().then(() => {
  console.log('App is ready')
  createWindow()
  createTray()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    } else {
      mainWindow.show()
    }
  })
})

// 所有窗口关闭
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// 应用退出前
app.on('before-quit', () => {
  app.isQuiting = true
  if (serverProcess) {
    serverProcess.kill()
  }
})

// IPC 通信
ipcMain.handle('get-app-version', () => {
  return app.getVersion()
})

ipcMain.handle('minimize-window', () => {
  mainWindow.minimize()
})

ipcMain.handle('maximize-window', () => {
  if (mainWindow.isMaximized()) {
    mainWindow.unmaximize()
  } else {
    mainWindow.maximize()
  }
})

ipcMain.handle('close-window', () => {
  mainWindow.hide()
})
