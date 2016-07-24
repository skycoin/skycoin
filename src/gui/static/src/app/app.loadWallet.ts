

//import {Component, OnInit, ViewChild} from 'app/angular2/core';
//import {ROUTER_DIRECTIVES, OnActivate} from 'app/angular2/router';
import {Component, OnInit, ViewChild} from 'angular2/core';
import {ROUTER_DIRECTIVES, OnActivate} from 'angular2/router';

import {Http, HTTP_BINDINGS, Response} from 'angular2/http';
import {HTTP_PROVIDERS, Headers} from 'angular2/http';
import {Observable} from 'rxjs/Observable';
import {Observer} from 'rxjs/Observer';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import {QRCodeComponent} from './ng2-qrcode.ts'

@Component({
    selector: 'load-wallet',
    directives: [ROUTER_DIRECTIVES, QRCodeComponent],
    providers: [],
    templateUrl: 'app/templates/wallet.html'
})

export class loadWalletComponent implements OnInit {
    //Declare default varialbes
    wallets : Array<any>;
    progress: any;
    spendid: string;
    spendaddress: string;
    sendDisable: boolean;
    readyDisable: boolean;
    displayMode: DisplayModeEnum;
    displayModeEnum = DisplayModeEnum;

    QrAddress: string;
    QrIsVisible: boolean;

    NewWalletIsVisible: boolean;
    EditWalletIsVisible: boolean;
    loadSeedIsVisible: boolean;

    walletname: string;
    walletId: string;

    historyTable: Array<any>;
    pendingTable: Array<any>;
    addresses: Array<any>;
    connections: Array<any>;

    //Constructor method for load HTTP object
    constructor(private http: Http) { }

    //Init function for load default value
    ngOnInit() {
        this.displayMode = DisplayModeEnum.first;
        this.loadWallet();
        this.loadConnections();
        this.loadProgress();

        //Set interval function for load wallet every 15 seconds
        setInterval(() => {
            this.loadWallet();
            console.log("Refreshing balance");
        }, 15000);

        //Enable Send tab "textbox" and "Ready" button by default
        this.sendDisable = true;
        this.readyDisable = false;
        this.pendingTable = [];

        if(localStorage.getItem('historyAddresses') != null){
            this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
        } else {
            localStorage.setItem('historyAddresses',JSON.stringify([]));
            this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
        }
    }

    //Ready button function for disable "textbox" and enable "Send" button for ready to send coin
    ready(spendId, spendaddress, spendamount){
        if(!spendId){
            alert("Please select from id");
            return false;
        }
        if(!spendaddress){
            alert("Please enter pay to");
            return false;
        }
        if(!spendamount){
            alert("Please enter amount");
            return false;
        }
        this.readyDisable = true;
        this.sendDisable = false;
    }

    //Load wallet function
    loadWallet(){
        this.http.post('/wallets', '')
            .map((res:Response) => res.json())
            .subscribe(
                data => {
                    this.wallets = data;

                    //Load Balance for each wallet
                    var inc = 0;
                    for(var item in data){
                        var address = data[inc].meta.filename;
                        this.loadWalletItem(address, inc);
                        inc++;
                    }
                    //Load Balance for each wallet end

                },
                err => console.log("Error on load wallet: "+err),
                () => console.log('Wallet load done')
            );
    }
    loadWalletItem(address, inc){
        //Set http headers
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        this.http.get('/wallet/balance?id=' + address, { headers: headers })
            .map((res) => res.json())
            .subscribe(
                //Response from API
                response => {
                    console.log('load done: ' + inc, response);
                    this.wallets[inc].balance = response.confirmed.coins / 1000000;
                }, err => console.log("Error on load balance: " + err), () => console.log('Balance load done'))
    }
    loadConnections() {
        this.http.post('/network/connections', '')
            .map((res) => res.json())
            .subscribe(data => {
                console.log("connections", data);
                this.connections = data.connections;
            }, err => console.log("Error on load wallet: " + err), () => console.log('Wallet load done'));
    }
    //Load progress function for Skycoin
    loadProgress(){
        //Post method executed
        this.http.post('/blockchain/progress', '')
            .map((res:Response) => res.json())
            .subscribe(
                //Response from API
                response => { this.progress = (parseInt(response.current,10)+1) / parseInt(response.highest,10) * 100 },
                err => console.log("Error on load progress: "+err),
                () => console.log('Progress load done:' + this.progress)
            );
    }

    //Switch tab function
    switchTab(mode: DisplayModeEnum, wallet) {
        //"Textbox" and "Ready" button enable in Send tab while switching tabs
        this.sendDisable = true;
        this.readyDisable = false;

        this.displayMode = mode;
        if(wallet){
            this.spendid = wallet.meta.filename;
        }
    }

