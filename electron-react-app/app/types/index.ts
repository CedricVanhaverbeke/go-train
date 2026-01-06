export interface WorkoutStep {
  duration: number
  start_power: number
  end_power: number
  type?: 'steady' | 'ramp'
}

export interface Workout {
  name: string
  url?: string
  steps: WorkoutStep[]
  totalDuration?: number
  baseTotalDuration?: number
  targetDuration?: number
  durationScale?: number
}

export interface TrainingFile {
  id: number
  name: string
  created_at: string
}

export type SortOrder = 'default' | 'asc' | 'desc' | 'duration_asc' | 'duration_desc'

export type ViewType = 'list' | 'create' | 'settings' | 'training-list'

export interface WorkoutAPI {
  startApp: (appName: string, ...args: string[]) => void
  stopApp: (appName: string) => void
  getGpxFiles: () => void
  getGpxFileData: (id: number) => void
  onAppStatus: (callback: (event: any, message: string) => void) => void
  onAppStdout: (callback: (event: any, data: string) => void) => void
  onAppStderr: (callback: (event: any, data: string) => void) => void
  onAppExited: (callback: (event: any) => void) => void
  onGpxFilesData: (callback: (event: any, files: TrainingFile[]) => void) => void
  onGpxFilesError: (callback: (event: any, error: string) => void) => void
  onGpxFileData: (callback: (event: any, fileData: any) => void) => void
  onGpxFileError: (callback: (event: any, error: string) => void) => void
  removeAllListeners: (channel: string) => void
}

declare global {
  interface Window {
    workoutAPI: WorkoutAPI
  }
}
