import React from 'react'
import { Workout } from '../types'
import WorkoutPreview from './WorkoutPreview'

interface WorkoutCardProps {
  workout: Workout
  onClick: () => void
  ftp: number
}

const WorkoutCard: React.FC<WorkoutCardProps> = ({ workout, onClick, ftp }) => {
  return (
    <div
      className="bg-slate-800 rounded-lg shadow-lg p-4 flex flex-col gap-4 cursor-pointer hover:bg-slate-700 transition-colors"
      onClick={onClick}
    >
      <h3 className="font-bold text-lg">{workout.name}</h3>
      <WorkoutPreview workout={workout} ftp={ftp} />
    </div>
  )
}

export default WorkoutCard
