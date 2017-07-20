import {Component, OnInit, ViewChild} from "@angular/core";
import {ROUTER_DIRECTIVES} from "@angular/router";
import {Http, Response, Headers} from "@angular/http";
import {Observable as ObservableRx} from "rxjs/Rx";
import "rxjs/add/operator/map";
import "rxjs/add/operator/catch";
import {QRCodeComponent} from "./ng2-qrcode";
import {SkyCoinEditComponent} from "./components/skycoin.edit.component";
import {SeedComponent} from "./components/seed.component";
import {SkyCoinOutputComponent} from "./components/outputs.component";
import {PendingTxnsComponent} from "./components/pending.transactions.component";
import {WalletBackupPageComponent} from "./components/wallet.backup.page.component";
import {SkycoinSyncWalletBlock} from "./components/progress.bannner.component";

declare var _: any;
declare var $: any;
declare var async: any;
declare var moment: any;
declare var toastr: any;

export class PagerService {
    getPager(totalItems: number, currentPage: number = 1, pageSize: number = 5) {
        // calculate total pages
        var totalPages = Math.ceil(totalItems / pageSize);

        var startPage, endPage;
        if (totalPages <= 10) {
            // less than 10 total pages so show all
            startPage = 1;
            endPage = totalPages;
        } else {
            // more than 10 total pages so calculate start and end pages
            if (currentPage <= 6) {
                startPage = 1;
                endPage = 10;
            } else if (currentPage + 4 >= totalPages) {
                startPage = totalPages - 9;
                endPage = totalPages;
            } else {
                startPage = currentPage - 5;
                endPage = currentPage + 4;
            }
        }

        // calculate start and end item indexes
        var startIndex = (currentPage - 1) * pageSize;
        var endIndex = Math.min(startIndex + pageSize - 1, totalItems - 1);

        // create an array of pages to ng-repeat in the pager control
        var pages = _.range(startPage, endPage + 1);

        // return object with all pager properties required by the view
        return {
            totalItems: totalItems,
            currentPage: currentPage,
            pageSize: pageSize,
            totalPages: totalPages,
            startPage: startPage,
            endPage: endPage,
            startIndex: startIndex,
            endIndex: endIndex,
            pages: pages
        };
    }
}

@Component({
    selector: 'load-wallet',
    directives: [ROUTER_DIRECTIVES, QRCodeComponent, SeedComponent, SkycoinSyncWalletBlock, SkyCoinEditComponent, SkyCoinOutputComponent, PendingTxnsComponent, WalletBackupPageComponent],
    providers: [PagerService],
    templateUrl: 'app/templates/wallet.html'
})

export class LoadWalletComponent implements OnInit {
    //Declare default varialbes
    wallets : Array<any>;
    walletsWithAddress : Array<any>;
    progress: any;
    spendid: string;
    spendaddress: string;
    sendDisable: boolean;
    readyDisable: boolean;
    displayMode: DisplayModeEnum;
    displayModeEnum = DisplayModeEnum;
    selectedMenu: string;
    historyTimeBy:string = 'desc';
    userTransactions:Array<any>;

    @ViewChild(SkyCoinOutputComponent)
    private outputComponent: SkyCoinOutputComponent;

    @ViewChild(PendingTxnsComponent)
    private pendingTxnComponent: PendingTxnsComponent;

    @ViewChild('spendaddress')
    private spendAddress:any;

    @ViewChild('walletname')
    private walletNameInLoadWallet:any;

    @ViewChild('walletseed')
    private walletSeedInLoadWallet:any;

    @ViewChild('spendamount')
    private spendAmount:any;

    @ViewChild('transactionNote')
    private transactionNote:any;

    QrAddress: string;
    QrIsVisible: boolean;

    @ViewChild(SeedComponent)
    private seedComponent: SeedComponent;

    NewWalletIsVisible: boolean;
    loadSeedIsVisible: boolean;

    walletname: string;
    walletId: string;

    historyTable: Array<any>;
    pendingTable: Array<any>;
    addresses: Array<any>;
    connections: Array<any>;
    defaultConnections: Array<any>;
    blockChain: any;
    elapsedTime: number;
    numberOfBlocks: number;
    outputs: Array<any>;
    NewDefaultConnectionIsVisible : boolean;
    EditDefaultConnectionIsVisible : boolean;
    noConnections:boolean;
    oldConnection:string;
    filterAddressVal:string;
    totalSky:any;
    totalWalletBalance:any;
    historySearchKey:string;
    selectedWallet:any;

