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
import {QRCodeComponent} from './ng2-qrcode.ts';

declare var _: any;
declare var $: any;
declare var moment: any;

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
    selectedMenu: string;

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
    defaultConnections: Array<any>;
    blockChain: any;
    outputs: Array<any>;
    NewDefaultConnectionIsVisible : boolean;
    EditDefaultConnectionIsVisible : boolean;
    oldConnection:string;
    filterAddressVal:string;
    totalSky:any;
    randomWords:any;

    sortDir:{};

    //Constructor method for load HTTP object
    constructor(private http: Http) { }

    //Init function for load default value
    ngOnInit() {
        this.displayMode = DisplayModeEnum.first;
        this.totalSky = 0;
        this.loadWallet();
        this.loadConnections();
        this.loadDefaultConnections();
        this.loadBlockChain();
        this.loadProgress();
        this.loadOutputs();

        //Set interval function for load wallet every 15 seconds
        setInterval(() => {
            this.loadWallet();
            //console.log("Refreshing balance");
        }, 15000);
        setInterval(() => {
            this.loadConnections();
            this.loadBlockChain();
            //console.log("Refreshing connections");
        }, 5000);

        //Enable Send tab "textbox" and "Ready" button by default
        this.sendDisable = true;
        this.readyDisable = false;
        this.pendingTable = [];
        this.selectedMenu = "Wallets";
        this.sortDir = {time:0, amount:0, address:0};
        this.filterAddressVal = '';

        if(localStorage.getItem('historyAddresses') != null){
            this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
        } else {
            localStorage.setItem('historyAddresses',JSON.stringify([]));
            this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
        }

        //Set local storage for history
        if(localStorage.getItem('historyTable') != null){
            this.historyTable = JSON.parse(localStorage.getItem('historyTable'));
        } else {
            localStorage.setItem('historyTable',JSON.stringify([]));
            this.historyTable = JSON.parse(localStorage.getItem('historyTable'));
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
        this.totalSky = 0;
        this.http.post('/wallets', '')
            .map((res:Response) => res.json())
            .subscribe(
                data => {
                    this.wallets = data;
                    //console.log(this.wallets);

                    //Load Balance for each wallet
                    var inc = 0;
                    for(var item in data){
                        var filename = data[inc].meta.filename;
                        this.loadWalletItem(filename, inc);
                        inc++;
                    }
                    //Load Balance for each wallet end

                },
                err => console.log("Error on load wallet: "+err),
                () => {
                  //console.log('Wallet load done')
                }
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
                    //console.log('load done: ' + inc, response);
                    this.wallets[inc].balance = response.confirmed.coins / 1000000;
                    this.totalSky += this.wallets[inc].balance;
                }, err => console.log("Error on load balance: " + err), () => {
                  //console.log('Balance load done')
                })
        //get address balances
        this.wallets[inc].entries.map((entry)=>{
          this.http.get('/balance?addrs=' + entry.address, { headers: headers })
              .map((res) => res.json())
              .subscribe(
                  //Response from API
                  response => {
                      //console.log('balance:' + entry.address, response);
                      entry.balance = response.confirmed.coins / 1000000;
                  }, err => console.log("Error on load balance: " + err), () => {
                    //console.log('Balance load done')
                  })
        })
    }
    loadConnections() {
        this.http.post('/network/connections', '')
            .map((res) => res.json())
            .subscribe(data => {
                //console.log("connections", data);
                this.connections = data.connections;
            }, err => console.log("Error on load connection: " + err), () => {
              //console.log('Connection load done')
            });
    }
    loadDefaultConnections() {
        this.http.post('/network/defaultConnections', '')
            .map((res) => res.json())
            .subscribe(data => {
                //console.log("default connections", data);
                this.defaultConnections = data;
            }, err => console.log("Error on load default connection: " + err), () => {
              //console.log('Default connections load done')
            });
    }
    loadOutputs() {
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        this.http.get('/outputs', { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
                this.outputs = _.sortBy(data, function(o){
                    return o.address;
                });
            }, err => console.log("Error on load outputs: " + err), () => {
              //console.log('Connection load done')
            });
    }
    loadBlockChain() {
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        this.http.get('/blocks?start=1&end=10000', { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
                //console.log("blockchain", data);
                this.blockChain = data.blocks;
            }, err => console.log("Error on load blockchain: " + err), () => {
              //console.log('blockchain load done');
            });
    }    //Load progress function for Skycoin
    loadProgress(){
        //Post method executed
        this.http.post('/blockchain/progress', '')
            .map((res:Response) => res.json())
            .subscribe(
                //Response from API
                response => { this.progress = (parseInt(response.current,10)+1) / parseInt(response.highest,10) * 100 },
                err => console.log("Error on load progress: "+err),
                () => {
                  //console.log('Progress load done:' + this.progress)
                }
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
    selectMenu(menu, event) {
        this.displayMode = this.displayModeEnum.fifth;
        event.preventDefault();
        this.selectedMenu = menu;
    }
    getDateTimeString(ts) {
        return moment.unix(ts).format("YYYY-MM-DD HH:mm")
    }
    getElapsedTime(ts) {
        return moment().unix() - ts;
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
        this.randomWords = this.getRandomWords();
    }
    //Hide wallet function for hide New wallet popup
    hideWalletPopup(){
        this.NewWalletIsVisible = false;
    }
    showNewDefaultConnectionDialog(){
        this.NewDefaultConnectionIsVisible = true;
    }
    hideNewDefaultConnectionDialog(){
        this.NewDefaultConnectionIsVisible = false;
    }
    showEditDefaultConnectionDialog(item){
        this.oldConnection = item;
        this.EditDefaultConnectionIsVisible = true;
    }
    hideEditDefaultConnectionDialog(){
        this.EditDefaultConnectionIsVisible = false;
    }
    createDefaultConnection(connectionValue){
        //console.log("new value", connectionValue);
        this.defaultConnections.push(connectionValue);
        this.NewDefaultConnectionIsVisible = false;
    }
    updateDefaultConnection(connectionValue){
        //console.log("old/new value", this.oldConnection, connectionValue);
        var idx = this.defaultConnections.indexOf(this.oldConnection);
        this.defaultConnections.splice(idx, 1);
        this.defaultConnections.splice(idx, 0, connectionValue);
        this.EditDefaultConnectionIsVisible = false;
    }
    deleteDefaultConnection(item){
        var idx = this.defaultConnections.indexOf(item);
        this.defaultConnections.splice(idx, 1);
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
                () => {
                  //console.log('New wallet create done')
                }
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
                () => {
                  //console.log('Update wallet done')
                }
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
                () => {
                  //console.log('Load wallet seed done')
                }
            );
    }

    sortHistory(key) {

      if(this.sortDir[key]==0)
        this.sortDir[key] = 1;
      else
        this.sortDir[key] = this.sortDir[key] * (-1);

      if(key == 'time'){
        this.sortDir['address'] = 0;
        this.sortDir['amount'] = 0;
      } else if(key == 'amount') {
        this.sortDir['time'] = 0;
        this.sortDir['address'] = 0;
      } else {
        this.sortDir['time'] = 0;
        this.sortDir['amount'] = 0;
      }

      var self = this;
      if(key != 'address') {
        this.historyTable = _.sortBy(this.historyTable, function(o){
          return Number(o[key]) * self.sortDir[key];
        });
      } else {
        this.historyTable = _.sortBy(this.historyTable, function(o){
          return o[key];
        });

        if(this.sortDir[key] == -1) {
          this.historyTable = this.historyTable.reverse();
        }
      }
    }

    filterHistory(address) {
      console.log("filterHistory", address);
      this.filterAddressVal = address;
    }

    spend(spendid, spendaddress, spendamount){
        var amount = Number(spendamount);
        if(amount < 1) {
          alert('Cannot send values less than 1.');
          return;
        }

        this.historyTable.push({address:spendaddress, amount:spendamount, time:Date.now()/1000});
        localStorage.setItem('historyTable',JSON.stringify(this.historyTable));

        var oldItem = _.find(this.addresses, function(o){
          return o.address === spendaddress;
        })

        if(!oldItem) {
          this.addresses.push({address:spendaddress, amount:spendamount});
          localStorage.setItem('historyAddresses',JSON.stringify(this.addresses));
        }

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
                    //console.log(response);
                    this.pendingTable.push({complete: 'Completed', address: spendaddress, amount: spendamount});
                    //Load wallet for refresh list
                    this.loadWallet();
                    this.readyDisable = false;
                    this.sendDisable = true;
                },
                err => {
                    this.readyDisable = false;
                    this.sendDisable = true;
                    var logBody = err._body;
                    if(logBody == 'Invalid "coins" value') {
                      alert('Incorrect amount value.');
                      return;
                    } else if (logBody == 'Invalid connection') {
                      alert(logBody);
                      return;
                    } else {
                      var logContent = JSON.parse(logBody.substring(logBody.indexOf("{")));
                      alert(logContent.error);
                    }

                    this.pendingTable.push({complete: 'Pending', address: spendaddress, amount: spendamount});
                },
                () => {
                  //console.log('Spend successfully')
                  $("#send_pay_to").val("");
                  $("#send_amount").val(0);
                }
            );
    }

    getRandomWords() {
      var ret = [];
      for(var i = 0 ; i < 11; i++) {
        var length = Math.round(Math.random() * 10);
        length = Math.max(length, 3);

        ret.push(this.createRandomWord(length));
      }

      return ret.join(" ");
    }

    createRandomWord(length) {
    var consonants = 'bcdfghjklmnpqrstvwxyz',
        vowels = 'aeiou',
        rand = function(limit) {
            return Math.floor(Math.random()*limit);
        },
        i, word='',
        consonants2 = consonants.split(''),
        vowels2 = vowels.split('');
    for (i=0;i<length/2;i++) {
        var randConsonant = consonants2[rand(consonants.length)],
            randVowel = vowels2[rand(vowels.length)];
        word += (i===0) ? randConsonant.toUpperCase() : randConsonant;
        word += i*2<length-1 ? randVowel : '';
    }
    return word;
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
