import {
  Component,
  ComponentRef,
  NO_ERRORS_SCHEMA,
  Pipe,
  PipeTransform,
  ViewChild,
  ViewContainerRef,
} from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Observable, Subject } from 'rxjs';

import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { MsgBarComponent } from './msg-bar.component';

/**
 * Behavioural guard for OnPush change detection in the web wallet.
 *
 * Every component here loads state from an observable or a timer and assigns it
 * to a field. Under ChangeDetectionStrategy.OnPush that assignment repaints
 * nothing on its own — the view has to be marked dirty, or it keeps showing what
 * it showed before. Nothing throws and nothing is logged; the view just stops
 * updating, which in a wallet means stale information in front of someone about
 * to move funds.
 *
 * ci-scripts/check-onpush-marks.js sweeps every component structurally so no
 * asynchronous callback ships without a mark. This spec is the other half,
 * proving the contract holds at runtime for a representative component, so the
 * structural sweep is known to be guarding something real.
 *
 * The setup mirrors the guards in /explorer and /src/gui/static:
 *
 *  - Change detection is driven from a host using the default strategy. Calling
 *    detectChanges() on a fixture of the OnPush component itself checks that
 *    component's view whether or not it is dirty, hiding the bug.
 *  - The component is created through a ViewContainerRef, so it is a child view
 *    of the host exactly as it is in the running application.
 *  - State changes after the first render, putting the assignment outside the
 *    initial pass.
 */

@Pipe({ name: 'translate', standalone: false })
class StubTranslatePipe implements PipeTransform {
  transform(value: string): string { return value; }
}

@Component({
  selector: 'app-cd-host',
  template: '<ng-container #slot></ng-container>',
  standalone: false,
})
class HostComponent {
  @ViewChild('slot', { read: ViewContainerRef, static: true })
  slot!: ViewContainerRef;

  create(): ComponentRef<MsgBarComponent> {
    return this.slot.createComponent(MsgBarComponent);
  }
}

class MockDialogService {
  /** Controls exactly when the "a dialog is open" state changes. */
  readonly showingDialog = new Subject<boolean>();
}

describe('MsgBarComponent change detection', () => {
  let fixture: ComponentFixture<HostComponent>;
  let componentRef: ComponentRef<MsgBarComponent>;
  let dialog: MockDialogService;

  beforeEach(async () => {
    dialog = new MockDialogService();

    await TestBed.configureTestingModule({
      declarations: [HostComponent, MsgBarComponent, StubTranslatePipe],
      providers: [{ provide: CustomMatDialogService, useValue: dialog }],
      // mat-icon is irrelevant to change detection.
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(HostComponent);
    componentRef = fixture.componentInstance.create();
  });

  function container(): HTMLElement | null {
    return fixture.nativeElement.querySelector('.main-container');
  }

  it('renders nothing while hidden', () => {
    fixture.detectChanges();
    expect(container()).toBeNull();
  });

  // This is the guard. Removing markForCheck from the showingDialog
  // subscription in msg-bar.component.ts turns it red, because the class stays
  // on the element that the first render put it on.
  //
  // The assertion is on a class the template derives from UiIsShowingModalWindow
  // rather than on the component field, so it fails only when the *view* is
  // stale — which is the actual bug — and not merely when the field is unset.
  it('repaints when dialog state arrives after the first render', () => {
    const component = componentRef.instance;
    component.config.text = 'some.message';
    // Set before the first pass, so the bar is rendered without needing a mark.
    component.show();

    fixture.detectChanges();
    const el = container();
    expect(el).not.toBeNull();
    expect(el!.classList).toContain('fix-small-screen-position');

    // Arrives outside any change detection pass, as a real event would.
    dialog.showingDialog.next(true);

    // Driven from the parent, so the OnPush child only re-renders if the
    // subscription marked it dirty.
    fixture.detectChanges();

    expect(component.UiIsShowingModalWindow).toBeTruthy();
    expect(container()!.classList).not.toContain('fix-small-screen-position');
  });
});
