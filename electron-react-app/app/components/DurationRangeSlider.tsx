import React, { useEffect, useState } from 'react'
import { Range } from 'react-range'
import { useDebounce } from '@uidotdev/usehooks'

interface DurationRangeSliderProps {
  min: number
  max: number
  lowerValue: number
  upperValue: number
  step?: number
  onMinChange: (value: number) => void
  onMaxChange: (value: number) => void
}

const DurationRangeSlider: React.FC<DurationRangeSliderProps> = ({
  min,
  max,
  lowerValue,
  upperValue,
  step = 60,
  onMinChange,
  onMaxChange,
}) => {
  const [values, setValues] = useState([lowerValue, upperValue])
  const [minLocal, maxLocal] = values

  const debouncedMin = useDebounce(minLocal, 1000)
  const debouncedMax = useDebounce(maxLocal, 1000)

  useEffect(() => {
    onMaxChange(debouncedMax)
  }, [debouncedMax, onMaxChange])

  useEffect(() => {
    onMinChange(debouncedMin)
  }, [debouncedMin, onMinChange])

  const formatDuration = (seconds: number): string => {
    if (!Number.isFinite(seconds) || seconds < 0) return '00:00:00'
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    const secs = Math.floor(seconds % 60)
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }

  const handleValuesChange = (newValues: number[]) => {
    setValues(newValues)
  }

  return (
    <div className="flex-1 space-y-2">
      <div className="flex items-center justify-between text-sm font-semibold text-slate-300">
        <span>Duration filter</span>
        <span>
          {formatDuration(values[0])} – {formatDuration(values[1])}
        </span>
      </div>
      <div className="relative h-10 flex items-center px-2 w-full">
        <Range
          values={values}
          step={step}
          min={min}
          max={max}
          onChange={handleValuesChange}
          renderTrack={({ props, children }) => (
            <div
              {...props}
              className="h-1 bg-slate-700 rounded-full w-full"
              style={{
                ...props.style,
              }}
            >
              <div
                className="h-1 bg-sky-400 rounded-full absolute"
                style={{
                  left: `${((values[0] - min) / (max - min)) * 100}%`,
                  width: `${((values[1] - values[0]) / (max - min)) * 100}%`,
                }}
              />
              {children}
            </div>
          )}
          renderThumb={({ props, isDragged }) => (
            <div
              {...props}
              className={`h-4 w-4 bg-white rounded-full shadow-lg border-2 ${isDragged ? 'border-sky-500' : 'border-slate-400'} focus:outline-none`}
              style={{
                ...props.style,
              }}
            />
          )}
        />
      </div>
      <div className="flex justify-between text-xs text-slate-500 px-2">
        <span>{formatDuration(min)}</span>
        <span>{formatDuration(max)}</span>
      </div>
    </div>
  )
}

export default DurationRangeSlider
