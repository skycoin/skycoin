import { TestBed, inject } from '@angular/core/testing';

import { AppService } from './app.service';
import { Http } from '@angular/http';

describe('AppService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [AppService, Http],
    });
  });

  it('should be created', inject([AppService], (service: AppService) => {
    expect(service).toBeTruthy();
  }));
});
