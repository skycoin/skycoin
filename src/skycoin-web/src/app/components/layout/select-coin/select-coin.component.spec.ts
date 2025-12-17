import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { FormsModule } from '@angular/forms';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { MatDialogModule } from '@angular/material/dialog';

import { SelectCoinComponent } from './select-coin.component';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { MockCustomMatDialogService } from '../../../utils/test-mocks';

describe('SelectCoinComponent', () => {
  let component: SelectCoinComponent;
  let fixture: ComponentFixture<SelectCoinComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SelectCoinComponent ],
      imports: [
        FormsModule,
        MatDialogModule
      ],
      providers: [
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService }
      ],
      schemas: [ NO_ERRORS_SCHEMA ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectCoinComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
