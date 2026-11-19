import { InjectionToken } from '@angular/core';

/**
 * False when a host app sends the node calls through its own authenticated
 * proxy, as the skywire hypervisor does, so the node's CSRF setting guards nothing.
 */
export const NODE_CSRF_GUARDS_ACCESS = new InjectionToken<boolean>('NODE_CSRF_GUARDS_ACCESS', {
  providedIn: 'root',
  factory: () => true,
});
