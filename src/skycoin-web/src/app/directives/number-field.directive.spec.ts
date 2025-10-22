import { TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';

import { NumberFieldDirective } from './number-field.directive';
import { TestComponent } from '../utils/test-mocks';

describe('Directive: NumberField', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [TestComponent, NumberFieldDirective]
    });
  });

  it('should create an instance', () => {
    const fixture = TestBed.createComponent(TestComponent);
    const directiveEl = fixture.debugElement.query(By.directive(NumberFieldDirective));
    expect(directiveEl).not.toBeNull();
  });
});