    sortDir:{};
    isValidAddress: boolean;

    // pager object
    historyPager: any = {};
    historyPagedItems: any[];

    blockPager: any = {};

    //Constructor method for load HTTP object
    constructor(private http: Http, private pagerService: PagerService) { }

    //Init function for load default value
    ngOnInit() {
        this.displayMode = DisplayModeEnum.first;
        this.totalSky = 0;
        this.totalWalletBalance =0;
        this.selectedWallet = {};
        this.userTransactions=[];
        this.loadWallet();
        this.loadConnections();
        this.loadDefaultConnections();
        this.loadBlockChain();
        this.loadNumberOfBlocks();
        this.loadProgress();
        this.isValidAddress = false;

        //Set interval function for load wallet every 15 seconds
        setInterval(() => {
            this.loadWallet();
        }, 30000);
        setInterval(() => {
            this.loadConnections();
            this.loadBlockChain();
            this.loadNumberOfBlocks();
            //console.log("Refreshing connections");
        }, 15000);

        //Enable Send tab "textbox" and "Ready" button by default
        this.sendDisable = true;
        this.readyDisable = false;
        this.pendingTable = [];
        this.selectedMenu = "Wallets";
        this.sortDir = {time:0, amount:0, address:0};
        this.filterAddressVal = '';
        this.historySearchKey = '';

        if(localStorage.getItem('historyAddresses') != null){
            this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
        } else {
            localStorage.setItem('historyAddresses',JSON.stringify([]));
            this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
        }

        /*$("#walletSelect").select2({
         templateResult: function(state) {
         return state.text;
         /!*if (!state.id) { return state.text; }
         var $state = $(
         '<span><img src="vendor/images/flags/' + state.element.value.toLowerCase() + '.png" class="img-flag" /> ' + state.text + '</span>'
         );
         return $state;*!/
         }
         });*/
    }

    //Ready button function for disable "textbox" and enable "Send" button for ready to send coin
    ready(spendId, spendaddress, spendamount){
        if(!spendId){
            toastr.error("Please select from id");
            return false;
        }
        if(!spendaddress){
            toastr.error("Please enter pay to");
            return false;
        }
        if(!spendamount){
            toastr.error("Please enter amount");
            return false;
        }
        this.readyDisable = true;
        this.sendDisable = false;
    }

    loadNumberOfBlocks(){
        this.numberOfBlocks=0;
        this.http.get('/blockchain/metadata')
        .map((res:Response)=>res.json())
        .subscribe(
            data=>{
                this.numberOfBlocks = data.head.seq;
            }
        )
    }

    loadTransactionsForWallet(){
        let addresses=[];

        this.userTransactions=[];

        _.each(this.wallets,(wallet)=>{
            _.each(wallet.entries,(entry)=>{
                addresses.push(entry.address);
            });
        });


        _.each(addresses,(address)=>{
            this.http.get('/explorer/address?address='+address, {})
            .map((res) => res.json())
            .subscribe(transactions => {
                _.each(transactions,(transaction)=>{
                    let confirmed = transaction.status.confirmed? 'Confirmed': 'UnConfirmed'
                    this.userTransactions.push({'type':'confirmed','transactionInputs':transaction.inputs,'transactionOutputs':transaction.outputs
                        ,'actualTransaction':transaction, 'confirmed': confirmed
                    });
                    this.sortTransactions();
                });
            });
        });
    }
    sortHistoryByTime(ev:Event) {
        this.historyTimeBy = this.historyTimeBy === 'desc' ? 'asc' : 'desc';
        this.sortTransactions()
    }
    sortTransactions() {
        if(this.userTransactions.length <= 0) {
            return;
        }
        switch (this.historyTimeBy) {
            case 'asc':
                this.userTransactions.sort((a, b) => {
                    if (a.actualTransaction.timestamp > b.actualTransaction.timestamp) {
                        return 1;
                    }
                    if (a.actualTransaction.timestamp < b.actualTransaction.timestamp) {
                        return -1;
                    }
                    return 0;
                });
                break;
            case 'desc':
                 this.userTransactions.sort((a, b) => {
                    if (a.actualTransaction.timestamp > b.actualTransaction.timestamp) {
                        return -1;
                    }
                    if (a.actualTransaction.timestamp < b.actualTransaction.timestamp) {
                        return 1;
                    }
                    return 0;
                });
                break;
        }
    }

