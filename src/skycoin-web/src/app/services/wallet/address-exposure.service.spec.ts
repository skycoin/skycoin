import { firstValueFrom, of, throwError } from 'rxjs';

import { AddressExposureService } from './address-exposure.service';

describe('AddressExposureService', () => {
  let calls: string[];
  let uxouts: { [address: string]: any };

  function newService(): AddressExposureService {
    const api = {
      get: (_url: string, params: any) => {
        calls.push(params.address);
        const answer = uxouts[params.address];

        return answer instanceof Error ? throwError(() => answer) : of(answer);
      },
    };

    return new AddressExposureService(api as any);
  }

  beforeEach(() => {
    calls = [];
    uxouts = {
      spent: [{ spent_block_seq: 0 }, { spent_block_seq: 195400 }],
      unspent: [{ spent_block_seq: 0 }],
      empty: [],
      broken: new Error('node down'),
    };
  });

  it('reports only addresses with a spent output', async () => {
    const exposed = await firstValueFrom(newService().getExposed(['spent', 'unspent', 'empty', 'broken']));
    expect(Array.from(exposed)).toEqual(['spent']);
  });

  it('asks again only for addresses that have not spent', async () => {
    const service = newService();
    await firstValueFrom(service.getExposed(['spent', 'unspent']));
    calls = [];
    const exposed = await firstValueFrom(service.getExposed(['spent', 'unspent']));
    expect(Array.from(exposed)).toEqual(['spent']);
    // Both are cached: spent for good, unspent for a while.
    expect(calls).toEqual([]);
  });
});
