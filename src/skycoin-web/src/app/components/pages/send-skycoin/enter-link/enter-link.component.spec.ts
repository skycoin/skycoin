import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { UntypedFormBuilder } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';

import { EnterLinkComponent } from './enter-link.component';
import { MockTranslatePipe } from '../../../../utils/test-mocks';

describe('EnterLinkComponent', () => {
  let component: EnterLinkComponent;
  let fixture: ComponentFixture<EnterLinkComponent>;
  let close: jasmine.Spy;

  beforeEach(waitForAsync(() => {
    close = jasmine.createSpy('close');
    TestBed.configureTestingModule({
      declarations: [ EnterLinkComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        UntypedFormBuilder,
        { provide: MatDialogRef, useValue: { close } },
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EnterLinkComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('closes with the data of a valid link', () => {
    component.form.get('link')!.setValue('skycoin:addr1?amount=1&hours=2');
    component.process();
    expect(close).toHaveBeenCalledWith({ address: 'addr1', amount: '1', hours: '2' });
  });

  it('shows an error for text that is not a link', () => {
    component.form.get('link')!.setValue('addr1');
    component.process();
    expect(close).not.toHaveBeenCalled();
    expect(component.error).toBe('send.fill-with-link.invalid-link-error');
  });
});
