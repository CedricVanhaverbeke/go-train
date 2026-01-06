import { useState, useEffect } from 'react'
import { getWorkouts, getTotalDuration } from './lib/utils'
import type { Workout } from './types'
import App from './app'

interface AppWrapperProps {
  initialWorkouts?: Workout[]
  initialDurationMin?: number
  initialDurationMax?: number
  initialFilterMin?: number
  initialFilterMax?: number
}

export function AppWrapper({
  initialWorkouts = [],
  initialDurationMin = 0,
  initialDurationMax = 0,
  initialFilterMin = 0,
  initialFilterMax = 0,
}: AppWrapperProps) {
  const [workouts, setWorkouts] = useState<Workout[]>(initialWorkouts)
  const [durationMin, setDurationMin] = useState(initialDurationMin)
  const [durationMax, setDurationMax] = useState(initialDurationMax)
  const [filterMin, setFilterMin] = useState(initialFilterMin)
  const [filterMax, setFilterMax] = useState(initialFilterMax)
  const [isLoading, setIsLoading] = useState(!initialWorkouts.length)

  useEffect(() => {
    if (initialWorkouts.length > 0) {
      setIsLoading(false)
      return
    }

    getWorkouts().then((fetchedWorkouts) => {
      const workoutsWithDuration = fetchedWorkouts.map((workout) => ({
        ...workout,
        totalDuration: getTotalDuration(workout),
      }))
      const durationValues = workoutsWithDuration.map((w) => w.totalDuration || 0)
      const minDuration = durationValues.length ? Math.min(...durationValues) : 0
      const maxDuration = durationValues.length ? Math.max(...durationValues) : 3600

      setWorkouts(workoutsWithDuration)
      setDurationMin(minDuration)
      setDurationMax(maxDuration)
      setFilterMin(minDuration)
      setFilterMax(maxDuration)
      setIsLoading(false)
    })
  }, [initialWorkouts.length])

  if (isLoading) {
    return <div className="flex items-center justify-center h-screen">Loading workouts...</div>
  }

  return (
    <App
      workouts={workouts}
      durationMin={durationMin}
      durationMax={durationMax}
      filterMin={filterMin}
      filterMax={filterMax}
      setWorkouts={setWorkouts}
      setDurationMin={setDurationMin}
      setDurationMax={setDurationMax}
      setFilterMin={setFilterMin}
      setFilterMax={setFilterMax}
    />
  )
}
