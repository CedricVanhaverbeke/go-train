import React, { useState, useCallback, useMemo } from 'react'
import { Workout, SortOrder, ViewType } from './types'
import { getTotalDuration, scaleWorkoutDuration } from './lib/utils'
import WorkoutList from './components/WorkoutList'
import WorkoutDetail from './components/WorkoutDetail'
import WorkoutCreator from './components/WorkoutCreator'
import Settings from './components/Settings'
import TrainingList from './components/TrainingList'
import DurationRangeSlider from './components/DurationRangeSlider'
import './styles/app.css'

interface AppProps {
  workouts: Workout[]
  durationMin: number
  durationMax: number
  filterMin: number
  filterMax: number
  setWorkouts: (workouts: Workout[]) => void
  setDurationMin: (min: number) => void
  setDurationMax: (max: number) => void
  setFilterMin: (min: number) => void
  setFilterMax: (max: number) => void
}

const App: React.FC<AppProps> = ({
  workouts,
  durationMin,
  durationMax,
  filterMin,
  filterMax,
  setWorkouts,
  setFilterMin,
  setFilterMax,
}) => {
  const [selectedWorkout, setSelectedWorkout] = useState<Workout | null>(null)
  const [selectedWorkoutTargetDuration, setSelectedWorkoutTargetDuration] = useState<number | null>(null)
  const [sortOrder, setSortOrder] = useState<SortOrder>('default')
  const [view, setView] = useState<ViewType>('list')
  const [ftp, setFtp] = useState<number>(Number(localStorage.getItem('userFtp') || 250))

  const handleFtpChange = useCallback((newFtp: number) => {
    setFtp(newFtp)
    localStorage.setItem('userFtp', newFtp.toString())
  }, [])

  const showCreateWorkoutView = useCallback(() => {
    setView('create')
  }, [])

  const showSettingsView = useCallback(() => {
    setView('settings')
  }, [])

  const showTrainingListView = useCallback(() => {
    setView('training-list')
  }, [])

  const showListView = useCallback(() => {
    setView('list')
  }, [])

  const handleSaveWorkout = useCallback(
    (newWorkout: Workout) => {
      const workoutWithDuration = {
        ...newWorkout,
        totalDuration: getTotalDuration(newWorkout),
      }

      setWorkouts([...workouts, workoutWithDuration])
      setView('list')
    },
    [workouts, setWorkouts]
  )

  const selectWorkout = useCallback(
    (workout: Workout, ftp: number) => {
      const totalDuration = workout.totalDuration ?? getTotalDuration(workout)
      const workoutString =
        workout.name +
        ';' +
        `${ftp}` +
        ';' +
        workout.steps
          .map(({ start_power, end_power, duration }) => {
            return `${Math.ceil((start_power / 100) * ftp)}-${Math.ceil((end_power / 100) * ftp)}-${duration}`
          })
          .join(';')

      setSelectedWorkout(workout)
      setSelectedWorkoutTargetDuration(totalDuration)

      // Set up event listeners after state update
      setTimeout(() => {
        const startButton = document.getElementById('start-app')
        const stopButton = document.getElementById('stop-app')
        const statusBox = document.getElementById('status-box')

        const updateStatus = (message: string) => {
          const timestamp = new Date().toLocaleTimeString()
          if (statusBox) {
            statusBox.textContent = `[${timestamp}] ${message}`
          }
        }

        if (startButton) {
          startButton.addEventListener('click', () => {
            const appName = 'overlay'
            updateStatus(`Renderer: requesting ${appName}...`)
            window.workoutAPI.startApp(appName, '-workout', workoutString)
            startButton.setAttribute('disabled', 'true')
            stopButton?.removeAttribute('disabled')
          })
        }

        if (stopButton) {
          stopButton.addEventListener('click', () => {
            const appName = 'overlay'
            updateStatus(`Renderer: stopping ${appName}...`)
            window.workoutAPI.stopApp(appName)
            startButton?.removeAttribute('disabled')
            stopButton.setAttribute('disabled', 'true')
          })
        }

        window.workoutAPI.onAppStatus((_event, message) => {
          updateStatus(`Main Process: ${message}`)
        })

        window.workoutAPI.onAppExited(() => {
          if (startButton) startButton.removeAttribute('disabled')
          if (stopButton) stopButton.setAttribute('disabled', 'true')
        })
      }, 0)
    },
    [ftp]
  )

  const unselectWorkout = useCallback(() => {
    window.workoutAPI.stopApp('overlay')
    setSelectedWorkout(null)
    setSelectedWorkoutTargetDuration(null)
  }, [])

  const handleSortChange = useCallback((event: React.ChangeEvent<HTMLSelectElement>) => {
    setSortOrder(event.target.value as SortOrder)
  }, [])

  const handleMinDurationChange = useCallback(
    (value: number) => {
      setFilterMin(Math.min(value, filterMax))
    },
    [filterMax, setFilterMin]
  )

  const handleMaxDurationChange = useCallback(
    (value: number) => {
      setFilterMax(Math.max(value, filterMin))
    },
    [filterMin, setFilterMax]
  )

  const handleSelectedWorkoutDurationChange = useCallback((newDuration: number) => {
    const parsed = Number(newDuration)
    if (!Number.isFinite(parsed)) return
    setSelectedWorkoutTargetDuration(parsed)
  }, [])

  const filteredWorkouts = useMemo(() => {
    const filtered = workouts.filter((workout) => {
      const duration = workout.totalDuration ?? getTotalDuration(workout)
      return duration >= filterMin && duration <= filterMax
    })

    if (sortOrder === 'asc') {
      return [...filtered].sort((a, b) => a.name.localeCompare(b.name))
    }

    if (sortOrder === 'desc') {
      return [...filtered].sort((a, b) => b.name.localeCompare(a.name))
    }

    if (sortOrder === 'duration_asc') {
      return [...filtered].sort((a, b) => (a.totalDuration ?? 0) - (b.totalDuration ?? 0))
    }

    if (sortOrder === 'duration_desc') {
      return [...filtered].sort((a, b) => (b.totalDuration ?? 0) - (a.totalDuration ?? 0))
    }

    return filtered
  }, [sortOrder, filterMin, filterMax, workouts])

  if (selectedWorkout) {
    const scaledWorkout = scaleWorkoutDuration(selectedWorkout, selectedWorkoutTargetDuration || 0)
    return (
      <WorkoutDetail
        workout={scaledWorkout!}
        onBack={unselectWorkout}
        desiredDuration={selectedWorkoutTargetDuration ?? scaledWorkout?.baseTotalDuration}
        onDurationChange={handleSelectedWorkoutDurationChange}
        ftp={ftp}
      />
    )
  }

  if (view === 'create') {
    return <WorkoutCreator onSave={handleSaveWorkout} onBack={showListView} ftp={ftp} onFtpChange={handleFtpChange} />
  }

  if (view === 'settings') {
    return <Settings ftp={ftp} onFtpChange={handleFtpChange} onBack={showListView} />
  }

  if (view === 'training-list') {
    return <TrainingList onBack={showListView} />
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Workouts</h1>
        <div className="flex items-center gap-4">
          <button
            onClick={showCreateWorkoutView}
            className="inline-flex items-center gap-2 rounded-lg bg-sky-500 hover:bg-sky-400 transition-colors px-4 py-2 font-semibold text-slate-900"
          >
            Create Workout
          </button>
          <button
            onClick={showTrainingListView}
            className="inline-flex items-center gap-2 rounded-lg bg-green-600 hover:bg-green-500 transition-colors px-4 py-2 font-semibold text-slate-100"
          >
            Training Sessions
          </button>
          <button
            onClick={showSettingsView}
            className="inline-flex items-center gap-2 rounded-lg bg-slate-600 hover:bg-slate-500 transition-colors px-4 py-2 font-semibold text-slate-100"
          >
            Settings
          </button>
        </div>
      </div>
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-4 space-y-4">
        <div className="flex flex-col gap-4 md:flex-row">
          <div className="flex-1">
            <label className="block text-sm font-semibold text-slate-300 mb-2">Sort</label>
            <select
              className="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-sky-400"
              value={sortOrder}
              onChange={handleSortChange}
            >
              <option value="default">Default order</option>
              <option value="asc">Name A → Z</option>
              <option value="desc">Name Z → A</option>
              <option value="duration_asc">Duration (shortest first)</option>
              <option value="duration_desc">Duration (longest first)</option>
            </select>
          </div>
          {durationMin == 0 && durationMax == 0 ? null : (
            <div className="flex-1 space-y-2">
              <DurationRangeSlider
                min={durationMin}
                max={durationMax}
                lowerValue={filterMin}
                upperValue={filterMax}
                step={60}
                onMinChange={handleMinDurationChange}
                onMaxChange={handleMaxDurationChange}
              />
            </div>
          )}
        </div>
      </div>
      <WorkoutList workouts={filteredWorkouts} onSelectWorkout={selectWorkout} ftp={ftp} />
    </div>
  )
}

export default App
