import { TestBed, inject } from '@angular/core/testing';

import { AppService } from './app.service';
import { HttpClient } from '@angular/common/http';

describe('AppService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [AppService, HttpClient],
    });
  });

  it('should be created', inject([AppService], (service: AppService) => {
    expect(service).toBeTruthy();
  }));
});
