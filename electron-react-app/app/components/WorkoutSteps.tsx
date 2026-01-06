import React from 'react'
import { WorkoutStep } from '../types'

interface WorkoutStepsProps {
  steps: WorkoutStep[]
  ftp: number
}

const WorkoutSteps: React.FC<WorkoutStepsProps> = ({ steps, ftp }) => {
  const formatStepDuration = (seconds: number): string => {
    if (!Number.isFinite(seconds) || seconds < 0) return '0s'
    const min = Math.floor(seconds / 60)
    const sec = Math.round(seconds % 60)
    return min > 0 ? `${min}m ${sec}s` : `${sec}s`
  }

  return (
    <div className="space-y-2">
      <h3 className="text-lg font-semibold text-slate-300">Workout Steps</h3>
      <ul className="divide-y divide-slate-700">
        {steps.map((step, index) => {
          const avgPower = (step.start_power + step.end_power) / 2
          const wattage = Math.ceil((avgPower / 100) * ftp)
          return (
            <li key={index} className="py-2 flex justify-between items-center">
              <span className="text-slate-400">
                Step {index + 1}: {formatStepDuration(step.duration)}
              </span>
              <span className="font-semibold text-sky-400">{wattage} W</span>
            </li>
          )
        })}
      </ul>
    </div>
  )
}

export default WorkoutSteps