    //Load wallet function
    loadWallet(){
        this.totalWalletBalance = 0;
        this.http.post('/wallets', '')
        .map((res:Response) => res.json())
        .subscribe(
            data => {
                if(this.wallets == null || this.wallets.length == 0) {
                    _.each(data, (o)=>{
                        o.showChild = false;
                    })
                    this.wallets = data;
                    if (this.wallets.length > 0) {
                        this.onSelectWallet(this.wallets[0].meta.filename);
                    }
                } else {
                    data.map((w)=>{
                        var old = _.find(this.wallets, (o)=>{
                            return o.meta.filename === w.meta.filename;
                        })

                        if(old) {
                            _.extend(old, w);
                        } else {
                            w.showChild = false;
                            this.wallets.push(w);
                        }
                    })
                }

                //console.log("this.wallets", this.wallets);

                //Load Balance for each wallet
                //var inc = 0;
                //console.log("data", data);
                _.map(data, (item, idx) => {
                    var filename = item.meta.filename;
                    this.loadWalletItem(filename, idx,idx ==  data.length-1);
                })
                this.walletsWithAddress = [];
                _.map(this.wallets, (o, idx) => {
                    this.walletsWithAddress.push({
                        wallet:o,
                        type:'wallet'
                    });
                    _.map(o.entries, (_o, idx) => {
                        this.walletsWithAddress.push({
                            entry:_o,
                            type:'address',
                            wallet:o,
                            idx:idx==0?'':'(' + idx + ')'
                        });
                    });
                });
                this.loadTransactionsForWallet();
            },
            err => console.log(err),
            () => {
                //console.log('Wallet load done')
            }
        );
    }
    checkValidAddress(address) {
        if(address === "") {
            this.isValidAddress = false;
        } else {
            var headers = new Headers();
            headers.append('Content-Type', 'application/x-www-form-urlencoded');
            this.http.get('/balance?addrs=' + address, { headers: headers })
            .map((res) => res.json())
            .subscribe(
                //Response from API
                response => {
                    this.isValidAddress = true;
                }, err => {
                    //console.log("Error on load balance: " + err)
                    this.isValidAddress = false;
                }, () => {

                })
        }
    }
    loadWalletItem(address, inc, endReached){
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
                this.totalWalletBalance += this.wallets[inc].balance;
                if(endReached){
                   if(this.totalSky != this.totalWalletBalance) {
                       this.totalSky = this.totalWalletBalance;
                       console.log("Updating the wallet balance!");
                   }
                   else{
                       console.log("Not updating the wallet balance as it is the same!");
                   }
                }
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
    loadTransactions() {
        this.historyTable = [];
        this.http.get('/lastTxs', {})
        .map((res) => res.json())
        .subscribe(data => {
            console.log("transactions", data);
            this.historyTable = this.historyTable.concat(data);
            this.setHistoryPage(1);
        }, err => console.log("Error on load transactions: " + err), () => {
            //console.log('Connection load done')
        });
        this.http.get('/pendingTxs', {})
        .map((res) => res.json())
        .subscribe(data => {
            console.log("pending transactions", data);
            this.historyTable = this.historyTable.concat(data);
            this.setHistoryPage(1);
        }, err => console.log("Error on pending transactions: " + err), () => {

        });
    }
    GetTransactionAmount(transaction) {
        var ret = 0;
        _.each(transaction.outputs, function(o){
            ret += Number(o.coins);
        })

        return ret;
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

    loadBlockChain() {
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        this.http.get('/last_blocks?num=10', { headers: headers })
        .map((res) => res.json())
        .subscribe(data => {
            //console.log("blockchain", data);
            this.blockChain = _.sortBy(data.blocks, function(o){
                return o.header.seq * (-1);
            });

            if (this.blockChain.length != 0) {
                this.elapsedTime = moment().unix() - this.blockChain[0].header.timestamp;
            }

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
    toggleShowChild(wallet) {
        wallet.showChild = !wallet.showChild;
    }

    //Switch tab function
    switchTab(mode: DisplayModeEnum, wallet) {
        this.selectedMenu = '';
        //"Textbox" and "Ready" button enable in Send tab while switching tabs
        this.sendDisable = true;
        this.readyDisable = false;

        this.displayMode = mode;
        if(wallet){
            this.spendid = wallet.meta.filename;
            this.selectedWallet = _.find(this.wallets, function(o){
                return o.meta.filename === wallet.meta.filename;
            })
            console.log("selected wallet", this.spendid, this.selectedWallet);
        }
    }
    selectMenu(menu, event) {
        this.displayMode = this.displayModeEnum.fifth;
        event.preventDefault();
        this.selectedMenu = menu;
        if(menu=='Outputs'){
            if(this.outputComponent){
                this.outputComponent.refreshOutputs();
            }
        }
        if(menu == 'PendingTxns'){
            if(this.pendingTxnComponent){
                this.pendingTxnComponent.refreshPendingTxns();
            }
        }
    }
    getDateTimeString(ts) {
        return moment.unix(ts).format("YYYY-MM-DD HH:mm")
    }
    //Show QR code function for show QR popup
    showQR(address){
        this.QrAddress = address;
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
    createNewWallet(label, seed, addressCount){
        if(addressCount < 1) {
            //alert("Please input correct address count");
            toastr.error('Please input correct address count');
            return;
        }

        // new wallet will have a default address
        addressCount--;

        //check if label is duplicated
        var old = _.find(this.wallets, function(o){
            return (o.meta.label == label)
        })

        if(old) {
            toastr.error('This wallet label is used already.');
            //alert("This wallet label is used already");
            return;
        }

        //Set http headers
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        if(this.seedComponent){
            seed=this.seedComponent.getCurrentSeed();
        }
        //Post method executed
        var stringConvert = 'label='+label+'&seed='+seed;
        this.http.post('/wallet/create', stringConvert, {headers: headers})
        .map((res:Response) => res.json())
        .subscribe(
            response => {
                console.log('response:', response);
                if (addressCount > 0) {
                    let param = 'id='+response.meta.filename+'&num='+addressCount
                    this.http.post('/wallet/newAddress', param, {headers: headers})
                    .map((res:Response) => res.json())
                    .subscribe(
                        response => {
                            //Hide new wallet popup
                            this.NewWalletIsVisible = false;
                            toastr.info("New wallet created successfully");
                            //Load wallet for refresh list
                            this.loadWallet();
                        },
                        err => {
                            console.log(err);
                        }
                    );
                    return
                }

                //Hide new wallet popup
                this.NewWalletIsVisible = false;
                toastr.info("New wallet created successfully");
                //Load wallet for refresh list
                this.loadWallet();
            },
            err => {
                // when node is down, the response header status is 200
                if (err.status == 200) {
                    return
                }

                if (err.status == 400) {
                    if(err._body.indexOf("duplicate wallet ") !=-1){
                        toastr.info("Can't load same wallet twice!");
                    }
                }
            },
            () => {}
        );
    }

    addNewAddress(wallet) {
        //Set http headers
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');

        //Post method executed
        var stringConvert = 'id='+wallet.meta.filename;
        this.http.post('/wallet/newAddress', stringConvert, {headers: headers})
        .map((res:Response) => res.json())
        .subscribe(
            response => {
                console.log(response)
                toastr.info("New address created successfully");
                //Load wallet for refresh list
                this.loadWallet();
            },
            err => {
                console.log(err);
            },
            () => {}
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
    createWalletSeed(walletName, seed, addressCount){
        if(addressCount < 1) {
            toastr.error('Please correct address count');
            return;
        }

        addressCount--;

        //Set http headers
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var stringConvert = 'label='+walletName+'&seed='+seed;
        //Post method executed
        this.http.post('/wallet/create', stringConvert, {headers: headers})
        .map((res:Response) => res.json())
        .subscribe(
            response => {
                if(addressCount > 0) {
                    let params='id=' + response.meta.filename + '&num=' + addressCount;
                    this.http.post('/wallet/newAddress', params, {headers: headers})
                    .map((res:Response) => res.json())
                    .subscribe(
                        response => {
                            //Hide new wallet popup
                            this.loadSeedIsVisible = false;
                            toastr.info("Wallet loaded successfully");
                            //Load wallet for refresh list
                            this.loadWallet();
                        },
                        err => {
                            console.log(err);
                        },
                        () => {}
                    );
                    return
                } 

                //Hide new wallet popup
                this.loadSeedIsVisible = false;
                toastr.info("Wallet loaded successfully");
                //Load wallet for refresh list
                this.loadWallet();
            },
            err => {
                // when node is down, the response header status is 200
                if (err.status == 200) {
                    return
                }

                if (err.status == 400) {
                    if(err._body.indexOf("duplicate wallet ") !=-1){
                        toastr.info("Can't load same wallet twice!");
                    }
                }
            },
            () => {
            }
        );
    }

    filterHistory(address) {
        console.log("filterHistory", address);
        this.filterAddressVal = address;
    }

    updateStatusOfTransaction(txid, metaData){

        let self= this;
        let transactionConfirmed = false;

        ObservableRx.timer(0,1000).map((i)=>{
            if(transactionConfirmed){
                throw new Error("Transaction confirmed");
            }

            var headers = new Headers();
            headers.append('Content-Type', 'application/x-www-form-urlencoded');
            this.http.get('/transaction?txid=' + txid, { headers: headers })
            .map((res) => res.json())
            .subscribe(
                res => {
                    transactionConfirmed = res.status.confirmed;
                    self.pendingTable=[];
                    self.pendingTable.push({'time':res.txn.timestamp,'status':res.status.confirmed?'Completed':'Unconfirmed','amount':metaData.amount,'txId':txid,'address':metaData.address});
                    self.loadWallet();
                }, err => {
                    console.log("Error on load transaction: " + err)
                }, () => {
                });
        }).subscribe(()=>{

        },(err)=>{
            console.log("Transaction confirmed");
        });


    }

    spend(spendid, spendaddress, spendamount){
        var amount = Number(spendamount);
        if(amount < 1) {
            toastr.error('Cannot send values less than 1.');
            return;
        }

        //this.historyTable.push({address:spendaddress, amount:spendamount, time:Date.now()/1000});
        //localStorage.setItem('historyTable',JSON.stringify(this.historyTable));

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
        var self = this;
        this.http.post('/wallet/spend', stringConvert, {headers: headers})
        .map((res:Response) => res.json())
        .subscribe(
            response => {
                console.log(response);
                this.updateStatusOfTransaction(response.txn.txid, {address:spendaddress,amount:amount});
                this.readyDisable = false;
                this.sendDisable = true;
                self.spendAddress.nativeElement.value = '';
                self.spendAmount.nativeElement.value =0;
                self.transactionNote.nativeElement.value = '';
                self.isValidAddress=false;
            },
            err => {



                this.readyDisable = false;
                this.sendDisable = true;
                var logBody = err._body;
                if(logBody == 'Invalid "coins" value') {
                    toastr.error('Incorrect amount value.');
                    return;
                } else if (logBody == 'Invalid connection') {
                    toastr.error(logBody);
                    return;
                } else {
                    var logContent = JSON.parse(logBody.substring(logBody.indexOf("{")));
                    toastr.error(logContent.error);
                }

                //this.pendingTable.push({complete: 'Pending', address: spendaddress, amount: spendamount});
            },
            () => {
                self.spendAddress.nativeElement.value = '';
                self.spendAmount.nativeElement.value =0;
                self.transactionNote.nativeElement.value = '';
                self.isValidAddress=false;
                $("#send_pay_to").val("");
                $("#send_amount").val(0);
            }
        );
    }

    setHistoryPage(page: number) {
        this.historyPager.totalPages = this.historyTable.length;

        if (page < 1 || page > this.historyPager.totalPages) {
            return;
        }

        // get pager object from service
        this.historyPager = this.pagerService.getPager(this.historyTable.length, page);

        console.log("this.historyPager", this.historyPager );
        // get current page of items
        this.historyPagedItems = this.historyTable.slice(this.historyPager.startIndex, this.historyPager.endIndex + 1);
        //console.log('this.pagedItems', this.historyTable, this.pagedItems);
    }


    searchHistory(searchKey){
        console.log(searchKey);

    }

    onSelectWallet(val) {
        console.log("onSelectWallet", val);
        //this.selectedWallet = val;
        this.spendid = val;
        this.selectedWallet = _.find(this.wallets, function(o){
            return o.meta.filename === val;
        })
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