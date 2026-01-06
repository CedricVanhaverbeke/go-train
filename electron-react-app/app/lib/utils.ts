import { Workout } from '../types'

// --- Data Fetching ---
export async function getWorkouts(): Promise<Workout[]> {
  const response = await fetch('./assets/data.json')
  return await response.json()
}

// --- Color Scale ---
export function powerToColor(power: number, ftp: number): string {
  const percentage = Math.min(power / ftp, 1)
  const hue = (1 - percentage) * 120 // 120 (green) -> 0 (red)
  return `hsl(${hue}, 100%, 50%)`
}

// --- Utilities ---
export function getTotalDuration(workout: Workout): number {
  return workout.steps.reduce((sum, step) => sum + step.duration, 0)
}

export function formatDuration(seconds: number): string {
  if (!Number.isFinite(seconds) || seconds < 0) return '00:00:00'
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

export function scaleWorkoutDuration(workout: Workout, targetDurationSeconds: number): Workout | null {
  if (!workout) return null
  const baseTotal = workout.totalDuration ?? getTotalDuration(workout)
  if (!baseTotal) {
    return {
      ...workout,
      baseTotalDuration: 0,
      durationScale: 1,
      steps: workout.steps.map((step) => ({ ...step })),
    }
  }

  const rawTarget = Number(targetDurationSeconds)
  const safeTarget = Number.isFinite(rawTarget) ? rawTarget : baseTotal
  const minPossible = workout.steps.length || 1
  const target = Math.max(Math.round(safeTarget), minPossible)
  const scale = target / baseTotal

  const rawDurations = workout.steps.map((step) => step.duration * scale)
  const flooredDurations = rawDurations.map((raw) => Math.max(1, Math.floor(raw)))
  const sum = flooredDurations.reduce((acc, duration) => acc + duration, 0)
  let remaining = Math.max(target - sum, 0)

  if (remaining > 0) {
    const remainderOrder = rawDurations
      .map((raw, index) => ({ index, remainder: raw - Math.floor(raw) }))
      .sort((a, b) => b.remainder - a.remainder)

    for (let i = 0; i < remainderOrder.length && remaining > 0; i += 1) {
      flooredDurations[remainderOrder[i].index] += 1
      remaining -= 1
    }
  }

  const scaledSteps = workout.steps.map((step, index) => ({
    ...step,
    duration: flooredDurations[index],
  }))
  const scaledTotalDuration = scaledSteps.reduce((sumDuration, step) => sumDuration + step.duration, 0)

  return {
    ...workout,
    steps: scaledSteps,
    totalDuration: scaledTotalDuration,
    targetDuration: target,
    baseTotalDuration: baseTotal,
    durationScale: scaledTotalDuration / baseTotal,
  }
}
