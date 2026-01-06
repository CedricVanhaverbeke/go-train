const { app, BrowserWindow, ipcMain } = require("electron");
const path = require("path");
const fs = require("fs");
const os = require("os");
const { spawn } = require("child_process");
const sqlite3 = require("sqlite3");

// Keep a global reference to the child process.
let child = null;

// Define the commands that the renderer is allowed to start.
const AVAILABLE_APPS = {
  overlay: () => {
    const basePath = app.isPackaged
      ? path.join(process.resourcesPath, "bin")
      : path.join(__dirname, "..", "bin");

    // Path to a compiled Go server binary (replace with your real path).
    const binary = path.join(basePath, `overlay-${os.platform()}-${os.arch()}`);

    if (!fs.existsSync(binary)) {
      throw new Error(`Go server binary not found at ${binary}`);
    }

    return {
      command: binary,
      args: [],
      options: { cwd: path.dirname(binary) },
    };
  },
};

function createWindow() {
  const mainWindow = new BrowserWindow({
    width: 960,
    height: 640,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false,
    },
  });

  mainWindow.webContents.openDevTools();

  mainWindow.loadFile("index.html");
}

function startApp(event, appName, ...args) {
  try {
    const builder = AVAILABLE_APPS[appName];

    if (!builder) {
      throw new Error(
        `The selected application is not registered in main.js: ${appName}`,
      );
    }

    const { command, args: builderArgs, options = {} } = builder();

    child = spawn(command, [...args, ...builderArgs], {
      ...options,
    });

    child.on("spawn", () => {
      event.reply(
        "APP_STATUS",
        `Successfully started ${appName} (PID: ${child.pid}).`,
      );
    });

    child.stdout.on("data", (data) => {
      console.log(data.toString());
      event.reply("APP_STDOUT", data.toString());
    });

    child.stderr.on("data", (data) => {
      console.log(data.toString());
      event.reply("APP_STDERR", data.toString());
    });

    child.on("error", (error) => {
      event.reply("APP_STATUS", `Failed to start ${appName}: ${error.message}`);
    });

    child.on("exit", (code) => {
      event.reply("APP_STATUS", `App exited with code: ${code}.`);
      event.reply("APP_EXITED");
      child = null;
    });
  } catch (error) {
    event.reply("APP_STATUS", `Unable to start ${appName}: ${error.message}`);
  }
}

function stopApp(event, appName) {
  if (!child) {
    event.reply("APP_STATUS", "No application is currently running.");
    return;
  }

  child.kill();
  child = null;

  event.reply("APP_STATUS", `Successfully stopped ${appName}.`);
}

function getGpxFiles(event) {
  const dbPath = path.join(__dirname, "..", "db");
  const db = new sqlite3.Database(dbPath);

  db.all(
    "SELECT id, name, created_at FROM gpx_files ORDER BY created_at DESC",
    [],
    (err, rows) => {
      if (err) {
        event.reply("GPX_FILES_ERROR", err.message);
      } else {
        console.log(rows);
        event.reply("GPX_FILES_DATA", rows);
      }
      db.close();
    },
  );
}

function getGpxFileData(event, id) {
  const dbPath = path.join(__dirname, "..", "db");
  const db = new sqlite3.Database(dbPath);

  db.get("SELECT name, data FROM gpx_files WHERE id = ?", [id], (err, row) => {
    if (err) {
      event.reply("GPX_FILE_ERROR", err.message);
    } else {
      event.reply("GPX_FILE_DATA", row);
    }
    db.close();
  });
}

app.whenReady().then(() => {
  createWindow();

  ipcMain.on("START_APP", (event, ...appName) => {
    startApp(event, ...appName);
  });

  ipcMain.on("STOP_APP", (event, appName) => {
    stopApp(event, appName);
  });

  ipcMain.on("GET_GPX_FILES", (event) => {
    getGpxFiles(event);
  });

  ipcMain.on("GET_GPX_FILE_DATA", (event, id) => {
    getGpxFileData(event, id);
  });

  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on("before-quit", () => {
  if (child) {
    child.kill();
    child = null;
  }
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});
