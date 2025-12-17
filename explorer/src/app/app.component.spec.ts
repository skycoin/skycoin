import { TestBed, waitForAsync } from '@angular/core/testing';
import { Router } from '@angular/router';
import { of } from 'rxjs';

import { AppComponent } from './app.component';
import { LanguageService } from './services/language/language.service';
import { ExplorerService } from './services/explorer/explorer.service';

class MockRouter {
  events = of(null);
}

class MockLanguageService {
  initialize() {}
}

class MockExplorerService {
  initialize() {}
}

describe('AppComponent', () => {
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [
        AppComponent
      ],
      providers: [
        { provide: Router, useClass: MockRouter },
        { provide: LanguageService, useClass: MockLanguageService },
        { provide: ExplorerService, useClass: MockExplorerService },
      ]
    }).compileComponents();
  }));

  it('should create the app', waitForAsync(() => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  }));
});
