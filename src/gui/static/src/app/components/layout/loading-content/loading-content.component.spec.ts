import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LoadingContentComponent } from './loading-content.component';

describe('LoadingContentComponent', () => {
  let component: LoadingContentComponent;
  let fixture: ComponentFixture<LoadingContentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LoadingContentComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoadingContentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
