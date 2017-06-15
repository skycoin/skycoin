import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { HeaderComponent } from './home/header/header.component';
import { NavigationComponent } from './home/navigation/navigation.component';
import { CoinbalanceComponent } from './home/coinbalance/coinbalance.component';
import {WalletService} from "./services/wallet.service";
import { RouterModule, Routes } from '@angular/router';
import { AddressComponent } from './address/address.component';
import { ReceiveComponent } from './receive/receive.component';
import { QRCodeModule } from 'angular2-qrcode';

const appRoutes: Routes = [
  { path: 'address', component: AddressComponent },
  {
    path: 'home',
    component: HomeComponent,
    data: { title: 'Sycoin New Wallet' }
  },
  { path: '',
    redirectTo: '/home',
    pathMatch: 'full'
  },
  { path: 'receive',
    component: ReceiveComponent,
  },
  { path: '*',
    redirectTo: '/home',
    pathMatch: 'full'
  },
];


@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    HeaderComponent,
    NavigationComponent,
    CoinbalanceComponent,
    AddressComponent,
    ReceiveComponent
  ],
  imports: [
    BrowserModule,
    RouterModule.forRoot(appRoutes),
    FormsModule,
    HttpModule,
    QRCodeModule,
  ],
  providers: [WalletService],
  bootstrap: [AppComponent]
})
export class AppModule { }
