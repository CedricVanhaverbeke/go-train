import React from 'react'
import { Workout } from '../types'
import { powerToColor, getTotalDuration, formatDuration } from '../lib/utils'

interface WorkoutPreviewProps {
  workout: Workout
  width?: number
  height?: number
  ftp: number
}

const WorkoutPreview: React.FC<WorkoutPreviewProps> = ({ workout, width = 200, height = 100, ftp }) => {
  const totalDuration = workout.totalDuration ?? getTotalDuration(workout)
  const barWidthScale = width / totalDuration
  const formattedDuration = formatDuration(totalDuration)

  let accumulatedDuration = 0

  return (
    <div className="relative inline-flex">
      <svg width={width} height={height} className="bg-slate-700 rounded">
        {workout.steps.map((step, index) => {
          const startPower = step.start_power || 0
          const endPower = step.end_power || 0
          const barWidth = step.duration * barWidthScale
          const startHeight = (startPower / ftp) * height
          const endHeight = (endPower / ftp) * height
          const x = accumulatedDuration * barWidthScale
          accumulatedDuration += step.duration
          const points = `${x},${height} ${x},${height - startHeight} ${x + barWidth},${height - endHeight} ${x + barWidth},${height}`
          const avgPower = (startPower + endPower) / 2

          return <polygon key={index} points={points} fill={powerToColor(avgPower, ftp)} />
        })}
      </svg>
      <span className="absolute bottom-2 right-2 rounded bg-slate-900/80 px-2 py-1 text-xs font-semibold text-slate-100">
        {formattedDuration}
      </span>
    </div>
  )
}

export default WorkoutPreview

