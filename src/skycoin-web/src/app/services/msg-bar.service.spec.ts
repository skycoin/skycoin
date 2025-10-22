import { TestBed, inject } from '@angular/core/testing';

import { MsgBarService } from './msg-bar.service';

describe('MsgBarService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [MsgBarService],
    });
  });

  it('should be created', inject([MsgBarService], (service: MsgBarService) => {
    expect(service).toBeTruthy();
  }));
});
