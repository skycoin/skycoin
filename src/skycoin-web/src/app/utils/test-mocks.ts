import { Pipe, PipeTransform, Component, ChangeDetectionStrategy, EventEmitter } from '@angular/core';
import { BehaviorSubject, Observable, Subject, ReplaySubject, of } from 'rxjs';
import { filter, first } from 'rxjs';

import { BaseCoin } from '../coins/basecoin';
import { Wallet } from '../app.datatypes';
import { BalanceEvent, BalanceStates } from '../services/wallet/balance.service';
import { CoinHealth } from '../services/node-health.service';
import { MsgBarComponent } from '../components/layout/msg-bar/msg-bar.component';

// -- Components
@Component({
    template: `
    <span clipboard="test data"></span>
    <input type="text" appNumberField>`,
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class TestComponent {
}

// --- Pipes
@Pipe({
    name: 'translate',
    standalone: false
})
export class MockTranslatePipe implements PipeTransform {
  transform() {
    return 'translated value';
  }
}

@Pipe({
    name: 'tellerStatus',
    standalone: false
})
export class MockTellerStatusPipe implements PipeTransform {
  transform() {
    return 'transformed value';
  }
}

@Pipe({
    name: 'dateTime',
    standalone: false
})
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
    return this.nodeVersion.pipe(filter(version => version !== null), first());
  }
}

export class MockPurchaseService {
  all(): Observable<any[]> {
    return of([]);
  }
}

export class MockCoinService {
  coins = [];
  customNodeUrls = {};

  /** Services wait on this before their first request; emit straight away. */
  coinsLoaded = new ReplaySubject<boolean>(1);

  currentCoin = new BehaviorSubject<BaseCoin>(new BaseCoin({
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
  }));

  constructor() {
    this.coinsLoaded.next(true);
  }

  removeTemporarilyAllowedCoin() {}
}

export class MockWalletService {
  haveWallets: Observable<boolean> = of(true);

  get addresses(): Observable<any[]> {
    return of([]);
  }

  wallets = new BehaviorSubject<any[]>([]);

  get currentWallets(): Observable<Wallet[]> {
    return of([]);
  }

  scanAddresses(wallet: Wallet, onProgressChanged: EventEmitter<any>): Observable<void> {
    return of();
  }
}

export class MockHistoryService {
  transactions(): Observable<any[]> {
    return of([]);
  }

  getAllPendingTransactions() {
    return of([]);
  }

  getTransactionDetails() {
    return of({});
  }
}

export class MockSpendingService {
  outputsWithWallets() {
    return of([]);
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
  price: Subject<number> = new BehaviorSubject<number>(null as any);
}

export class MockBlockchainService {
  get progress() {
    return of();
  }

  get currentMaxDecimals() {
    return 3;
  }

  get burnRate() {
    return 2;
  }

  lastBlock(): Observable<any> {
    return of({});
  }

  coinSupply(): Observable<any> {
    return of({});
  }

  loadBlockchainBlocks() {
  }
}

export class MockNodeHealthService {
  /** Reports the node as reachable so forms leave their disabled state. */
  check(coin: BaseCoin): Observable<CoinHealth> {
    return of(this.health(coin));
  }

  watch(coin: BaseCoin, intervalMs = 20000): Observable<CoinHealth> {
    return of(this.health(coin));
  }

  private health(coin: BaseCoin): CoinHealth {
    return {
      coinId: coin.id,
      coinName: coin.coinName,
      coinSymbol: coin.coinSymbol,
      isBitcoin: false,
      available: true,
      height: 1,
    };
  }
}

export class MockHwWalletService {
  /** SpendingService only reaches for these two when spending from a device. */
  checkIfCorrectHwConnected(address: string): Observable<boolean> {
    return of(true);
  }

  signTransaction(inputs: any[], outputs: any[]): Observable<any> {
    return of(null);
  }
}

export class MockMsgBarService {
  set msgBarComponent(value: MsgBarComponent) { }
  show() { }
  hide() { }
  showError(body: string, duration = 20000) { }
  showWarning(body: string, duration = 20000) { }
  showDone(body: string, duration = 10000) { }
}

export class MockTranslateService {
  onLangChange = of({});

  get(key: string | Array<string>, interpolateParams?: Object): Observable<string | any> {
    return of({});
  }

  addLangs(langs: Array<string>): void {
  }

  setFallbackLang(lang: string): void {
  }

  use(lang: string): Observable<any> {
    return of({});
  }
}

export class MockApiService {
  get(url: string) {
    if (url === 'network/connections') {
      return of({ connections: [] });
    } else {
      return of({});
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
    return of(false);
  }
}

export class MockMatDialogRef<T> {
  close(dialogResult?: any) {
  }
}
