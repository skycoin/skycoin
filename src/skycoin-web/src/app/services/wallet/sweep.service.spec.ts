import { of } from 'rxjs';

import { SweepService, sweepKind, sweepActive } from './sweep.service';
import { ApiService } from '../api.service';

describe('SweepService', () => {
  let api: jasmine.SpyObj<ApiService>;
  let service: SweepService;
  const data = { rule: null, status: { wallet: 'w.wlt', enabled: false } };

  beforeEach(() => {
    api = jasmine.createSpyObj('ApiService', ['get', 'post']);
    api.get.and.returnValue(of({ data: data }));
    api.post.and.returnValue(of({ data: data }));
    service = new SweepService(api);
  });

  it('reads the rule and status from the v2 API', () => {
    let got: any;
    service.get('w.wlt').subscribe(s => got = s);
    expect(api.get).toHaveBeenCalledWith('wallet/sweep', { id: 'w.wlt' }, {}, true);
    expect(got).toEqual(data as any);
  });

  it('sets a rule', () => {
    const rule = { wallet: 'w.wlt', enabled: true, shared_accounts: [0], destination: { kind: 'account' as const } };
    service.setRule(rule).subscribe();
    expect(api.post).toHaveBeenCalledWith('wallet/sweep', rule, {}, true);
  });

  it('unlocks and locks', () => {
    service.unlock('w.wlt', 'pw', 0).subscribe();
    expect(api.post).toHaveBeenCalledWith('wallet/sweep/unlock', { id: 'w.wlt', password: 'pw', ttl_seconds: 0 }, {}, true);
    service.lock('w.wlt').subscribe();
    expect(api.post).toHaveBeenCalledWith('wallet/sweep/lock', { id: 'w.wlt' }, {}, true);
  });

  it('says who sweeps a wallet', () => {
    const coin: any = { isGateway: () => false };
    expect(sweepKind({ label: 'n', filename: 'n.wlt', addresses: [] }, true, coin)).toBe('node');
    expect(sweepKind({ label: 'b', addresses: [] }, false, coin)).toBe('browser');
    expect(sweepKind({ label: 'h', isHardware: true, addresses: [] }, false, coin)).toBe('hardware');
    expect(sweepKind({ label: 't', temporary: true, addresses: [] }, false, coin)).toBeNull();
    expect(sweepKind({ label: 'x', walletType: 'xpub', addresses: [] }, true, coin)).toBeNull();
    expect(sweepKind({ label: 'b', addresses: [] }, false, { isGateway: () => true } as any)).toBeNull();
  });

  it('says a rule with only mappings is active', () => {
    expect(sweepActive(null)).toBeFalse();
    expect(sweepActive({ wallet: 'w', enabled: false, destination: { kind: 'wallet' } })).toBeFalse();
    expect(sweepActive({ wallet: 'w', enabled: false, destination: { kind: 'wallet' }, mappings: [{ from: { address: 'a' }, to: { kind: 'wallet' } }] })).toBeTrue();
  });
});
