import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ArrowLinkComponent } from './arrow-link.component';

describe('ArrowLinkComponent', () => {
  let component: ArrowLinkComponent;
  let fixture: ComponentFixture<ArrowLinkComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ArrowLinkComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ArrowLinkComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
