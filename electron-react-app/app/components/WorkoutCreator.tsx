import React, { useState } from 'react'
import { Workout } from '../types'
import VisualWorkoutEditor from './VisualWorkoutEditor'

interface WorkoutCreatorProps {
  onSave: (workout: Workout) => void
  onBack: () => void
  ftp: number
  onFtpChange: (ftp: number) => void
}

const WorkoutCreator: React.FC<WorkoutCreatorProps> = ({ onSave, onBack, ftp, onFtpChange }) => {
  const [name, setName] = useState('')
  const [steps, setSteps] = useState<Workout['steps']>([{ duration: 60, start_power: 100, end_power: 100 }])

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value)
  }

  const handleStepsChange = (newSteps: Workout['steps']) => {
    setSteps(newSteps)
  }

  const handleSave = () => {
    if (name.trim() === '' || steps.length === 0) {
      alert('Please provide a name and at least one step for the workout.')
      return
    }
    const newWorkout: Workout = {
      name: name.trim(),
      steps,
    }
    onSave(newWorkout)
  }

  return (
    <div className="space-y-6">
      <button onClick={onBack} className="text-sky-400 hover:underline">
        ← Back to Workouts
      </button>
      <h2 className="text-3xl font-bold">Create Workout</h2>
      <div className="w-full max-w-3xl rounded-2xl bg-slate-800 shadow-xl border border-slate-700 p-8 space-y-6">
        <div className="space-y-2">
          <label className="block text-sm font-semibold text-slate-300">Workout Name</label>
          <input
            type="text"
            value={name}
            onChange={handleNameChange}
            className="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-sky-400"
            placeholder="e.g., Morning Spin"
          />
        </div>

        <VisualWorkoutEditor steps={steps} onStepsChange={handleStepsChange} ftp={ftp} onFtpChange={onFtpChange} />

        <div className="flex items-center gap-4">
          <button
            onClick={handleSave}
            className="inline-flex items-center gap-2 rounded-lg bg-sky-500 hover:bg-sky-400 transition-colors px-6 py-3 font-semibold text-slate-900"
          >
            Save Workout
          </button>
          <button onClick={onBack} className="text-slate-400 hover:text-slate-300">
            Cancel
          </button>
        </div>
      </div>
    </div>
  )
}

export default WorkoutCreator

