import React from 'react'

interface SettingsProps {
  ftp: number
  onFtpChange: (ftp: number) => void
  onBack: () => void
}

const Settings: React.FC<SettingsProps> = ({ ftp, onFtpChange, onBack }) => {
  const handleFtpChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onFtpChange(Number(e.target.value))
  }

  const saveFtp = () => {
    onFtpChange(ftp)
    onBack()
  }

  return (
    <div className="space-y-6">
      <button onClick={onBack} className="text-sky-400 hover:underline">
        ← Back to Workouts
      </button>
      <h2 className="text-3xl font-bold">Settings</h2>
      <div className="w-full max-w-3xl rounded-2xl bg-slate-800 shadow-xl border border-slate-700 p-8 space-y-6">
        <p className="text-slate-400">
          Update your Functional Threshold Power (FTP).
        </p>
        <div className="flex items-center gap-4">
          <input
            type="number"
            value={ftp}
            onChange={handleFtpChange}
            className="w-32 rounded-md bg-slate-900 border border-slate-600 px-3 py-2 text-center text-lg"
          />
          <button
            onClick={saveFtp}
            className="rounded-lg bg-sky-600 hover:bg-sky-500 transition-colors px-6 py-2 font-semibold text-white"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  )
}

export default Settings