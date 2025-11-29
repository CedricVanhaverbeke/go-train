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

// --- Components ---

function WorkoutPreview({ workout, maxPower }) {
  const totalDuration = workout.steps.reduce((sum, step) => sum + step.duration, 0);
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
  };

  componentDidMount() {
    getWorkouts().then(workouts => {
      const maxPower = Math.max(...workouts.flatMap(w => w.steps.map(s => s.power)));
      this.setState({ workouts, maxPower });
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

  render(_, { workouts, selectedWorkout, maxPower }) {
    if (selectedWorkout) {
      return html`<${WorkoutDetail} workout=${selectedWorkout} onBack=${this.unselectWorkout} />`;
    }

    return html`
      <div class="space-y-4">
        <h1 class="text-3xl font-bold">Workouts</h1>
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
