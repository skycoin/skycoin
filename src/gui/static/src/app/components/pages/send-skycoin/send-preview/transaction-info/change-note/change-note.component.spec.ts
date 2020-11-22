import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ChangeNoteComponent } from './change-note.component';

describe('ChangeNoteComponent', () => {
  let component: ChangeNoteComponent;
  let fixture: ComponentFixture<ChangeNoteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ChangeNoteComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ChangeNoteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
