import React, { useState, useEffect } from 'react'
import { WorkoutStep } from '../types'
import { formatDuration } from '../lib/utils'

interface VisualWorkoutEditorProps {
  steps: WorkoutStep[]
  onStepsChange: (steps: WorkoutStep[]) => void
  ftp: number
  onFtpChange: (ftp: number) => void
}

const VisualWorkoutEditor: React.FC<VisualWorkoutEditorProps> = ({ steps, onStepsChange, ftp, onFtpChange }) => {
  const [editingStep, setEditingStep] = useState<number | null>(null)
  const [draggingIndex, setDraggingIndex] = useState<number | null>(null)
  const [dropIndex, setDropIndex] = useState<number | null>(null)

  useEffect(() => {
    const style = document.createElement('style')
    style.textContent = `
      .range-thumb {
        -webkit-appearance: none;
        appearance: none;
        background: transparent;
        width: 100%;
        height: 100%;
        pointer-events: none;
      }
      .range-thumb::-webkit-slider-runnable-track {
        background: transparent;
      }
      .range-thumb::-moz-range-track {
        background: transparent;
        border: none;
      }
      .range-thumb::-webkit-slider-thumb {
        -webkit-appearance: none;
        height: 18px;
        width: 18px;
        border-radius: 9999px;
        background: #38bdf8;
        border: 2px solid #0f172a;
        pointer-events: auto;
        cursor: pointer;
        box-shadow: 0 0 0 2px rgba(15, 23, 42, 0.5);
      }
      .range-thumb::-moz-range-thumb {
        height: 18px;
        width: 18px;
        border-radius: 9999px;
        background: #38bdf8;
        border: 2px solid #0f172a;
        pointer-events: auto;
        cursor: pointer;
        box-shadow: 0 0 0 2px rgba(15, 23, 42, 0.5);
      }
    `
    document.head.appendChild(style)
    return () => {
      document.head.removeChild(style)
    }
  }, [])

  const handleStepChange = (index: number, field: keyof WorkoutStep, value: number) => {
    const newSteps = [...steps]
    newSteps[index] = { ...newSteps[index], [field]: Math.max(0, value) }
    onStepsChange(newSteps)
  }

  const addStep = () => {
    const newSteps = [...steps, { duration: 60, start_power: 100, end_power: 100 }]
    onStepsChange(newSteps)
  }

  const removeStep = (index: number) => {
    const newSteps = steps.filter((_, i) => i !== index)
    onStepsChange(newSteps)
    setEditingStep(null)
  }

  const handleStepClick = (e: React.MouseEvent, index: number) => {
    e.stopPropagation()
    setEditingStep(editingStep === index ? null : index)
  }

  // --- Drag and Drop Handlers ---
  const handleDragStart = (e: React.DragEvent, index: number) => {
    setDraggingIndex(index)
    setEditingStep(null)
    e.dataTransfer.effectAllowed = 'move'
  }

  const handleDragOver = (e: React.DragEvent, index: number) => {
    e.preventDefault()
    if (draggingIndex === null || draggingIndex === index) {
      setDropIndex(null)
      return
    }
    setDropIndex(index)
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    if (draggingIndex === null || dropIndex === null || draggingIndex === dropIndex) {
      return
    }

    const newSteps = [...steps]
    const [draggedStep] = newSteps.splice(draggingIndex, 1)
    newSteps.splice(dropIndex, 0, draggedStep)

    onStepsChange(newSteps)
    setDraggingIndex(null)
    setDropIndex(null)
  }

  const handleDragEnd = () => {
    setDraggingIndex(null)
    setDropIndex(null)
  }

  // --- Resize Handlers ---
  const handleDurationResize = (e: React.MouseEvent, index: number) => {
    e.stopPropagation()
    const startX = e.clientX
    const startDuration = steps[index].duration

    const onMouseMove = (moveEvent: MouseEvent) => {
      const newDuration = startDuration + (moveEvent.clientX - startX)
      handleStepChange(index, 'duration', newDuration)
    }

    const onMouseUp = () => {
      document.removeEventListener('mousemove', onMouseMove)
      document.removeEventListener('mouseup', onMouseUp)
    }

    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
  }

  const handlePowerResize = (e: React.MouseEvent, index: number) => {
    e.stopPropagation()
    const startY = e.clientY
    const startPower = steps[index].start_power
    const containerElement = (e.currentTarget as HTMLElement).closest('.h-64') as HTMLElement
    const containerHeight = containerElement?.offsetHeight || 256
    const maxPower = ftp

    const onMouseMove = (moveEvent: MouseEvent) => {
      const deltaY = startY - moveEvent.clientY
      const powerChange = (deltaY / containerHeight) * maxPower
      handleStepChange(index, 'start_power', startPower + powerChange)
      handleStepChange(index, 'end_power', startPower + powerChange)
    }

    const onMouseUp = () => {
      document.removeEventListener('mousemove', onMouseMove)
      document.removeEventListener('mouseup', onMouseUp)
    }

    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
  }

  const totalDuration = steps.reduce((sum, step) => sum + step.duration, 0)

  return (
    <div className="space-y-4" onClick={() => setEditingStep(null)}>
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold text-slate-300">Visual Workout Editor</h3>
        <div className="flex items-center gap-4">
          <span className="text-sm text-slate-400">FTP: {ftp}W</span>
        </div>
      </div>
      <div
        className="relative bg-slate-700 h-64 w-full rounded-lg overflow-x-auto flex"
        onDragOver={(e) => e.preventDefault()}
        onDrop={handleDrop}
      >
        {steps.map((step, index) => {
          const isDragging = draggingIndex === index
          const isDropTarget = dropIndex === index
          return (
            <div
              key={index}
              draggable
              onDragStart={(e) => handleDragStart(e, index)}
              onDragOver={(e) => handleDragOver(e, index)}
              onDragEnd={handleDragEnd}
              onClick={(e) => handleStepClick(e, index)}
              className={`relative h-full group transition-all duration-150 ${
                isDragging ? 'opacity-50' : 'opacity-100'
              } ${isDropTarget ? 'bg-green-500/20' : ''}`}
              style={{
                width: `${step.duration}px`,
                cursor: 'grab',
              }}
            >
              <div
                className="absolute bottom-0 w-full bg-sky-500 group-hover:bg-sky-400 transition-colors"
                style={{
                  height: `${(step.start_power / ftp) * 100}%`,
                  pointerEvents: 'none',
                }}
              />

              {editingStep === index && (
                <>
                  <div
                    onMouseDown={(e) => handleDurationResize(e, index)}
                    className="absolute top-1/2 right-0 w-2 h-4 bg-white rounded-sm cursor-ew-resize z-20"
                    style={{ transform: 'translate(50%, -50%)' }}
                  />
                  <div
                    onMouseDown={(e) => handlePowerResize(e, index)}
                    className="absolute left-1/2 w-4 h-2 bg-white rounded-sm cursor-ns-resize z-20"
                    style={{
                      top: `${100 - (step.start_power / ftp) * 100}%`,
                      transform: 'translate(-50%, -50%)',
                    }}
                  />

                  <div
                    className="absolute top-0 left-0 bg-slate-800 p-2 rounded-lg shadow-lg z-30 w-48"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <p className="text-xs font-bold text-slate-200 mb-2">Edit Step {index + 1}</p>
                    <div className="space-y-2">
                      <div>
                        <label className="block text-xs font-semibold text-slate-400 mb-1">Duration (s)</label>
                        <input
                          type="number"
                          value={step.duration}
                          onChange={(e) => handleStepChange(index, 'duration', Number(e.target.value))}
                          className="w-full rounded-md bg-slate-900 border border-slate-600 px-2 py-1"
                        />
                      </div>
                      <div>
                        <label className="block text-xs font-semibold text-slate-400 mb-1">Power (%)</label>
                        <input
                          type="number"
                          value={step.start_power}
                          onChange={(e) => handleStepChange(index, 'start_power', Number(e.target.value))}
                          className="w-full rounded-md bg-slate-900 border border-slate-600 px-2 py-1"
                        />
                      </div>
                      <div>
                        <label className="block text-xs font-semibold text-slate-400 mb-1">Power (W)</label>
                        <input
                          type="number"
                          value={Math.round((step.start_power / 100) * ftp)}
                          className="w-full rounded-md bg-slate-900 border border-slate-600 px-2 py-1"
                          readOnly
                        />
                      </div>
                    </div>
                    <button onClick={() => removeStep(index)} className="mt-3 text-xs text-red-400 hover:text-red-300">
                      Remove Step
                    </button>
                  </div>
                </>
              )}
            </div>
          )
        })}
      </div>
      <div className="flex justify-between items-center">
        <div className="flex gap-2">
          <button
            onClick={addStep}
            className="inline-flex items-center gap-2 rounded-lg bg-slate-600 hover:bg-slate-500 transition-colors px-4 py-2 font-semibold text-slate-100"
          >
            Add Step
          </button>
        </div>
        <div className="text-sm text-slate-400">Total Duration: {formatDuration(totalDuration)}</div>
      </div>
    </div>
  )
}

export default VisualWorkoutEditor

