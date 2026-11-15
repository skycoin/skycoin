import { fakeAsync, tick } from '@angular/core/testing';

import { GlobalsService } from './globals.service';

describe('GlobalsService', () => {
  const translate: any = { instant: (key: string) => key };

  it('gives the node version once it is set', fakeAsync(() => {
    const service = new GlobalsService(translate);
    let version = '';
    service.getValidNodeVersion().subscribe(v => version = v);
    tick(1000);
    service.setNodeVersion('0.28.7');
    expect(version).toBe('0.28.7');
  }));

  it('fails instead of waiting forever when the node was never checked', fakeAsync(() => {
    const service = new GlobalsService(translate);
    let error: Error | undefined;
    service.getValidNodeVersion().subscribe({ error: e => error = e });
    tick(service.nodeVersionWaitMs - 1);
    expect(error).toBeUndefined();
    tick(1);
    expect(error!.message).toBe('service.api.node-not-ready');
  }));
});
