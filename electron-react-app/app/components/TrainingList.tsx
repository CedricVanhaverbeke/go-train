import React, { useState, useEffect } from 'react'
import { TrainingFile } from '../types'
import dayjs from 'dayjs'
import customParseFormat from 'dayjs/plugin/customParseFormat'
dayjs.extend(customParseFormat)

interface TrainingListProps {
  onBack: () => void
}

const TrainingList: React.FC<TrainingListProps> = ({ onBack: _onBack }) => {
  // onBack is kept for interface consistency though not currently used
  const [trainings, setTrainings] = useState<TrainingFile[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadTrainings()
  }, [])

  const loadTrainings = () => {
    setLoading(true)
    setError(null)

    // Request GPX files from main process
    window.workoutAPI.getGpxFiles()

    // Set up listeners for the response
    const handleFilesData = (_event: any, trainingsData: TrainingFile[]) => {
      setTrainings(trainingsData || [])
      setLoading(false)
      window.workoutAPI.removeAllListeners('GPX_FILES_DATA')
      window.workoutAPI.removeAllListeners('GPX_FILES_ERROR')
    }

    const handleFilesError = (_event: any, errorMessage: string) => {
      setError(errorMessage)
      setLoading(false)
      window.workoutAPI.removeAllListeners('GPX_FILES_DATA')
      window.workoutAPI.removeAllListeners('GPX_FILES_ERROR')
    }

    window.workoutAPI.onGpxFilesData(handleFilesData)
    window.workoutAPI.onGpxFilesError(handleFilesError)
  }

  const downloadTraining = (id: number, name: string) => {
    // Request GPX file data from main process
    window.workoutAPI.getGpxFileData(id)

    // Set up listeners for the response
    const handleFileData = (_event: any, fileData: any) => {
      if (fileData && fileData.data) {
        // Create blob and download link
        const blob = new Blob([fileData.data], { type: 'application/gpx+xml' })
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `${name || 'training'}.gpx`
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
        URL.revokeObjectURL(url)
      }
      window.workoutAPI.removeAllListeners('GPX_FILE_DATA')
      window.workoutAPI.removeAllListeners('GPX_FILE_ERROR')
    }

    const handleFileError = (_event: any, errorMessage: string) => {
      console.error('Error downloading file:', errorMessage)
      setError(`Failed to download ${name}: ${errorMessage}`)
      window.workoutAPI.removeAllListeners('GPX_FILE_DATA')
      window.workoutAPI.removeAllListeners('GPX_FILE_ERROR')
    }

    window.workoutAPI.onGpxFileData(handleFileData)
    window.workoutAPI.onGpxFileError(handleFileError)
  }

  const formatDate = (dateString: string): string => {
    const date = dayjs(dateString, 'YYYY-MM-DD HH:mm:ss.SSSSSS ZZ')
    return date.format('DD MMMM, YYYY')
  }

  return (
    <div className="space-y-4 max-h-[calc(100vh-12rem)] overflow-y-auto">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold text-slate-300">Training Files</h3>
        <button onClick={loadTrainings} className="text-sm text-slate-400 hover:text-slate-200">
          Refresh
        </button>
      </div>

      {error && (
        <div className="bg-red-500/20 border border-red-500/50 text-red-300 px-4 py-3 rounded-lg">
          <p className="text-sm">{error}</p>
        </div>
      )}

      {loading && (
        <div className="text-center py-8">
          <p className="text-slate-400">Loading training files...</p>
        </div>
      )}

      {!loading && !error && trainings.length === 0 && (
        <div className="text-center py-8">
          <p className="text-slate-400">No training files found</p>
        </div>
      )}

      {!loading && !error && trainings.length > 0 && (
        <div className="space-y-2">
          {trainings.map((training) => (
            <div key={training.id} className="bg-slate-700 rounded-lg p-4 flex items-center justify-between">
              <div className="flex-1">
                <h4 className="font-medium text-slate-100">{training.name || 'Untitled Training'}</h4>
                <p className="text-sm text-slate-400">{formatDate(training.created_at)}</p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs text-slate-500">ID: {training.id}</span>
                <button
                  onClick={() => downloadTraining(training.id, training.name)}
                  className="inline-flex items-center gap-2 rounded-lg bg-sky-600 hover:bg-sky-500 transition-colors px-4 py-2 font-semibold text-white text-sm"
                >
                  Download
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="text-sm text-slate-400 text-center">
        {trainings.length} training file{trainings.length !== 1 ? 's' : ''} found
      </div>
    </div>
  )
}

export default TrainingList
