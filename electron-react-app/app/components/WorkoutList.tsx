import React from 'react'
import { Workout } from '../types'
import WorkoutCard from './WorkoutCard'

interface WorkoutListProps {
  workouts: Workout[]
  onSelectWorkout: (workout: Workout, ftp: number) => void
  ftp: number
}

const WorkoutList: React.FC<WorkoutListProps> = ({ workouts, onSelectWorkout, ftp }) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {workouts.map((workout, i) => (
        <WorkoutCard key={i} workout={workout} onClick={() => onSelectWorkout(workout, ftp)} ftp={ftp} />
      ))}
    </div>
  )
}

export default WorkoutList
