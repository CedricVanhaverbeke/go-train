
const { h, Component } = window;
const html = window.htm.bind(h);

class VisualWorkoutEditor extends Component {
  constructor(props) {
    super(props);
    this.state = {
      steps: props.steps || [],
      editingStep: null,
      draggingIndex: null,
      dropIndex: null,
    };
    this.draggedStepNode = null;
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.steps !== this.state.steps) {
      this.setState({ steps: nextProps.steps });
    }
  }

  handleStepChange = (index, field, value) => {
    const { steps } = this.state;
    const newSteps = [...steps];
    newSteps[index][field] = Math.max(0, Number(value));
    this.setState({ steps: newSteps });
    this.props.onStepsChange(newSteps);
  };

  addStep = () => {
    const newSteps = [...this.state.steps, { duration: 60, power: 100 }];
    this.setState({ steps: newSteps });
    this.props.onStepsChange(newSteps);
  };

  removeStep = (index) => {
    const newSteps = this.state.steps.filter((_, i) => i !== index);
    this.setState({ steps: newSteps, editingStep: null });
    this.props.onStepsChange(newSteps);
  };

  handleStepClick = (e, index) => {
    e.stopPropagation();
    this.setState({ editingStep: this.state.editingStep === index ? null : index });
  };

  // --- Drag and Drop Handlers ---
  handleDragStart = (e, index) => {
    this.setState({ draggingIndex: index, editingStep: null });
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/html', e.currentTarget);
    this.draggedStepNode = e.currentTarget;
  };

  handleDragOver = (e, index) => {
    e.preventDefault();
    if (this.state.draggingIndex === null || this.state.draggingIndex === index) {
      this.setState({ dropIndex: null });
      return;
    }
    this.setState({ dropIndex: index });
  };

  handleDrop = (e) => {
    e.preventDefault();
    const { steps, draggingIndex, dropIndex } = this.state;
    if (draggingIndex === null || dropIndex === null || draggingIndex === dropIndex) {
      return;
    }

    const newSteps = [...steps];
    const [draggedStep] = newSteps.splice(draggingIndex, 1);
    newSteps.splice(dropIndex, 0, draggedStep);

    this.setState({ steps: newSteps, draggingIndex: null, dropIndex: null });
    this.props.onStepsChange(newSteps);
  };

  handleDragEnd = () => {
    this.draggedStepNode.style.opacity = '1';
    this.setState({ draggingIndex: null, dropIndex: null });
  };

  // --- Resize Handlers ---
  handleDurationResize = (e, index) => {
    e.stopPropagation();
    const startX = e.clientX;
    const startDuration = this.state.steps[index].duration;
    const onMouseMove = (moveEvent) => {
      const newDuration = startDuration + (moveEvent.clientX - startX);
      this.handleStepChange(index, 'duration', newDuration);
    };
    const onMouseUp = () => {
      document.removeEventListener('mousemove', onMouseMove);
      document.removeEventListener('mouseup', onMouseUp);
    };
    document.addEventListener('mousemove', onMouseMove);
    document.addEventListener('mouseup', onMouseUp);
  };

  handlePowerResize = (e, index) => {
    e.stopPropagation();
    const startY = e.clientY;
    const startPower = this.state.steps[index].power;
    const containerHeight = e.currentTarget.closest('.h-64').offsetHeight;
    const maxPower = this.state.steps.reduce((max, step) => Math.max(max, step.power), 100);

    const onMouseMove = (moveEvent) => {
      const deltaY = startY - moveEvent.clientY;
      const powerChange = (deltaY / containerHeight) * maxPower;
      this.handleStepChange(index, 'power', startPower + powerChange);
    };
    const onMouseUp = () => {
      document.removeEventListener('mousemove', onMouseMove);
      document.removeEventListener('mouseup', onMouseUp);
    };
    document.addEventListener('mousemove', onMouseMove);
    document.addEventListener('mouseup', onMouseUp);
  };

  render() {
    const { steps, editingStep, draggingIndex, dropIndex } = this.state;
    const totalDuration = steps.reduce((sum, step) => sum + step.duration, 0);
    const maxPower = Math.max(100, ...steps.map(s => s.power));

    return html`
      <div class="space-y-4" onClick=${() => this.setState({ editingStep: null })}>
        <h3 class="text-lg font-semibold text-slate-300">Visual Workout Editor</h3>
        <div class="relative bg-slate-700 h-64 w-full rounded-lg overflow-x-auto flex" onDragOver=${e => e.preventDefault()} onDrop=${this.handleDrop}>
            ${steps.map((step, index) => {
              const isDragging = draggingIndex === index;
              const isDropTarget = dropIndex === index;
              return html`
                <div
                  key=${index}
                  draggable="true"
                  onDragStart=${(e) => this.handleDragStart(e, index)}
                  onDragOver=${(e) => this.handleDragOver(e, index)}
                  onDragEnd=${this.handleDragEnd}
                  onClick=${(e) => this.handleStepClick(e, index)}
                  class="relative h-full group transition-all duration-150 ${isDragging ? 'opacity-50' : 'opacity-100'} ${isDropTarget ? 'bg-green-500/20' : ''}"
                  style=${{ width: `${step.duration}px`, cursor: 'grab' }}
                >
                  <div
                    class="absolute bottom-0 w-full bg-sky-500 group-hover:bg-sky-400 transition-colors"
                    style=${{ height: `${(step.power / maxPower) * 100}%`, pointerEvents: 'none' }}
                  ></div>

                  ${editingStep === index && html`
                    <div
                      onMouseDown=${(e) => this.handleDurationResize(e, index)}
                      class="absolute top-1/2 right-0 w-2 h-4 bg-white rounded-sm cursor-ew-resize z-20"
                      style=${{ transform: 'translate(50%, -50%)' }}
                    ></div>
                    <div
                      onMouseDown=${(e) => this.handlePowerResize(e, index)}
                      class="absolute left-1/2 w-4 h-2 bg-white rounded-sm cursor-ns-resize z-20"
                      style=${{ top: `${100 - (step.power / maxPower) * 100}%`, transform: 'translate(-50%, -50%)' }}
                    ></div>

                    <div class="absolute top-0 left-0 bg-slate-800 p-2 rounded-lg shadow-lg z-30 w-48" onClick=${e => e.stopPropagation()}>
                      <p class="text-xs font-bold text-slate-200 mb-2">Edit Step ${index + 1}</p>
                      <div class="space-y-2">
                        <div>
                          <label class="block text-xs font-semibold text-slate-400 mb-1">Duration (s)</label>
                          <input
                            type="number"
                            value=${step.duration}
                            onInput=${(e) => this.handleStepChange(index, 'duration', e.target.value)}
                            class="w-full rounded-md bg-slate-900 border border-slate-600 px-2 py-1"
                          />
                        </div>
                        <div>
                          <label class="block text-xs font-semibold text-slate-400 mb-1">Power (%FTP)</label>
                          <input
                            type="number"
                            value=${step.power}
                            onInput=${(e) => this.handleStepChange(index, 'power', e.target.value)}
                            class="w-full rounded-md bg-slate-900 border border-slate-600 px-2 py-1"
                          />
                        </div>
                      </div>
                      <button onClick=${() => this.removeStep(index)} class="mt-3 text-xs text-red-400 hover:text-red-300">Remove Step</button>
                    </div>
                  `}
                </div>
              `;
            })}
        </div>
        <div class="flex justify-between items-center">
            <div class="flex gap-2">
                <button onClick=${this.addStep} class="inline-flex items-center gap-2 rounded-lg bg-slate-600 hover:bg-slate-500 transition-colors px-4 py-2 font-semibold text-slate-100">
                    Add Step
                </button>
            </div>
            <div class="text-sm text-slate-400">
                Total Duration: ${formatDuration(totalDuration)}
            </div>
        </div>
      </div>
    `;
  }
}

// Helper function to format duration
function formatDuration(seconds) {
  if (!Number.isFinite(seconds) || seconds < 0) return '0s';
  const min = Math.floor(seconds / 60);
  const sec = Math.round(seconds % 60);
  return min > 0 ? `${min}m ${sec}s` : `${sec}s`;
}

window.VisualWorkoutEditor = VisualWorkoutEditor;
