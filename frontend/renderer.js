const { ipcRenderer } = require("electron");

const { h, Component, render } = window;
const html = window.htm.bind(h);

// --- Data Fetching ---
async function getWorkouts() {
  const response = await fetch("./data.json");
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
  if (!Number.isFinite(seconds) || seconds < 0) return "0 min";
  const minutes = Math.round(seconds / 60);
  return minutes >= 60
    ? `${(minutes / 60).toFixed(minutes % 60 === 0 ? 0 : 1)} hr`
    : `${minutes} min`;
}

function scaleWorkoutDuration(workout, targetDurationSeconds) {
  if (!workout) return null;
  const baseTotal = workout.totalDuration ?? getTotalDuration(workout);
  if (!baseTotal) {
    return {
      ...workout,
      baseTotalDuration: 0,
      durationScale: 1,
      steps: workout.steps.map((step) => ({ ...step })),
    };
  }

  const rawTarget = Number(targetDurationSeconds);
  const safeTarget = Number.isFinite(rawTarget) ? rawTarget : baseTotal;
  const minPossible = workout.steps.length || 1;
  const target = Math.max(Math.round(safeTarget), minPossible);
  const scale = target / baseTotal;

  const rawDurations = workout.steps.map((step) => step.duration * scale);
  const flooredDurations = rawDurations.map((raw) =>
    Math.max(1, Math.floor(raw)),
  );
  let sum = flooredDurations.reduce((acc, duration) => acc + duration, 0);
  let remaining = Math.max(target - sum, 0);

  if (remaining > 0) {
    const remainderOrder = rawDurations
      .map((raw, index) => ({ index, remainder: raw - Math.floor(raw) }))
      .sort((a, b) => b.remainder - a.remainder);

    for (let i = 0; i < remainderOrder.length && remaining > 0; i += 1) {
      flooredDurations[remainderOrder[i].index] += 1;
      remaining -= 1;
    }
  }

  const scaledSteps = workout.steps.map((step, index) => ({
    ...step,
    duration: flooredDurations[index],
  }));
  const scaledTotalDuration = scaledSteps.reduce(
    (sumDuration, step) => sumDuration + step.duration,
    0,
  );

  return {
    ...workout,
    steps: scaledSteps,
    totalDuration: scaledTotalDuration,
    targetDuration: target,
    baseTotalDuration: baseTotal,
    durationScale: scaledTotalDuration / baseTotal,
  };
}

// --- Components ---

function WorkoutPreview({ workout, maxPower, width = 200, height = 100 }) {
  const totalDuration = workout.totalDuration ?? getTotalDuration(workout);
  const barWidthScale = width / totalDuration; // Scale bars to the provided width
  const formattedDuration = formatDuration(totalDuration);

  let accumulatedDuration = 0;
  return html`
    <div class="relative inline-flex">
      <svg width=${width} height=${height} class="bg-slate-700 rounded">
        ${workout.steps.map((step) => {
          const barWidth = step.duration * barWidthScale;
          const barHeight = (step.power / maxPower) * height;
          const x = accumulatedDuration * barWidthScale;
          accumulatedDuration += step.duration;
          return html`
            <rect
              x=${x}
              y=${height - barHeight}
              width=${barWidth}
              height=${barHeight}
              fill=${powerToColor(step.power, maxPower)}
            />
          `;
        })}
      </svg>
      <span
        class="absolute bottom-2 right-2 rounded bg-slate-900/80 px-2 py-1 text-xs font-semibold text-slate-100"
      >
        ${formattedDuration}
      </span>
    </div>
  `;
}

