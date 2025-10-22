// Polyfill for Go WASM (requires 'global' to be defined)
if (typeof global === 'undefined') {
  window.global = window;
}

window.removeSplash = function() {
  var element = document.getElementById('splashScreen');
  element.parentNode.removeChild(element);
}

