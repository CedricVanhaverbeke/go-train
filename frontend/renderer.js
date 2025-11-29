const { ipcRenderer } = require('electron');

const { h, Component, render } = window;
const html = window.htm.bind(h);

// --- Data Fetching ---
async function getWorkouts() {
  const response = await fetch('./data.json');
  return await response.json();
}

// --- Color Scale ---
function powerToColor(power, maxPower) {
  const percentage = Math.min(power / maxPower, 1);
  const hue = (1 - percentage) * 120; // 120 (green) -> 0 (red)
  return `hsl(${hue}, 100%, 50%)`;
}

// --- Utilities ---
function getTotalDuration(workout) {
  return workout.steps.reduce((sum, step) => sum + step.duration, 0);
}

function formatDuration(seconds) {
  if (!Number.isFinite(seconds) || seconds < 0) return '0 min';
  const minutes = Math.round(seconds / 60);
  return minutes >= 60
    ? `${(minutes / 60).toFixed(minutes % 60 === 0 ? 0 : 1)} hr`
    : `${minutes} min`;
}

// --- Components ---

function WorkoutPreview({ workout, maxPower }) {
  const totalDuration = workout.totalDuration ?? getTotalDuration(workout);
  const barWidthScale = 200 / totalDuration; // Fixed width for the preview

  let accumulatedDuration = 0;
  return html`
    <svg width="200" height="100" class="bg-slate-700 rounded">
      ${workout.steps.map(step => {
        const barWidth = step.duration * barWidthScale;
        const barHeight = (step.power / maxPower) * 100;
        const x = accumulatedDuration * barWidthScale;
        accumulatedDuration += step.duration;
        return html`
          <rect
            x=${x}
            y=${100 - barHeight}
            width=${barWidth}
            height=${barHeight}
            fill=${powerToColor(step.power, maxPower)}
          />
        `;
      })}
    </svg>
  `;
}

function WorkoutCard({ workout, maxPower, onClick }) {
  return html`
    <div class="bg-slate-800 rounded-lg shadow-lg p-4 flex flex-col gap-4 cursor-pointer hover:bg-slate-700 transition-colors" onClick=${onClick}>
      <h3 class="font-bold text-lg">${workout.name}</h3>
      <${WorkoutPreview} workout=${workout} maxPower=${maxPower} />
    </div>
  `;
}

function DurationRangeSlider({
  min,
  max,
  lowerValue,
  upperValue,
  step = 60,
  onMinChange,
  onMaxChange,
}) {
  const safeRange = Math.max(max - min, 1);
  const clampedLower = Math.min(lowerValue, upperValue);
  const clampedUpper = Math.max(lowerValue, upperValue);
  const lowerPercent = ((clampedLower - min) / safeRange) * 100;
  const upperPercent = ((clampedUpper - min) / safeRange) * 100;

  return html`
    <div class="flex-1 space-y-2">
      <div class="flex items-center justify-between text-sm font-semibold text-slate-300">
        <span>Duration filter</span>
        <span>${formatDuration(clampedLower)} – ${formatDuration(clampedUpper)}</span>
      </div>
      <div class="relative h-10 flex items-center">
        <div class="absolute inset-x-0 h-1 bg-slate-700 rounded-full"></div>
        <div
          class="absolute h-1 bg-sky-400 rounded-full"
          style=${{
            left: `${lowerPercent}%`,
            width: `${Math.max(upperPercent - lowerPercent, 0)}%`,
          }}
        ></div>
        <input
          type="range"
          min=${min}
          max=${max}
          value=${clampedLower}
          step=${step}
          onInput=${onMinChange}
          class="range-thumb absolute inset-0 w-full h-full focus:outline-none"
        />
        <input
          type="range"
          min=${min}
          max=${max}
          value=${clampedUpper}
          step=${step}
          onInput=${onMaxChange}
          class="range-thumb absolute inset-0 w-full h-full focus:outline-none"
        />
      </div>
      <div class="flex justify-between text-xs text-slate-500">
        <span>${formatDuration(min)}</span>
        <span>${formatDuration(max)}</span>
      </div>
    </div>
  `;
}

function WorkoutDetail({ workout, onBack }) {
  const totalDuration = workout.steps.reduce((sum, step) => sum + step.duration, 0);
  const maxPower = Math.max(...workout.steps.map(s => s.power));

  return html`
    <div class="space-y-6">
      <button onClick=${onBack} class="text-sky-400 hover:underline">← Back to Workouts</button>
      <h2 class="text-3xl font-bold">${workout.name}</h2>
      <div class="w-full max-w-3xl rounded-2xl bg-slate-800 shadow-xl border border-slate-700 p-8 space-y-6">
        <section class="space-y-3">
          <div class="flex items-center gap-4">
            <button id="start-app" class="inline-flex items-center gap-2 rounded-lg bg-sky-500 hover:bg-sky-400 transition-colors px-6 py-3 font-semibold text-slate-900">
              <span>Start App</span>
            </button>
            <button id="stop-app" class="inline-flex items-center gap-2 rounded-lg bg-red-500 hover:bg-red-400 transition-colors px-6 py-3 font-semibold text-slate-900 disabled:bg-slate-600 disabled:cursor-not-allowed" disabled>
              <span>Stop App</span>
            </button>
          </div>
          <div>
            <p class="text-sm font-semibold text-slate-200 mb-2">Status</p>
            <pre id="status-box" class="bg-slate-900 rounded-lg border border-slate-700 p-4 text-green-400 text-sm overflow-auto min-h-[4rem]">Waiting for user action...</pre>
          </div>
        </section>
      </div>
    </div>
  `;
}

