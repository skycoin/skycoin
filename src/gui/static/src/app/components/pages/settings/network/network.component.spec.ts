import {
  Component,
  NO_ERRORS_SCHEMA,
  Pipe,
  PipeTransform,
  ViewChild,
  ViewContainerRef,
  ComponentRef,
} from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Observable, Subject } from 'rxjs';

import { NetworkService } from '../../../../services/network.service';
import { NetworkComponent } from './network.component';

/**
 * Behavioural guard for OnPush change detection in the desktop wallet.
 *
 * Every page here loads data from an observable and assigns it to a field. Under
 * ChangeDetectionStrategy.OnPush that assignment repaints nothing on its own —
 * the view has to be marked dirty, or it keeps showing what it showed before.
 * Nothing throws and nothing is logged; the page simply never updates, which in
 * a wallet means a stale balance or a transaction list that never appears.
 *
 * ci-scripts/check-onpush-marks.js sweeps every component structurally, so no
 * asynchronous callback can quietly ship without a mark. This spec is the other
 * half: it proves the contract actually holds at runtime for a representative
 * component, so the structural check is known to be guarding something real.
 *
 * Three parts of the setup are load-bearing:
 *
 *  - Change detection is driven from a host using the default strategy. Calling
 *    detectChanges() on a fixture of the OnPush component itself checks that
 *    component's view whether or not it is dirty, hiding the very bug this
 *    exists to catch.
 *
 *  - The component is created through a ViewContainerRef, so it behaves as a
 *    child view of the host exactly as it does under the router.
 *
 *  - The peers are emitted after the first render, putting the assignment
 *    outside the initial pass, where a real HTTP response lands.
 */

@Pipe({ name: 'translate', standalone: false })
class StubTranslatePipe implements PipeTransform {
  transform(value: string): string { return value; }
}

@Pipe({ name: 'dateFromNow', standalone: false })
class StubDateFromNowPipe implements PipeTransform {
  transform(value: any): string { return String(value); }
}

@Component({
  selector: 'app-cd-host',
  template: '<ng-container #slot></ng-container>',
  standalone: false,
})
class HostComponent {
  @ViewChild('slot', { read: ViewContainerRef, static: true })
  slot!: ViewContainerRef;

  create(): ComponentRef<NetworkComponent> {
    return this.slot.createComponent(NetworkComponent);
  }
}

class MockNetworkService {
  /** Controls exactly when the connection list arrives. */
  readonly peers = new Subject<any[]>();
  connections(): Observable<any[]> { return this.peers.asObservable(); }
}

describe('NetworkComponent change detection', () => {
  let fixture: ComponentFixture<HostComponent>;
  let componentRef: ComponentRef<NetworkComponent>;
  let network: MockNetworkService;

  beforeEach(async () => {
    network = new MockNetworkService();

    await TestBed.configureTestingModule({
      declarations: [
        HostComponent,
        NetworkComponent,
        StubTranslatePipe,
        StubDateFromNowPipe,
      ],
      providers: [
        { provide: NetworkService, useValue: network },
      ],
      // app-header, app-loading-content, mat-icon and matTooltip are all
      // irrelevant to change detection.
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(HostComponent);
    componentRef = fixture.componentInstance.create();
  });

  function renderedRows(): string[] {
    return Array.from(
      fixture.nativeElement.querySelectorAll('.-row') as NodeListOf<HTMLElement>,
    ).map(el => el.textContent || '');
  }

  it('renders no peers before the response arrives', () => {
    fixture.detectChanges();
    expect(renderedRows().length).toEqual(0);
  });

  // This is the guard. Removing markForCheck from the subscription in
  // network.component.ts turns it red with "Expected 0 to equal 2".
  it('renders peers that arrive after the first change detection pass', () => {
    fixture.detectChanges();
    expect(renderedRows().length).toEqual(0);

    network.peers.next([
      { address: '10.0.0.1:6000', listenPort: 6000, source: 'x', height: 42, outgoing: true, lastSent: 1, lastReceived: 2 },
      { address: '10.0.0.2:6001', listenPort: 6001, source: 'y', height: 43, outgoing: false, lastSent: 3, lastReceived: 4 },
    ]);

    // Driven from the parent, so the OnPush child only re-renders if it was
    // marked dirty. Without markForCheck this stays empty.
    fixture.detectChanges();

    expect(componentRef.instance.peers.length).toEqual(2);
    const rows = renderedRows();
    expect(rows.length).toEqual(2);
    expect(rows[0]).toContain('10.0.0.1');
    expect(rows[1]).toContain('10.0.0.2');
  });
});