function WorkoutCard({ workout, maxPower, onClick }) {
  return html`
    <div
      class="bg-slate-800 rounded-lg shadow-lg p-4 flex flex-col gap-4 cursor-pointer hover:bg-slate-700 transition-colors"
      onClick=${onClick}
    >
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
      <div
        class="flex items-center justify-between text-sm font-semibold text-slate-300"
      >
        <span>Duration filter</span>
        <span
          >${formatDuration(clampedLower)} –
          ${formatDuration(clampedUpper)}</span
        >
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

function DurationAdjuster({
  baseDuration,
  currentDuration,
  onChange,
  stepsCount = 1,
}) {
  if (!Number.isFinite(baseDuration) || baseDuration <= 0) {
    return null;
  }

  const safeBase = Math.max(baseDuration, 1);
  const safeCurrent = Math.max(Number(currentDuration) || safeBase, 1);
  const sliderLowerBound = Math.max(stepsCount, 10);
  const sliderMin = Math.max(
    sliderLowerBound,
    Math.min(Math.round(safeBase * 0.5), Math.floor(safeCurrent * 0.8)),
  );
  const sliderMax = Math.max(
    Math.round(safeBase * 1.5),
    Math.ceil(safeCurrent * 1.2),
    sliderMin + 60,
  );
  const sliderRange = Math.max(sliderMax - sliderMin, 1);
  const clampedCurrent = Math.min(Math.max(safeCurrent, sliderMin), sliderMax);
  const sliderPercent = ((clampedCurrent - sliderMin) / sliderRange) * 100;
  const scalePercent = Math.round((safeCurrent / safeBase) * 100);

  const handleSliderChange = (event) => {
    onChange?.(Number(event.target.value));
  };

  const handleMinutesChange = (event) => {
    const minutes = Number(event.target.value);
    if (!Number.isFinite(minutes)) return;
    const secondsValue = minutes * 60;
    onChange?.(Math.max(sliderMin, secondsValue));
  };

  return html`
    <div class="space-y-4">
      <div
        class="flex items-center justify-between text-sm font-semibold text-slate-200"
      >
        <span>Adjust duration</span>
        <span
          >${formatDuration(safeCurrent)} (${scalePercent}% of original)</span
        >
      </div>
      <div class="relative h-10 flex items-center">
        <div class="absolute inset-x-0 h-1 bg-slate-700 rounded-full"></div>
        <div
          class="absolute h-1 bg-sky-400 rounded-full"
          style=${{ left: "0%", width: `${sliderPercent}%` }}
        ></div>
        <input
          type="range"
          min=${sliderMin}
          max=${sliderMax}
          step=${30}
          value=${clampedCurrent}
          onInput=${handleSliderChange}
          class="range-thumb absolute inset-0 w-full h-full focus:outline-none"
        />
      </div>
      <div class="flex items-center gap-3 text-sm text-slate-300">
        <label
          class="text-xs font-semibold uppercase tracking-wide text-slate-400"
          >Minutes</label
        >
        <input
          type="number"
          min=${sliderMin / 60}
          step=${0.5}
          value=${Number((safeCurrent / 60).toFixed(2))}
          onInput=${handleMinutesChange}
          class="w-20 rounded-lg bg-slate-900 border border-slate-700 px-2 py-1 text-right"
        />
        <span class="text-xs text-slate-500">min</span>
        <button
          type="button"
          onClick=${() => onChange?.(safeBase)}
          class="ml-auto text-xs font-semibold text-sky-400 hover:text-sky-300"
        >
          Reset
        </button>
      </div>
      <p class="text-xs text-slate-500">
        Steps are scaled linearly so the workout profile stays intact.
      </p>
    </div>
  `;
}

function WorkoutDetail({
  workout,
  onBack,
  maxPower: globalMaxPower,
  desiredDuration,
  onDurationChange,
}) {
  const totalDuration = workout.steps.reduce(
    (sum, step) => sum + step.duration,
    0,
  );
  const workoutMaxPower = workout.steps.reduce(
    (max, step) => Math.max(max, step.power),
    0,
  );
  const previewMaxPower = globalMaxPower || workoutMaxPower;
  const baseDuration = workout.baseTotalDuration ?? totalDuration;
  const durationScalePercent = baseDuration
    ? Math.round((totalDuration / baseDuration) * 100)
    : 100;

  return html`
    <div class="space-y-6">
      <button onClick=${onBack} class="text-sky-400 hover:underline">
        ← Back to Workouts
      </button>
      <h2 class="text-3xl font-bold">${workout.name}</h2>
      <div
        class="w-full max-w-3xl rounded-2xl bg-slate-800 shadow-xl border border-slate-700 p-8 space-y-6"
      >
        <section class="space-y-3">
          <p class="text-sm font-semibold text-slate-200">Workout Preview</p>
          <div
            class="bg-slate-900 border border-slate-700 rounded-xl p-4 flex justify-center"
          >
            <${WorkoutPreview}
              workout=${workout}
              maxPower=${previewMaxPower}
              width=${480}
              height=${160}
            />
          </div>
        </section>
        <section class="space-y-3">
          <p class="text-sm font-semibold text-slate-200">Duration</p>
          <${DurationAdjuster}
            baseDuration=${baseDuration}
            currentDuration=${desiredDuration ??
            workout.targetDuration ??
            totalDuration}
            stepsCount=${workout.steps.length}
            onChange=${onDurationChange}
          />
          <div class="text-xs text-slate-400 flex flex-wrap gap-x-4 gap-y-1">
            <span>Original: ${formatDuration(baseDuration)}</span>
            <span>Current: ${formatDuration(totalDuration)}</span>
            <span>Scale: ${durationScalePercent}%</span>
          </div>
        </section>
        <section class="space-y-3">
          <div class="flex items-center gap-4">
            <button
              id="start-app"
              class="inline-flex items-center gap-2 rounded-lg bg-sky-500 hover:bg-sky-400 transition-colors px-6 py-3 font-semibold text-slate-900"
            >
              <span>Start App</span>
            </button>
            <button
              id="stop-app"
              class="inline-flex items-center gap-2 rounded-lg bg-red-500 hover:bg-red-400 transition-colors px-6 py-3 font-semibold text-slate-900 disabled:bg-slate-600 disabled:cursor-not-allowed"
              disabled
            >
              <span>Stop App</span>
            </button>
          </div>
          <div>
            <p class="text-sm font-semibold text-slate-200 mb-2">Status</p>
            <pre
              id="status-box"
              class="bg-slate-900 rounded-lg border border-slate-700 p-4 text-green-400 text-sm overflow-auto min-h-[4rem]"
            >
Waiting for user action...</pre
            >
          </div>
        </section>
      </div>
    </div>
  `;
}

class WorkoutCreator extends Component {
  state = {
    name: "",
    steps: [{ duration: 60, power: 100 }],
  };

  handleNameChange = (e) => {
    this.setState({ name: e.target.value });
  };

  handleStepsChange = (steps) => {
    this.setState({ steps });
  };

  handleSave = () => {
    const { name, steps } = this.state;
    if (name.trim() === "" || steps.length === 0) {
      alert("Please provide a name and at least one step for the workout.");
      return;
    }
    const newWorkout = {
      name: name.trim(),
      steps,
    };
    this.props.onSave(newWorkout);
  };

  render(_, { name, steps }) {
    const { VisualWorkoutEditor } = window;
    return html`
      <div class="space-y-6">
        <button
          onClick=${this.props.onBack}
          class="text-sky-400 hover:underline"
        >
          ← Back to Workouts
        </button>
        <h2 class="text-3xl font-bold">Create Workout</h2>
        <div
          class="w-full max-w-3xl rounded-2xl bg-slate-800 shadow-xl border border-slate-700 p-8 space-y-6"
        >
          <div class="space-y-2">
            <label class="block text-sm font-semibold text-slate-300"
              >Workout Name</label
            >
            <input
              type="text"
              value=${name}
              onInput=${this.handleNameChange}
              class="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-sky-400"
              placeholder="e.g., Morning Spin"
            />
          </div>

          <${VisualWorkoutEditor}
            steps=${steps}
            onStepsChange=${this.handleStepsChange}
          />

          <div class="flex items-center gap-4">
            <button
              onClick=${this.handleSave}
              class="inline-flex items-center gap-2 rounded-lg bg-sky-500 hover:bg-sky-400 transition-colors px-6 py-3 font-semibold text-slate-900"
            >
              Save Workout
            </button>
            <button
              onClick=${this.props.onBack}
              class="text-slate-400 hover:text-slate-300"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    `;
  }
}

class App extends Component {
  state = {
    workouts: [],
    selectedWorkout: null,
    selectedWorkoutTargetDuration: null,
    maxPower: 0,
    sortOrder: "default",
    durationMin: 0,
    durationMax: 0,
    filterMin: 0,
    filterMax: 0,
    view: "list",
  };

  showCreateWorkoutView = () => this.setState({ view: "create" });

  showListView = () => this.setState({ view: "list" });

  handleSaveWorkout = (newWorkout) => {
    const workoutWithDuration = {
      ...newWorkout,
      totalDuration: getTotalDuration(newWorkout),
    };

    this.setState((prevState) => ({
      workouts: [...prevState.workouts, workoutWithDuration],
      view: "list",
    }));
  };

  componentDidMount() {
    getWorkouts().then((workouts) => {
      const workoutsWithDuration = workouts.map((workout) => ({
        ...workout,
        totalDuration: getTotalDuration(workout),
      }));
      const powerValues = workoutsWithDuration.flatMap((w) =>
        w.steps.map((s) => s.power),
      );
      const maxPower = powerValues.length ? Math.max(...powerValues) : 0;
      const durationValues = workoutsWithDuration.map((w) => w.totalDuration);
      const durationMin = durationValues.length
        ? Math.min(...durationValues)
        : 0;
      const durationMax = durationValues.length
        ? Math.max(...durationValues)
        : 0;

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

  selectWorkout = (workout) => {
    const totalDuration = workout.totalDuration ?? getTotalDuration(workout);
    const workoutString = workout.steps
      .map(({ power, duration }) => `${power}-${duration}`)
      .join(";");
    this.setState(
      {
        selectedWorkout: workout,
        selectedWorkoutTargetDuration: totalDuration,
      },
      () => {
        // we need to re-add the event listeners after the DOM is updated
        const startButton = document.getElementById("start-app");
        const stopButton = document.getElementById("stop-app");
        const statusBox = document.getElementById("status-box");

        const updateStatus = (message) => {
          const timestamp = new Date().toLocaleTimeString();
          if (statusBox) {
            statusBox.textContent = `[${timestamp}] ${message}`;
          }
        };

        if (startButton) {
          startButton.addEventListener("click", () => {
            const appName = "overlay";
            updateStatus(`Renderer: requesting ${appName}...`);
            ipcRenderer.send("START_APP", appName, "-workout", workoutString);
            startButton.disabled = true;
            stopButton.disabled = false;
          });
        }

        if (stopButton) {
          stopButton.addEventListener("click", () => {
            const appName = "overlay";
            updateStatus(`Renderer: stopping ${appName}...`);
            ipcRenderer.send("STOP_APP", appName);
            startButton.disabled = false;
            stopButton.disabled = true;
          });
        }

        ipcRenderer.on("APP_STATUS", (_event, message) => {
          updateStatus(`Main Process: ${message}`);
        });

        ipcRenderer.on("APP_OUTPUT", (_event, data) => {
          console.log(data);
        });

        ipcRenderer.on("APP_EXITED", () => {
          if (startButton) startButton.disabled = false;
          if (stopButton) stopButton.disabled = true;
        });
      },
    );
  };

  unselectWorkout = () =>
    this.setState({
      selectedWorkout: null,
      selectedWorkoutTargetDuration: null,
    });

  handleSortChange = (event) => {
    this.setState({ sortOrder: event.target.value });
  };

  handleMinDurationChange = (event) => {
    const newMin = Number(event.target.value);
    this.setState((prev) => ({ filterMin: Math.min(newMin, prev.filterMax) }));
  };

  handleMaxDurationChange = (event) => {
    const newMax = Number(event.target.value);
    this.setState((prev) => ({ filterMax: Math.max(newMax, prev.filterMin) }));
  };

  handleSelectedWorkoutDurationChange = (newDuration) => {
    const parsed = Number(newDuration);
    if (!Number.isFinite(parsed)) return;
    this.setState({ selectedWorkoutTargetDuration: parsed });
  };

  getFilteredWorkouts() {
    const { workouts, sortOrder, filterMin, filterMax } = this.state;
    const filtered = workouts.filter((workout) => {
      const duration = workout.totalDuration ?? getTotalDuration(workout);
      return duration >= filterMin && duration <= filterMax;
    });

    if (sortOrder === "asc") {
      return [...filtered].sort((a, b) => a.name.localeCompare(b.name));
    }

    if (sortOrder === "desc") {
      return [...filtered].sort((a, b) => b.name.localeCompare(a.name));
    }

    return filtered;
  }

  render(
    _,
    {
      selectedWorkout,
      selectedWorkoutTargetDuration,
      maxPower,
      sortOrder,
      durationMin,
      durationMax,
      filterMin,
      filterMax,
      view,
    },
  ) {
    const workouts = this.getFilteredWorkouts();

    if (selectedWorkout) {
      const scaledWorkout = scaleWorkoutDuration(
        selectedWorkout,
        selectedWorkoutTargetDuration,
      );
      return html`<${WorkoutDetail}
        workout=${scaledWorkout}
        maxPower=${maxPower}
        onBack=${this.unselectWorkout}
        desiredDuration=${selectedWorkoutTargetDuration ??
        scaledWorkout?.baseTotalDuration}
        onDurationChange=${this.handleSelectedWorkoutDurationChange}
      />`;
    }

    if (view === "create") {
      return html`<${WorkoutCreator}
        onSave=${this.handleSaveWorkout}
        onBack=${this.showListView}
      />`;
    }

    return html`
      <div class="space-y-6">
        <div class="flex items-center justify-between">
          <h1 class="text-3xl font-bold">Workouts</h1>
          <button
            onClick=${this.showCreateWorkoutView}
            class="inline-flex items-center gap-2 rounded-lg bg-sky-500 hover:bg-sky-400 transition-colors px-4 py-2 font-semibold text-slate-900"
          >
            Create Workout
          </button>
        </div>
        <div
          class="bg-slate-800 border border-slate-700 rounded-xl p-4 space-y-4"
        >
          <div class="flex flex-col gap-4 md:flex-row">
            <div class="flex-1">
              <label class="block text-sm font-semibold text-slate-300 mb-2"
                >Sort</label
              >
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
          ${workouts.map(
            (workout) => html`
              <${WorkoutCard}
                workout=${workout}
                maxPower=${maxPower}
                onClick=${() => this.selectWorkout(workout)}
              />
            `,
          )}
        </div>
      </div>
    `;
  }
}

render(html`<${App} />`, document.getElementById("app"));
