import { environment } from '../environments/environment';

// The wallet registers its own service worker only as a top-level page that no
// other worker controls. Inside the skywire desk, the desk's worker owns requests.
export function useOwnServiceWorker(): boolean {
  if (!environment.production || window.self !== window.top || !('serviceWorker' in navigator)) {
    return false;
  }
  const controller = navigator.serviceWorker.controller;

  return !controller || controller.scriptURL.endsWith('/ngsw-worker.js');
}