    //Show QR code function for show QR popup
    showQR(wallet){
        this.QrAddress = wallet.entries[0].address;
        this.QrIsVisible = true;
    }
    //Hide QR code function for hide QR popup
    hideQrPopup(){
        this.QrIsVisible = false;
    }

    //Show wallet function for view New wallet popup
    showNewWalletDialog(){
        this.NewWalletIsVisible = true;
    }
    //Hide wallet function for hide New wallet popup
    hideWalletPopup(){
        this.NewWalletIsVisible = false;
    }

    //Add new wallet function for generate new wallet in Skycoin
    createNewWallet(){
        //Set http headers
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');

        //Post method executed
        this.http.post('/wallet/create', JSON.stringify({name: ''}), {headers: headers})
            .map((res:Response) => res.json())
            .subscribe(
                response => {
                    //Hide new wallet popup
                    this.NewWalletIsVisible = false;
                    alert("New wallet created successfully");
                    //Load wallet for refresh list
                    this.loadWallet();
                },
                err => console.log("Error on create new wallet: "+err),
                () => console.log('New wallet create done')
            );
    }

    //Edit existing wallet function
    editWallet(wallet){
        this.EditWalletIsVisible = true;
        this.walletId = wallet.meta.filename;
    }
    //Hide edit wallet function
    hideEditWalletPopup(){
        this.EditWalletIsVisible = false;
    }

    //Update wallet function for update wallet label
    updateWallet(walletid, walletName){
        //Set http headers
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var stringConvert = 'name='+walletName+'&id='+walletid;
        //Post method executed
        this.http.post('/wallet/update', stringConvert, {headers: headers})
            .map((res:Response) => res.json())
            .subscribe(
                response => {
                    //Hide new wallet popup
                    this.EditWalletIsVisible = false;
                    alert("Wallet updated successfully");
                    //Load wallet for refresh list
                    this.loadWallet();
                },
                err => console.log("Error on update wallet: "+JSON.stringify(err)),
                () => console.log('Update wallet done')
            );
    }

    //Load wallet seed function
    openLoadWallet(walletName, seed){
        this.loadSeedIsVisible = true;
    }
    //Hide load wallet seed function
    hideLoadSeedWalletPopup(){
        this.loadSeedIsVisible = false;
    }

    //Load wallet seed function for create new wallet with name and seed
    createWalletSeed(walletName, seed){
        //Set http headers
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var stringConvert = 'name='+walletName+'&seed='+seed;
        //Post method executed
        this.http.post('/wallet/create', stringConvert, {headers: headers})
            .map((res:Response) => res.json())
            .subscribe(
                response => {
                    //Hide load wallet seed popup
                    this.loadSeedIsVisible = false;
                    //Load wallet for refresh list
                    this.loadWallet();
                },
                err => console.log("Error on create load wallet seed: "+JSON.stringify(err)),
                () => console.log('Load wallet seed done')
            );
    }

    spend(spendid, spendaddress, spendamount){
        //Set local storage for history
        if(localStorage.getItem('historyTable') != null){
            this.historyTable = JSON.parse(localStorage.getItem('historyTable'));
        } else {
            localStorage.setItem('historyTable',JSON.stringify([]));
            this.historyTable = JSON.parse(localStorage.getItem('historyTable'));
        }

        this.historyTable.push({address:spendaddress, amount:spendamount});
        localStorage.setItem('historyTable',JSON.stringify(this.historyTable));

        //Set local storage for addresses history
        if(localStorage.getItem('historyAddresses') != null){
            this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
        } else {
            localStorage.setItem('historyAddresses',JSON.stringify([]));
            this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
        }

        this.addresses.push({address:spendaddress, amount:spendamount});
        localStorage.setItem('historyAddresses',JSON.stringify(this.addresses));


        this.readyDisable = true;
        this.sendDisable = true;
        //Set http headers
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var stringConvert = 'id='+spendid+'&coins='+spendamount*1000000+"&fee=1&hours=1&dst="+spendaddress;
        //Post method executed
        this.http.post('/wallet/spend', stringConvert, {headers: headers})
            .map((res:Response) => res.json())
            .subscribe(
                response => {
                    this.pendingTable.push({complete: 'Completed', address: spendaddress, amount: spendamount});
                    //Load wallet for refresh list
                    this.loadWallet();
                },
                err => {
                    alert(err._body);
                    this.readyDisable = false;
                    this.sendDisable = true;
                    this.pendingTable.push({complete: 'Pending', address: spendaddress, amount: spendamount});
                },
                () => console.log('Spend successfully')
            );
    }

}

//Set default enum value for tabs
enum DisplayModeEnum {
    first = 0,
    second = 1,
    third = 2,
    fourth = 3,
    fifth = 4
}