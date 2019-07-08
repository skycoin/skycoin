import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MultipleDestinationsDialogComponent } from './multiple-destinations-dialog.component';

describe('MultipleDestinationsDialogComponent', () => {
  let component: MultipleDestinationsDialogComponent;
  let fixture: ComponentFixture<MultipleDestinationsDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MultipleDestinationsDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MultipleDestinationsDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
