const { h, Component } = window;
const html = window.htm.bind(h);

class TrainingList extends Component {
  constructor(props) {
    super(props);
    this.state = {
      trainings: [],
      loading: false,
      error: null,
    };
  }

  componentDidMount() {
    this.loadTrainings();
  }

  loadTrainings = () => {
    this.setState({ loading: true, error: null });
    
    // Request GPX files from main process
    window.ipcRenderer.send('GET_GPX_FILES');
    
    // Set up listeners for the response
    const handleFilesData = (event, trainings) => {
      this.setState({ 
        trainings: trainings || [], 
        loading: false 
      });
      window.ipcRenderer.removeListener('GPX_FILES_DATA', handleFilesData);
      window.ipcRenderer.removeListener('GPX_FILES_ERROR', handleFilesError);
    };
    
    const handleFilesError = (event, error) => {
      this.setState({ 
        error: error, 
        loading: false 
      });
      window.ipcRenderer.removeListener('GPX_FILES_DATA', handleFilesData);
      window.ipcRenderer.removeListener('GPX_FILES_ERROR', handleFilesError);
    };
    
    window.ipcRenderer.on('GPX_FILES_DATA', handleFilesData);
    window.ipcRenderer.on('GPX_FILES_ERROR', handleFilesError);
  };

  downloadTraining = (id, name) => {
    // Request GPX file data from main process
    window.ipcRenderer.send('GET_GPX_FILE_DATA', id);
    
    // Set up listeners for the response
    const handleFileData = (event, fileData) => {
      if (fileData && fileData.data) {
        // Create blob and download link
        const blob = new Blob([fileData.data], { type: 'application/gpx+xml' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${name || 'training'}.gpx`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
      }
      window.ipcRenderer.removeListener('GPX_FILE_DATA', handleFileData);
      window.ipcRenderer.removeListener('GPX_FILE_ERROR', handleFileError);
    };
    
    const handleFileError = (event, error) => {
      console.error('Error downloading file:', error);
      this.setState({ error: `Failed to download ${name}: ${error}` });
      window.ipcRenderer.removeListener('GPX_FILE_DATA', handleFileData);
      window.ipcRenderer.removeListener('GPX_FILE_ERROR', handleFileError);
    };
    
    window.ipcRenderer.on('GPX_FILE_DATA', handleFileData);
    window.ipcRenderer.on('GPX_FILE_ERROR', handleFileError);
  };

  formatDate = (dateString) => {
    if (!dateString) return 'Unknown date';
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
    } catch (error) {
      return 'Invalid date';
    }
  };

  render() {
    const { trainings, loading, error } = this.state;

    return html`
      <div class="space-y-4">
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-semibold text-slate-300">
            Training Files
          </h3>
          <button
            onClick=${this.loadTrainings}
            class="text-sm text-slate-400 hover:text-slate-200"
          >
            Refresh
          </button>
        </div>

        ${error && html`
          <div class="bg-red-500/20 border border-red-500/50 text-red-300 px-4 py-3 rounded-lg">
            <p class="text-sm">${error}</p>
          </div>
        `}

        ${loading && html`
          <div class="text-center py-8">
            <p class="text-slate-400">Loading training files...</p>
          </div>
        `}

        ${!loading && !error && trainings.length === 0 && html`
          <div class="text-center py-8">
            <p class="text-slate-400">No training files found</p>
          </div>
        `}

        ${!loading && !error && trainings.length > 0 && html`
          <div class="space-y-2">
            ${trainings.map(training => html`
              <div class="bg-slate-700 rounded-lg p-4 flex items-center justify-between">
                <div class="flex-1">
                  <h4 class="font-medium text-slate-100">${training.name || 'Untitled Training'}</h4>
                  <p class="text-sm text-slate-400">
                    ${this.formatDate(training.created_at)}
                  </p>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-xs text-slate-500">ID: ${training.id}</span>
                  <button
                    onClick=${() => this.downloadTraining(training.id, training.name)}
                    class="inline-flex items-center gap-2 rounded-lg bg-sky-600 hover:bg-sky-500 transition-colors px-4 py-2 font-semibold text-white text-sm"
                  >
                    Download
                  </button>
                </div>
              </div>
            `)}
          </div>
        `}

        <div class="text-sm text-slate-400 text-center">
          ${trainings.length} training file${trainings.length !== 1 ? 's' : ''} found
        </div>
      </div>
    `;
  }
}

window.TrainingList = TrainingList;