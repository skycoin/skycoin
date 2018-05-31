import { TestBed, inject } from '@angular/core/testing';

import { WizardGuardService } from './wizard-guard.service';

describe('WizardGuardService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [WizardGuardService],
    });
  });

  it('should be created', inject([WizardGuardService], (service: WizardGuardService) => {
    expect(service).toBeTruthy();
  }));
});
