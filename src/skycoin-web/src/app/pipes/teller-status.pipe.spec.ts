import { TellerStatusPipe } from './teller-status.pipe';
import { TestBed, inject } from '@angular/core/testing';

describe('TellerStatusPipe', () => {
  let pipe;

  beforeEach(() => TestBed.configureTestingModule({
    providers: [TellerStatusPipe]
  }));

  beforeEach(inject([TellerStatusPipe], p => {
    pipe = p;
  }));

  it('create an instance', () => {
    expect(pipe).toBeTruthy();
  });

  it('should work with empty string', () => {
    expect(pipe.transform('')).toBe('teller.unknown');
  });

  it('should transform waiting deposit', () => {
    expect(pipe.transform('waiting_deposit')).toBe('teller.waiting-deposit');
  });

  it('should transform waiting send', () => {
    expect(pipe.transform('waiting_send')).toBe('teller.waiting-send');
  });

  it('should transform waiting confirm', () => {
    expect(pipe.transform('waiting_confirm')).toBe('teller.waiting-confirm');
  });

  it('should transform done', () => {
    expect(pipe.transform('done')).toBe('teller.done');
  });
});
