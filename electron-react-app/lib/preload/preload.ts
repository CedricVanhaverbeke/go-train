import { contextBridge, ipcRenderer } from 'electron'
import { conveyor } from '@/lib/conveyor/api'

// Custom API for workout app functionality
const workoutAPI = {
  // App control
  startApp: (appName: string, ...args: string[]) => ipcRenderer.send('START_APP', appName, ...args),
  stopApp: (appName: string) => ipcRenderer.send('STOP_APP', appName),
  
  // GPX file management
  getGpxFiles: () => ipcRenderer.send('GET_GPX_FILES'),
  getGpxFileData: (id: number) => ipcRenderer.send('GET_GPX_FILE_DATA', id),
  
  // Event listeners
  onAppStatus: (callback: (event: any, message: string) => void) => ipcRenderer.on('APP_STATUS', callback),
  onAppStdout: (callback: (event: any, data: string) => void) => ipcRenderer.on('APP_STDOUT', callback),
  onAppStderr: (callback: (event: any, data: string) => void) => ipcRenderer.on('APP_STDERR', callback),
  onAppExited: (callback: (event: any) => void) => ipcRenderer.on('APP_EXITED', callback),
  onGpxFilesData: (callback: (event: any, files: any[]) => void) => ipcRenderer.on('GPX_FILES_DATA', callback),
  onGpxFilesError: (callback: (event: any, error: string) => void) => ipcRenderer.on('GPX_FILES_ERROR', callback),
  onGpxFileData: (callback: (event: any, fileData: any) => void) => ipcRenderer.on('GPX_FILE_DATA', callback),
  onGpxFileError: (callback: (event: any, error: string) => void) => ipcRenderer.on('GPX_FILE_ERROR', callback),
  
  // Remove listeners
  removeAllListeners: (channel: string) => ipcRenderer.removeAllListeners(channel)
}

// Use `contextBridge` APIs to expose APIs to
// renderer only if context isolation is enabled, otherwise
// just add to the DOM global.
if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld('conveyor', conveyor)
    contextBridge.exposeInMainWorld('workoutAPI', workoutAPI)
  } catch (error) {
    console.error(error)
  }
} else {
  window.conveyor = conveyor
  window.workoutAPI = workoutAPI
}
