import { Pipe, PipeTransform, Component } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { Observable } from 'rxjs';
import { Subject } from 'rxjs';
import { ReplaySubject } from 'rxjs';

import { BaseCoin } from '../coins/basecoin';
import { Wallet } from '../app.datatypes';
import { BalanceEvent, BalanceStates } from '../services/wallet/balance.service';

// -- Components
@Component({
  template: `
    <span clipboard="test data"></span>
    <input type="text" appNumberField>`
})
export class TestComponent {
}

// --- Pipes
@Pipe({name: 'translate'})
export class MockTranslatePipe implements PipeTransform {
  transform() {
    return 'translated value';
  }
}

@Pipe({name: 'tellerStatus'})
export class MockTellerStatusPipe implements PipeTransform {
  transform() {
    return 'transformed value';
  }
}

@Pipe({ name: 'dateTime' })
export class MockDateTimePipe implements PipeTransform {
  transform() {
    return 'transformed value';
  }
}

// --- Services

export class MockNavBarService {
  activeComponent = new BehaviorSubject({});

  showSwitch(leftText: any, rightText: any) {}

  hideSwitch() {}

  setActiveComponent() {}
}

export class MockGlobalsService {
  private nodeVersion = new BehaviorSubject<string>('0.24.0');

  setNodeVersion(version: string) {
    this.nodeVersion.next(version);
  }

  getValidNodeVersion(): Observable<string> {
    return this.nodeVersion.filter(version => version !== null).first();
  }
}

export class MockPurchaseService {
  all(): Observable<any[]> {
    return Observable.of([]);
  }
}

export class MockCoinService {
  coins = [];
  customNodeUrls = {};

  currentCoin = new BehaviorSubject<BaseCoin>({
    id: 1,
    priceTickerId: 'btc-bitcoin',
    nodeUrl: 'nodeUrl',
    coinName: 'test coin',
    coinSymbol: 'test',
    hoursName: 'Test Hours',
    coinExplorer: 'testUrl',
    imageName: 'imageName.png',
    gradientName: 'imageName.png',
    iconName: 'icon.png',
    bigIconName: 'big-icon.png',
  });

  removeTemporarilyAllowedCoin() {}
}

export class MockWalletService {
  haveWallets: Observable<boolean> = Observable.of(true);

  get addresses(): Observable<any[]> {
    return Observable.of([]);
  }

  wallets = new BehaviorSubject<any[]>([]);

  get currentWallets(): Observable<Wallet[]> {
    return Observable.of([]);
  }

  scanAddresses(wallet, onProgressChanged): Observable<void> {
    return Observable.of();
  }
}

export class MockHistoryService {
  transactions(): Observable<any[]> {
    return Observable.of([]);
  }

  getAllPendingTransactions() {
    return Observable.of([]);
  }

  getTransactionDetails() {
    return Observable.of({});
  }
}

export class MockSpendingService {
  outputsWithWallets() {
    return Observable.of([]);
  }
}

export class MockBalanceService {
  hasPendingTransactions: Subject<boolean> = new ReplaySubject<boolean>();
  lastBalancesUpdateTime: Date = new Date();

  get totalBalance(): BehaviorSubject<BalanceEvent> {
    return new BehaviorSubject<BalanceEvent>({
      state: BalanceStates.Updating
    });
  }

  startGettingBalances() { }

  stopGettingBalances() { }
}

export class MockPriceService {
  price: Subject<number> = new BehaviorSubject<number>(null);
}

export class MockBlockchainService {
  get progress() {
    return Observable.of();
  }

  get currentMaxDecimals() {
    return 3;
  }

  get burnRate() {
    return 2;
  }

  lastBlock(): Observable<any> {
    return Observable.of({});
  }

  coinSupply(): Observable<any> {
    return Observable.of({});
  }

  loadBlockchainBlocks() {
  }
}

export class MockMsgBarService {
  set msgBarComponent(value) { }
  show() { }
  hide() { }
  showError(body: string, duration = 20000) { }
  showWarning(body: string, duration = 20000) { }
  showDone(body: string, duration = 10000) { }
}

export class MockTranslateService {
  onLangChange = Observable.of({});

  get(key: string | Array<string>, interpolateParams?: Object): Observable<string | any> {
    return Observable.of({});
  }

  addLangs(langs: Array<string>): void {
  }

  setDefaultLang(lang: string): void {
  }

  use(lang: string): Observable<any> {
    return Observable.of({});
  }
}

export class MockApiService {
  get(url: string) {
    if (url === 'network/connections') {
      return Observable.of({ connections: [] });
    } else {
      return Observable.of({});
    }
  }
}

export class MockClipboardService {
}

export class MockLanguageService {
  currentLanguage = new BehaviorSubject<string>('en');

  loadLanguageSettings() {
  }
}

export class MockCustomMatDialogService {
  get showingDialog() {
    return Observable.of(false);
  }
}

export class MockMatDialogRef<T> {
  close(dialogResult?: any) {
  }
}
