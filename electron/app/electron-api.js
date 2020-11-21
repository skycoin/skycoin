// Warning: for this to work the contextIsolation property of the browser window must be set to false.

window.isElectron = true;
window.ipcRenderer = require('electron').ipcRenderer;