class App extends Component {
  state = {
    workouts: [],
    selectedWorkout: null,
    maxPower: 0,
    sortOrder: 'default',
    durationMin: 0,
    durationMax: 0,
    filterMin: 0,
    filterMax: 0,
  };

  componentDidMount() {
    getWorkouts().then(workouts => {
      const workoutsWithDuration = workouts.map(workout => ({
        ...workout,
        totalDuration: getTotalDuration(workout),
      }));
      const powerValues = workoutsWithDuration.flatMap(w => w.steps.map(s => s.power));
      const maxPower = powerValues.length ? Math.max(...powerValues) : 0;
      const durationValues = workoutsWithDuration.map(w => w.totalDuration);
      const durationMin = durationValues.length ? Math.min(...durationValues) : 0;
      const durationMax = durationValues.length ? Math.max(...durationValues) : 0;

      this.setState({
        workouts: workoutsWithDuration,
        maxPower,
        durationMin,
        durationMax,
        filterMin: durationMin,
        filterMax: durationMax,
      });
    });
  }
  
  selectWorkout = workout => {
    this.setState({ selectedWorkout: workout }, () => {
        // we need to re-add the event listeners after the DOM is updated
        const startButton = document.getElementById('start-app');
        const stopButton = document.getElementById('stop-app');
        const statusBox = document.getElementById('status-box');

        const updateStatus = (message) => {
            const timestamp = new Date().toLocaleTimeString();
            if (statusBox) {
                statusBox.textContent = `[${timestamp}] ${message}`;
            }
        };

        if (startButton) {
            startButton.addEventListener('click', () => {
                const appName = "overlay";
                updateStatus(`Renderer: requesting ${appName}...`);
                ipcRenderer.send('START_APP', appName);
                startButton.disabled = true;
                stopButton.disabled = false;
            });
        }

        if (stopButton) {
            stopButton.addEventListener('click', () => {
                const appName = "overlay";
                updateStatus(`Renderer: stopping ${appName}...`);
                ipcRenderer.send('STOP_APP', appName);
                startButton.disabled = false;
                stopButton.disabled = true;
            });
        }

        ipcRenderer.on('APP_STATUS', (_event, message) => {
            updateStatus(`Main Process: ${message}`);
        });

        ipcRenderer.on('APP_OUTPUT', (_event, data) => {
            console.log(data);
        });

        ipcRenderer.on('APP_EXITED', () => {
            if (startButton) startButton.disabled = false;
            if (stopButton) stopButton.disabled = true;
        });
    });
  };

  unselectWorkout = () => this.setState({ selectedWorkout: null });

  handleSortChange = event => {
    this.setState({ sortOrder: event.target.value });
  };

  handleMinDurationChange = event => {
    const newMin = Number(event.target.value);
    this.setState(prev => ({ filterMin: Math.min(newMin, prev.filterMax) }));
  };

  handleMaxDurationChange = event => {
    const newMax = Number(event.target.value);
    this.setState(prev => ({ filterMax: Math.max(newMax, prev.filterMin) }));
  };

  getFilteredWorkouts() {
    const { workouts, sortOrder, filterMin, filterMax } = this.state;
    const filtered = workouts.filter(workout => {
      const duration = workout.totalDuration ?? getTotalDuration(workout);
      return duration >= filterMin && duration <= filterMax;
    });

    if (sortOrder === 'asc') {
      return [...filtered].sort((a, b) => a.name.localeCompare(b.name));
    }

    if (sortOrder === 'desc') {
      return [...filtered].sort((a, b) => b.name.localeCompare(a.name));
    }

    return filtered;
  }

  render(
    _,
    {
      selectedWorkout,
      maxPower,
      sortOrder,
      durationMin,
      durationMax,
      filterMin,
      filterMax,
    }
  ) {
    const workouts = this.getFilteredWorkouts();

    if (selectedWorkout) {
      return html`<${WorkoutDetail} workout=${selectedWorkout} onBack=${this.unselectWorkout} />`;
    }

    return html`
      <div class="space-y-6">
        <h1 class="text-3xl font-bold">Workouts</h1>
        <div class="bg-slate-800 border border-slate-700 rounded-xl p-4 space-y-4">
          <div class="flex flex-col gap-4 md:flex-row">
            <div class="flex-1">
              <label class="block text-sm font-semibold text-slate-300 mb-2">Sort</label>
              <select
                class="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-sky-400"
                value=${sortOrder}
                onChange=${this.handleSortChange}
              >
                <option value="default">Default order</option>
                <option value="asc">Name A → Z</option>
                <option value="desc">Name Z → A</option>
              </select>
            </div>
            <div class="flex-1 space-y-2">
              <${DurationRangeSlider}
                min=${durationMin}
                max=${durationMax}
                lowerValue=${filterMin}
                upperValue=${filterMax}
                step=${60}
                onMinChange=${this.handleMinDurationChange}
                onMaxChange=${this.handleMaxDurationChange}
              />
            </div>
          </div>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          ${workouts.map(workout => html`
            <${WorkoutCard}
              workout=${workout}
              maxPower=${maxPower}
              onClick=${() => this.selectWorkout(workout)}
            />
          `)}
        </div>
      </div>
    `;
  }
}

render(html`<${App} />`, document.getElementById('app'));
