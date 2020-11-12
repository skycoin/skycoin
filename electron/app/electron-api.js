// Warning: for this to work the contextIsolation property of the browser window must be set to true.

window.isElectron = true;
window.ipcRenderer = require('electron').ipcRenderer;