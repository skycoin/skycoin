import {
  Component,
  ComponentRef,
  Input,
  NO_ERRORS_SCHEMA,
  Pipe,
  PipeTransform,
  ViewChild,
  ViewContainerRef,
} from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Observable, Subject } from 'rxjs';

import { ApiService } from '../../../services/api/api.service';
import { ExplorerService } from '../../../services/explorer/explorer.service';
import { RichlistComponent } from './richlist.component';

/**
 * Guards the change detection contract of an OnPush component.
 *
 * Every page in this project loads its data from an observable and assigns the
 * result to a field. Under ChangeDetectionStrategy.OnPush that assignment
 * repaints nothing on its own: the component has to be marked dirty, or the view
 * keeps showing whatever it showed before. The failure is silent — no error, no
 * console warning, just a page that never leaves "Waiting for data".
 *
 * The e2e suite cannot be relied on to catch this. It runs against the pinned
 * blockchain-180 database, which answers fast enough that the response can
 * arrive inside the first change detection pass, rendering the view whether or
 * not it was ever marked dirty. The same broken build does fail those specs
 * against a mainnet node, where the query takes about 1.8 seconds. Catching this
 * dependably means controlling *when* the value arrives instead of racing the
 * network, which is what the Subject below is for.
 *
 * Three details matter and are all easy to get wrong:
 *
 *  - Change detection is driven from a host component that uses the default
 *    strategy. Calling detectChanges() on a fixture of the OnPush component
 *    itself checks that component's own view whether or not it is dirty, which
 *    hides the exact bug this spec exists to catch.
 *
 *  - RichlistComponent is reached only through the router and declares no
 *    selector, so it cannot be placed in a host template. It is created through
 *    a ViewContainerRef instead; the resulting view is still a child of the
 *    host's view, so OnPush is honoured the same way.
 *
 *  - The data is emitted after the first render, so the assignment happens
 *    outside the initial pass, as it would after a real HTTP response.
 */

@Pipe({ name: 'translate', standalone: false })
class StubTranslatePipe implements PipeTransform {
  transform(value: string): string { return value; }
}

@Pipe({ name: 'amount', standalone: false })
class StubAmountPipe implements PipeTransform {
  transform(value: any): string { return String(value); }
}

/**
 * Stands in for LoadingComponent. It renders the message it is given so the
 * error path can be asserted against the DOM rather than against a field, which
 * is what makes that case a change detection assertion too.
 */
@Component({
  selector: 'app-loading',
  template: '<span class="stub-error">{{ longErrorMsg }}</span>',
  standalone: false,
})
class StubLoadingComponent {
  @Input() longErrorMsg!: string;
}

@Component({
  selector: 'app-cd-host',
  template: '<ng-container #slot></ng-container>',
  standalone: false,
})
class HostComponent {
  @ViewChild('slot', { read: ViewContainerRef, static: true })
  slot!: ViewContainerRef;

  create(): ComponentRef<RichlistComponent> {
    return this.slot.createComponent(RichlistComponent);
  }
}

class MockExplorerService {
  getAddressName(): string { return ''; }
}

class MockApiService {
  /** Controls exactly when the richlist response arrives. */
  readonly richlist = new Subject<any[]>();
  getRichlist(): Observable<any[]> { return this.richlist.asObservable(); }
}

describe('RichlistComponent change detection', () => {
  let fixture: ComponentFixture<HostComponent>;
  let componentRef: ComponentRef<RichlistComponent>;
  let api: MockApiService;

  beforeEach(async () => {
    api = new MockApiService();

    // PageBaseComponent persists scroll position and server responses onto
    // window.history.state and writes to it unconditionally. A fresh karma frame
    // starts with a null state, which the component would throw on; in the
    // running application the router always leaves an object there.
    window.history.replaceState({}, '', window.location.pathname);

    await TestBed.configureTestingModule({
      declarations: [
        HostComponent,
        RichlistComponent,
        StubLoadingComponent,
        StubTranslatePipe,
        StubAmountPipe,
      ],
      providers: [
        { provide: ApiService, useValue: api },
        { provide: ExplorerService, useClass: MockExplorerService },
      ],
      // routerLink is not relevant to change detection.
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(HostComponent);
    componentRef = fixture.componentInstance.create();
  });

  function renderedRows(): string[] {
    return Array.from(
      fixture.nativeElement.querySelectorAll('a.-row') as NodeListOf<HTMLElement>,
    ).map(el => el.textContent || '');
  }

  it('renders nothing before the response arrives', () => {
    fixture.detectChanges();
    expect(renderedRows().length).toEqual(0);
  });

  // This is the guard. Removing markForCheck from the subscription in
  // richlist.component.ts turns it red with "Expected 0 to equal 2", which is
  // the silent stale-view failure reproduced deterministically.
  it('renders entries that arrive after the first change detection pass', () => {
    // First pass: the component subscribes; the response has not arrived.
    fixture.detectChanges();
    expect(renderedRows().length).toEqual(0);

    // The response lands outside any change detection pass, as a real HTTP
    // response would.
    api.richlist.next([
      { address: 'AddressOne', coins: 100 },
      { address: 'AddressTwo', coins: 50 },
    ]);

    // Driven from the parent, so the OnPush child is only re-rendered if it was
    // marked dirty. Without markForCheck in the subscription this stays empty.
    fixture.detectChanges();

    expect(componentRef.instance.entries.length).toEqual(2);
    const rows = renderedRows();
    expect(rows.length).toEqual(2);
    expect(rows[0]).toContain('AddressOne');
    expect(rows[1]).toContain('AddressTwo');
  });

  // Documents the error path. Note that this one is not a change detection
  // guard: it still passes with markForCheck removed, because the message is
  // rendered by a child component whose input is refreshed on the pass that
  // follows. Only the test above reliably detects a missing mark.
  it('renders the error message when the request fails', () => {
    fixture.detectChanges();

    api.richlist.error(new Error('request failed'));
    fixture.detectChanges();

    expect(componentRef.instance.longErrorMsg).toEqual('general.longLoadingErrorMsg');
    expect(fixture.nativeElement.querySelector('.stub-error').textContent)
      .toContain('general.longLoadingErrorMsg');
  });
});
