import React from 'react'
import { Workout } from '../types'
import { formatDuration } from '../lib/utils'
import WorkoutPreview from './WorkoutPreview'
import WorkoutSteps from './WorkoutSteps'
import DurationAdjuster from './DurationAdjuster'

interface WorkoutDetailProps {
  workout: Workout
  onBack: () => void
  desiredDuration?: number
  onDurationChange: (duration: number) => void
  ftp: number
}

const WorkoutDetail: React.FC<WorkoutDetailProps> = ({ workout, onBack, desiredDuration, onDurationChange, ftp }) => {
  const totalDuration = workout.steps.reduce((sum, step) => sum + step.duration, 0)
  const baseDuration = workout.baseTotalDuration ?? totalDuration
  const durationScalePercent = baseDuration ? Math.round((totalDuration / baseDuration) * 100) : 100

  return (
    <div className="space-y-6">
      <button onClick={onBack} className="text-sky-400 hover:underline">
        ← Back to Workouts
      </button>
      <h2 className="text-3xl font-bold">{workout.name}</h2>
      <div className="w-full max-w-3xl rounded-2xl bg-slate-800 shadow-xl border border-slate-700 p-8 space-y-6">
        <section className="space-y-3">
          <p className="text-sm font-semibold text-slate-200">Workout Preview</p>
          <div className="bg-slate-900 border border-slate-700 rounded-xl p-4 flex justify-center">
            <WorkoutPreview workout={workout} width={480} height={160} ftp={ftp} />
          </div>
        </section>
        <section className="space-y-3">
          <WorkoutSteps steps={workout.steps} ftp={ftp} />
        </section>
        <section className="space-y-3">
          <p className="text-sm font-semibold text-slate-200">Duration</p>
          <DurationAdjuster
            baseDuration={baseDuration}
            currentDuration={desiredDuration ?? workout.targetDuration ?? totalDuration}
            stepsCount={workout.steps.length}
            onChange={onDurationChange}
          />
          <div className="text-xs text-slate-400 flex flex-wrap gap-x-4 gap-y-1">
            <span>Original: {formatDuration(baseDuration)}</span>
            <span>Current: {formatDuration(totalDuration)}</span>
            <span>Scale: {durationScalePercent}%</span>
          </div>
        </section>
        <section className="space-y-3">
          <div className="flex items-center gap-4">
            <button
              id="start-app"
              className="inline-flex items-center gap-2 rounded-lg bg-sky-500 hover:bg-sky-400 transition-colors px-6 py-3 font-semibold text-slate-900"
            >
              <span>Start App</span>
            </button>
            <button
              id="stop-app"
              className="inline-flex items-center gap-2 rounded-lg bg-red-500 hover:bg-red-400 transition-colors px-6 py-3 font-semibold text-slate-900 disabled:bg-slate-600 disabled:cursor-not-allowed"
              disabled
            >
              <span>Stop App</span>
            </button>
          </div>
          <div>
            <p className="text-sm font-semibold text-slate-200 mb-2">Status</p>
            <pre
              id="status-box"
              className="bg-slate-900 rounded-lg border border-slate-700 p-4 text-green-400 text-sm overflow-auto min-h-[4rem]"
            >
              Waiting for user action...
            </pre>
          </div>
        </section>
      </div>
    </div>
  )
}

export default WorkoutDetail
