import React from 'react'

interface DurationAdjusterProps {
  baseDuration: number
  currentDuration: number
  onChange: (duration: number) => void
  stepsCount?: number
}

const DurationAdjuster: React.FC<DurationAdjusterProps> = ({
  baseDuration,
  currentDuration,
  onChange,
  stepsCount = 1,
}) => {
  if (!Number.isFinite(baseDuration) || baseDuration <= 0) {
    return null
  }

  const safeBase = Math.max(baseDuration, 1)
  const safeCurrent = Math.max(Number(currentDuration) || safeBase, 1)
  const sliderLowerBound = Math.max(stepsCount, 10)
  const sliderMin = Math.max(sliderLowerBound, Math.min(Math.round(safeBase * 0.5), Math.floor(safeCurrent * 0.8)))
  const sliderMax = Math.max(Math.round(safeBase * 1.5), Math.ceil(safeCurrent * 1.2), sliderMin + 60)
  const sliderRange = Math.max(sliderMax - sliderMin, 1)
  const clampedCurrent = Math.min(Math.max(safeCurrent, sliderMin), sliderMax)
  const sliderPercent = ((clampedCurrent - sliderMin) / sliderRange) * 100
  const scalePercent = Math.round((safeCurrent / safeBase) * 100)

  const formatDuration = (seconds: number): string => {
    if (!Number.isFinite(seconds) || seconds < 0) return '00:00:00'
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    const secs = Math.floor(seconds % 60)
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between text-sm font-semibold text-slate-200">
        <span>Adjust duration</span>
        <span>
          {formatDuration(safeCurrent)} ({scalePercent}% of original)
        </span>
      </div>
      <div className="relative h-10 flex items-center">
        <div className="absolute inset-x-0 h-1 bg-slate-700 rounded-full"></div>
        <div className="absolute h-1 bg-sky-400 rounded-full" style={{ left: '0%', width: `${sliderPercent}%` }}></div>
        <input
          type="range"
          min={sliderMin}
          max={sliderMax}
          step={30}
          value={clampedCurrent}
          onChange={(e) => onChange?.(Number(e.target.value))}
          className="range-thumb absolute inset-0 w-full h-full focus:outline-none"
          style={{
            WebkitAppearance: 'none',
            appearance: 'none',
            background: 'transparent',
            width: '100%',
            height: '100%',
            pointerEvents: 'auto',
          }}
        />
      </div>
      <div className="flex items-center gap-3 text-sm text-slate-300">
        <label className="text-xs font-semibold uppercase tracking-wide text-slate-400">Minutes</label>
        <input
          type="number"
          min={sliderMin / 60}
          step={0.5}
          value={Number((safeCurrent / 60).toFixed(2))}
          onChange={(e) => {
            const minutes = Number(e.target.value)
            if (!Number.isFinite(minutes)) return
            const secondsValue = minutes * 60
            onChange?.(Math.max(sliderMin, secondsValue))
          }}
          className="w-20 rounded-lg bg-slate-900 border border-slate-700 px-2 py-1 text-right"
        />
        <span className="text-xs text-slate-500">min</span>
        <button
          type="button"
          onClick={() => onChange?.(safeBase)}
          className="ml-auto text-xs font-semibold text-sky-400 hover:text-sky-300"
        >
          Reset
        </button>
      </div>
      <p className="text-xs text-slate-500">Steps are scaled linearly so the workout profile stays intact.</p>
    </div>
  )
}

export default DurationAdjuster
