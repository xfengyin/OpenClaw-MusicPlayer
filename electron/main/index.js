const { app, BrowserWindow, ipcMain, shell } = require('electron')
const path = require('path')
const { spawn } = require('child_process')

let mainWindow
let serverProcess

function createWindow () {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
     /nodeIntegration: false,
      contextIsolation: true,
      enableRemoteModule: true
    },
    title: 'OpenClaw Music Player',
    icon: path.join(__dirname, '../public/logo.png')
  })

  mainWindow.loadFile(path.join(__dirname, '../dist/index.html'))
  
  // 初始化 Go 服务
  startServer()
}

function startServer() {
  const serverPath = path.join(__dirname, '../../../server/server')
  
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
  if (process.platform !== 'darwin') {
    if (serverProcess) serverProcess.kill()
    app.quit()
  }
})

app.on('before-quit', () => {
  if (serverProcess) serverProcess.kill()
})
