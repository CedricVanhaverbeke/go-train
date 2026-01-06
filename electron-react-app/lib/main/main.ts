import { app, BrowserWindow, ipcMain } from 'electron'
import { electronApp, optimizer } from '@electron-toolkit/utils'
import { createAppWindow } from './app'
import path from 'path'
import fs from 'fs'
import os from 'os'
import { spawn } from 'child_process'
import sqlite3 from 'sqlite3'

// Keep a global reference to the child process.
let child = null

// Define the commands that the renderer is allowed to start.
const AVAILABLE_APPS = {
  overlay: () => {
    const basePath = app.isPackaged
      ? path.join(process.resourcesPath, 'bin')
      : path.join(__dirname, '..', '..', '..', 'bin')

    // Path to a compiled Go server binary (replace with your real path).
    const binary = path.join(basePath, `overlay-${os.platform()}-${os.arch()}`)

    if (!fs.existsSync(binary)) {
      throw new Error(`Go server binary not found at ${binary}`)
    }

    return {
      command: binary,
      args: [],
      options: { cwd: path.dirname(binary) },
    }
  },
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  // Set app user model id for windows
  electronApp.setAppUserModelId('com.electron')
  // Create app window
  createAppWindow()

  // Default open or close DevTools by F12 in development
  // and ignore CommandOrControl + R in production.
  // see https://github.com/alex8088/electron-toolkit/tree/master/packages/utils
  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  app.on('activate', function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) {
      createAppWindow()
    }
  })
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// IPC handlers
function startApp(event, appName, ...args) {
  try {
    const builder = AVAILABLE_APPS[appName]

    if (!builder) {
      throw new Error(`The selected application is not registered in main.ts: ${appName}`)
    }

    const { command, args: builderArgs, options = {} } = builder()

    child = spawn(command, [...args, ...builderArgs], {
      ...options,
    })

    child.on('spawn', () => {
      event.reply('APP_STATUS', `Successfully started ${appName} (PID: ${child.pid}).`)
    })

    child.stdout.on('data', (data) => {
      console.log(data.toString())
      event.reply('APP_STDOUT', data.toString())
    })

    child.stderr.on('data', (data) => {
      console.log(data.toString())
      event.reply('APP_STDERR', data.toString())
    })

    child.on('error', (error) => {
      event.reply('APP_STATUS', `Failed to start ${appName}: ${error.message}`)
    })

    child.on('exit', (code) => {
      event.reply('APP_STATUS', `App exited with code: ${code}.`)
      event.reply('APP_EXITED')
      child = null
    })
  } catch (error) {
    event.reply('APP_STATUS', `Unable to start ${appName}: ${error.message}`)
  }
}

function stopApp(event, appName) {
  if (!child) {
    event.reply('APP_STATUS', 'No application is currently running.')
    return
  }

  child.kill()
  child = null

  event.reply('APP_STATUS', `Successfully stopped ${appName}.`)
}

function getGpxFiles(event) {
  const dbPath = path.join(__dirname, '..', '..', '..', 'db')
  const db = new sqlite3.Database(dbPath)

  db.all('SELECT id, name, created_at FROM gpx_files ORDER BY created_at DESC', [], (err, rows) => {
    if (err) {
      event.reply('GPX_FILES_ERROR', err.message)
    } else {
      event.reply('GPX_FILES_DATA', rows)
    }
    db.close()
  })
}

function getGpxFileData(event, id) {
  const dbPath = path.join(__dirname, '..', '..', '..', 'db')
  const db = new sqlite3.Database(dbPath)

  db.get('SELECT name, data FROM gpx_files WHERE id = ?', [id], (err, row) => {
    if (err) {
      event.reply('GPX_FILE_ERROR', err.message)
    } else {
      event.reply('GPX_FILE_DATA', row)
    }
    db.close()
  })
}

// Set up IPC handlers
app.whenReady().then(() => {
  ipcMain.on('START_APP', (event, ...appName) => {
    startApp(event, ...appName)
  })

  ipcMain.on('STOP_APP', (event, appName) => {
    stopApp(event, appName)
  })

  ipcMain.on('GET_GPX_FILES', (event) => {
    getGpxFiles(event)
  })

  ipcMain.on('GET_GPX_FILE_DATA', (event, id) => {
    getGpxFileData(event, id)
  })
})

app.on('before-quit', () => {
  if (child) {
    child.kill()
    child = null
  }
})
