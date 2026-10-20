import { firstValueFrom, of, throwError } from 'rxjs';

import { ConnectionSources, NetworkService } from './network.service';

describe('NetworkService', () => {
  function service(answers: { [url: string]: any }): NetworkService {
    const api = {
      get: (url: string) => answers[url] instanceof Error ? throwError(() => answers[url]) : of(answers[url]),
    };

    return new NetworkService(api as any);
  }

  it('lists the peers sorted, marking the default ones', async () => {
    const peers = await firstValueFrom(service({
      'network/defaultConnections': ['2.2.2.2:6000'],
      'network/connections': { connections: [
        { address: '3.3.3.3:6000', listen_port: 6000, outgoing: false, height: 9, last_sent: 1, last_received: 2 },
        { address: '2.2.2.2:6000', listen_port: 6000, outgoing: true, height: 8, last_sent: 3, last_received: 4 },
      ] },
    }).connections());

    expect(peers!.map(p => p.address)).toEqual(['2.2.2.2:6000', '3.3.3.3:6000']);
    expect(peers!.map(p => p.source)).toEqual([ConnectionSources.Default, ConnectionSources.Exchange]);
    expect(peers![0].outgoing).toBeTrue();
  });

  it('emits null when the node does not answer', async () => {
    const peers = await firstValueFrom(service({
      'network/defaultConnections': new Error('down'),
      'network/connections': { connections: [] },
    }).connections());

    expect(peers).toBeNull();
  });
});
