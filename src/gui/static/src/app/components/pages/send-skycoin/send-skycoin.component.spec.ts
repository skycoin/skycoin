import * as sinon from 'sinon';
import { async, inject, ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpModule, XHRBackend } from '@angular/http';
import { MdCardModule, MdMenuModule, MdSelectModule, MdSnackBarModule, MdTooltipModule } from '@angular/material';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { RouterTestingModule } from '@angular/router/testing';
import { MockBackend } from '@angular/http/testing';

import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import { ButtonComponent } from '../../layout/button/button.component';
import { HeaderComponent } from '../../layout/header/header.component';
import { NavBarComponent } from '../../layout/header/nav-bar/nav-bar.component';
import { TopBarComponent } from '../../layout/header/top-bar/top-bar.component';
import { SendSkycoinComponent } from './send-skycoin.component';
import { DateTimePipe } from '../../../pipes/date-time.pipe';
import { ApiService } from '../../../services/api.service';
import { PriceService } from '../../../price.service';
import { WalletService } from '../../../services/wallet.service';

describe('SendSkycoinComponent', () => {
  let component: SendSkycoinComponent;
  let mockBackend: MockBackend;
  let fixture: ComponentFixture<SendSkycoinComponent>;
  let sendButton: any;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        ButtonComponent,
        DateTimePipe,
        HeaderComponent,
        NavBarComponent,
        SendSkycoinComponent,
        TopBarComponent
      ],
      imports: [
        HttpModule,
        MdCardModule,
        MdMenuModule,
        MdSelectModule,
        MdSnackBarModule,
        MdTooltipModule,
        NgxDatatableModule,
        NoopAnimationsModule,
        ReactiveFormsModule,
        RouterTestingModule
      ],
      providers: [
        ApiService,
        PriceService,
        WalletService,
        { provide: XHRBackend, useClass: MockBackend }
      ]
    })
      .compileComponents();
  }));

  beforeEach(inject([XHRBackend], (backend: MockBackend) => {
    mockBackend = backend;
    fixture = TestBed.createComponent(SendSkycoinComponent);
    component = fixture.componentInstance;
    component.ngOnInit();
    fixture.detectChanges();

    sendButton = fixture.debugElement.nativeElement.querySelector('app-button');
  }));

  afterEach(() => {
    fixture.destroy();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('send should not make any backend call if wallet is not chosen', () => {
    const backendCallSpy = sinon.spy();
    mockBackend.connections.subscribe(backendCallSpy);

    fixture.detectChanges();
    component.button.onClick();

    sinon.assert.notCalled(backendCallSpy);
  });

  it('send should not make any backend call if address is not set', () => {
    const backendCallSpy = sinon.spy();
    mockBackend.connections.subscribe(backendCallSpy);

    component.form.controls.wallet.setValue({
      meta: { filename: 'test' },
      balance: 1000000
    });

    fixture.detectChanges();
    component.button.onClick();

    sinon.assert.notCalled(backendCallSpy);
  });

  it('send should not make any backend call if amount is not set', () => {
    const backendCallSpy = sinon.spy();
    mockBackend.connections.subscribe(backendCallSpy);

    component.form.controls.wallet.setValue({
      meta: { filename: 'test' },
      balance: 1000000
    });
    component.form.controls.address.setValue('0xABC');

    fixture.detectChanges();
    component.button.onClick();

    sinon.assert.notCalled(backendCallSpy);
  });

  it('send should not make any backend call if amount exceeds balance', () => {
    const backendCallSpy = sinon.spy();
    mockBackend.connections.subscribe(backendCallSpy);

    component.form.controls.wallet.setValue({
      meta: { filename: 'test' },
      balance: 1000000
    });
    component.form.controls.address.setValue('0xABC');
    component.form.controls.amount.setValue(2);

    fixture.detectChanges();
    component.button.onClick();

    sinon.assert.notCalled(backendCallSpy);
  });

  it('send should make backend call if send form is valid', () => {
    let request;
    mockBackend.connections.subscribe((connection) => {
      request = connection.request;
    });

    component.form.controls.wallet.setValue({
      filename: 'test',
      coins: 100000
    });
    component.form.controls.address.setValue('0xABC');
    component.form.controls.amount.setValue(1);

    fixture.detectChanges();
    component.button.onClick();

    expect(request).toBeDefined();
    expect(request.url).toBe('http://127.0.0.1:6420/wallet/spend?');
    expect(request.getBody()).toBe('id=test&dst=0xABC&coins=1000000');
  });
});
