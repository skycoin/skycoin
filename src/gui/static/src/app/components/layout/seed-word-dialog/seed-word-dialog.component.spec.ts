import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SeedWordDialogComponent } from './seed-word-dialog.component';

describe('SeedWordDialogComponent', () => {
  let component: SeedWordDialogComponent;
  let fixture: ComponentFixture<SeedWordDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SeedWordDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SeedWordDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
