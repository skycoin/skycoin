System.register(['angular2/core', 'angular2/router', 'angular2/http', 'rxjs/add/operator/map', 'rxjs/add/operator/catch', './ng2-qrcode.js'], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, router_1, http_1, http_2, ng2_qrcode_ts_1;
    var loadWalletComponent, DisplayModeEnum;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (router_1_1) {
                router_1 = router_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
                http_2 = http_1_1;
            },
            function (_1) {},
            function (_2) {},
            function (ng2_qrcode_ts_1_1) {
                ng2_qrcode_ts_1 = ng2_qrcode_ts_1_1;
            }],
        execute: function() {
            let loadWalletComponent = class loadWalletComponent {
                //Constructor method for load HTTP object
                constructor(http) {
                    this.http = http;
                    this.displayModeEnum = DisplayModeEnum;
                }
                //Init function for load default value
                ngOnInit() {
                    this.displayMode = DisplayModeEnum.first;
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
                    this.sortDir = { time: 0, amount: 0, address: 0 };
                    this.filterAddressVal = '';
                    if (localStorage.getItem('historyAddresses') != null) {
                        this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
                    }
                    else {
                        localStorage.setItem('historyAddresses', JSON.stringify([]));
                        this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
                    }
                    //Set local storage for history
                    if (localStorage.getItem('historyTable') != null) {
                        this.historyTable = JSON.parse(localStorage.getItem('historyTable'));
                    }
                    else {
                        localStorage.setItem('historyTable', JSON.stringify([]));
                        this.historyTable = JSON.parse(localStorage.getItem('historyTable'));
                    }
                }
                //Ready button function for disable "textbox" and enable "Send" button for ready to send coin
                ready(spendId, spendaddress, spendamount) {
                    if (!spendId) {
                        alert("Please select from id");
                        return false;
                    }
                    if (!spendaddress) {
                        alert("Please enter pay to");
                        return false;
                    }
                    if (!spendamount) {
                        alert("Please enter amount");
                        return false;
                    }
                    this.readyDisable = true;
                    this.sendDisable = false;
                }
                //Load wallet function
                loadWallet() {
                    this.http.post('/wallets', '')
                        .map((res) => res.json())
                        .subscribe(data => {
                        this.wallets = data;
                        //Load Balance for each wallet
                        var inc = 0;
                        for (var item in data) {
                            var address = data[inc].meta.filename;
                            this.loadWalletItem(address, inc);
                            inc++;
                        }
                        //Load Balance for each wallet end
                    }, err => console.log("Error on load wallet: " + err), () => {
                        //console.log('Wallet load done')
                    });
                }
                loadWalletItem(address, inc) {
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    this.http.get('/wallet/balance?id=' + address, { headers: headers })
                        .map((res) => res.json())
                        .subscribe(
                    //Response from API
                    response => {
                        //console.log('load done: ' + inc, response);
                        this.wallets[inc].balance = response.confirmed.coins / 1000000;
                    }, err => console.log("Error on load balance: " + err), () => {
                        //console.log('Balance load done')
                    });
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
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    this.http.get('/outputs', { headers: headers })
                        .map((res) => res.json())
                        .subscribe(data => {
                        this.outputs = _.sortBy(data, function (o) {
                            return o.address;
                        });
                    }, err => console.log("Error on load outputs: " + err), () => {
                        //console.log('Connection load done')
                    });
                }
                loadBlockChain() {
                    this.http.post('/blockchain', '')
                        .map((res) => res.json())
                        .subscribe(data => {
                        //console.log("blockchain", data);
                        if (data.head) {
                            this.blockChain = data;
                        }
                    }, err => console.log("Error on load blockchain: " + err), () => {
                        //console.log('blockchain load done');
                    });
                } //Load progress function for Skycoin
                loadProgress() {
                    //Post method executed
                    this.http.post('/blockchain/progress', '')
                        .map((res) => res.json())
                        .subscribe(
                    //Response from API
                    response => { this.progress = (parseInt(response.current, 10) + 1) / parseInt(response.highest, 10) * 100; }, err => console.log("Error on load progress: " + err), () => {
                        //console.log('Progress load done:' + this.progress)
                    });
                }
                //Switch tab function
                switchTab(mode, wallet) {
                    //"Textbox" and "Ready" button enable in Send tab while switching tabs
                    this.sendDisable = true;
                    this.readyDisable = false;
                    this.displayMode = mode;
                    if (wallet) {
                        this.spendid = wallet.meta.filename;
                    }
                }
                selectMenu(menu, event) {
                    this.displayMode = this.displayModeEnum.fourth;
                    event.preventDefault();
                    this.selectedMenu = menu;
                }
                getDateTimeString(ts) {
                    return moment.unix(ts).format("YYYY-MM-DD HH:mm");
                }
                getElapsedTime(ts) {
                    return moment().unix() - ts;
                }
                //Show QR code function for show QR popup
                showQR(wallet) {
                    this.QrAddress = wallet.entries[0].address;
                    this.QrIsVisible = true;
                }
                //Hide QR code function for hide QR popup
                hideQrPopup() {
                    this.QrIsVisible = false;
                }
                //Show wallet function for view New wallet popup
                showNewWalletDialog() {
                    this.NewWalletIsVisible = true;
                }
                //Hide wallet function for hide New wallet popup
                hideWalletPopup() {
                    this.NewWalletIsVisible = false;
                }
                showNewDefaultConnectionDialog() {
                    this.NewDefaultConnectionIsVisible = true;
                }
                hideNewDefaultConnectionDialog() {
                    this.NewDefaultConnectionIsVisible = false;
                }
                showEditDefaultConnectionDialog(item) {
                    this.oldConnection = item;
                    this.EditDefaultConnectionIsVisible = true;
                }
                hideEditDefaultConnectionDialog() {
                    this.EditDefaultConnectionIsVisible = false;
                }
                createDefaultConnection(connectionValue) {
                    //console.log("new value", connectionValue);
                    this.defaultConnections.push(connectionValue);
                    this.NewDefaultConnectionIsVisible = false;
                }
                updateDefaultConnection(connectionValue) {
                    //console.log("old/new value", this.oldConnection, connectionValue);
                    var idx = this.defaultConnections.indexOf(this.oldConnection);
                    this.defaultConnections.splice(idx, 1);
                    this.defaultConnections.splice(idx, 0, connectionValue);
                    this.EditDefaultConnectionIsVisible = false;
                }
                deleteDefaultConnection(item) {
                    var idx = this.defaultConnections.indexOf(item);
                    this.defaultConnections.splice(idx, 1);
                }
                //Add new wallet function for generate new wallet in Skycoin
                createNewWallet() {
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    //Post method executed
                    this.http.post('/wallet/create', JSON.stringify({ name: '' }), { headers: headers })
                        .map((res) => res.json())
                        .subscribe(response => {
                        //Hide new wallet popup
                        this.NewWalletIsVisible = false;
                        alert("New wallet created successfully");
                        //Load wallet for refresh list
                        this.loadWallet();
                    }, err => console.log("Error on create new wallet: " + err), () => {
                        //console.log('New wallet create done')
                    });
                }
                //Edit existing wallet function
                editWallet(wallet) {
                    this.EditWalletIsVisible = true;
                    this.walletId = wallet.meta.filename;
                }
                //Hide edit wallet function
                hideEditWalletPopup() {
                    this.EditWalletIsVisible = false;
                }
                //Update wallet function for update wallet label
                updateWallet(walletid, walletName) {
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    var stringConvert = 'name=' + walletName + '&id=' + walletid;
                    //Post method executed
                    this.http.post('/wallet/update', stringConvert, { headers: headers })
                        .map((res) => res.json())
                        .subscribe(response => {
                        //Hide new wallet popup
                        this.EditWalletIsVisible = false;
                        alert("Wallet updated successfully");
                        //Load wallet for refresh list
                        this.loadWallet();
                    }, err => console.log("Error on update wallet: " + JSON.stringify(err)), () => {
                        //console.log('Update wallet done')
                    });
                }
                //Load wallet seed function
                openLoadWallet(walletName, seed) {
                    this.loadSeedIsVisible = true;
                }
                //Hide load wallet seed function
                hideLoadSeedWalletPopup() {
                    this.loadSeedIsVisible = false;
                }
                //Load wallet seed function for create new wallet with name and seed
                createWalletSeed(walletName, seed) {
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    var stringConvert = 'name=' + walletName + '&seed=' + seed;
                    //Post method executed
                    this.http.post('/wallet/create', stringConvert, { headers: headers })
                        .map((res) => res.json())
                        .subscribe(response => {
                        //Hide load wallet seed popup
                        this.loadSeedIsVisible = false;
                        //Load wallet for refresh list
                        this.loadWallet();
                    }, err => console.log("Error on create load wallet seed: " + JSON.stringify(err)), () => {
                        //console.log('Load wallet seed done')
                    });
                }
                sortHistory(key) {
                    if (this.sortDir[key] == 0)
                        this.sortDir[key] = 1;
                    else
                        this.sortDir[key] = this.sortDir[key] * (-1);
                    if (key == 'time') {
                        this.sortDir['address'] = 0;
                        this.sortDir['amount'] = 0;
                    }
                    else if (key == 'amount') {
                        this.sortDir['time'] = 0;
                        this.sortDir['address'] = 0;
                    }
                    else {
                        this.sortDir['time'] = 0;
                        this.sortDir['amount'] = 0;
                    }
                    var self = this;
                    if (key != 'address') {
                        this.historyTable = _.sortBy(this.historyTable, function (o) {
                            return Number(o[key]) * self.sortDir[key];
                        });
                    }
                    else {
                        this.historyTable = _.sortBy(this.historyTable, function (o) {
                            return o[key];
                        });
                        if (this.sortDir[key] == -1) {
                            this.historyTable = this.historyTable.reverse();
                        }
                    }
                }
                filterHistory(address) {
                    console.log("filterHistory", address);
                    this.filterAddressVal = address;
                }
                spend(spendid, spendaddress, spendamount) {
                    var amount = Number(spendamount);
                    if (amount < 1) {
                        alert('Cannot send values less than 1.');
                        return;
                    }
                    this.historyTable.push({ address: spendaddress, amount: spendamount, time: Date.now() / 1000 });
                    localStorage.setItem('historyTable', JSON.stringify(this.historyTable));
                    var oldItem = _.find(this.addresses, function (o) {
                        return o.address === spendaddress;
                    });
                    if (!oldItem) {
                        this.addresses.push({ address: spendaddress, amount: spendamount });
                        localStorage.setItem('historyAddresses', JSON.stringify(this.addresses));
                    }
                    this.readyDisable = true;
                    this.sendDisable = true;
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    var stringConvert = 'id=' + spendid + '&coins=' + spendamount * 1000000 + "&fee=1&hours=1&dst=" + spendaddress;
                    //Post method executed
                    this.http.post('/wallet/spend', stringConvert, { headers: headers })
                        .map((res) => res.json())
                        .subscribe(response => {
                        //console.log(response);
                        this.pendingTable.push({ complete: 'Completed', address: spendaddress, amount: spendamount });
                        //Load wallet for refresh list
                        this.loadWallet();
                        this.readyDisable = false;
                        this.sendDisable = true;
                    }, err => {
                        this.readyDisable = false;
                        this.sendDisable = true;
                        var logBody = err._body;
                        if (logBody == 'Invalid "coins" value') {
                            alert('Incorrect amount value.');
                            return;
                        }
                        else if (logBody == 'Invalid connection') {
                            alert(logBody);
                            return;
                        }
                        else {
                            var logContent = JSON.parse(logBody.substring(logBody.indexOf("{")));
                            alert(logContent.error);
                        }
                        this.pendingTable.push({ complete: 'Pending', address: spendaddress, amount: spendamount });
                    }, () => {
                        //console.log('Spend successfully')
                        $("#send_pay_to").val("");
                        $("#send_amount").val(0);
                    });
                }
            };
            loadWalletComponent = __decorate([
                core_1.Component({
                    selector: 'load-wallet',
                    directives: [router_1.ROUTER_DIRECTIVES, ng2_qrcode_ts_1.QRCodeComponent],
                    providers: [],
                    templateUrl: 'app/templates/wallet.html'
                }), 
                __metadata('design:paramtypes', [http_1.Http])
            ], loadWalletComponent);
            exports_1("loadWalletComponent", loadWalletComponent);
            //Set default enum value for tabs
            (function (DisplayModeEnum) {
                DisplayModeEnum[DisplayModeEnum["first"] = 0] = "first";
                DisplayModeEnum[DisplayModeEnum["second"] = 1] = "second";
                DisplayModeEnum[DisplayModeEnum["third"] = 2] = "third";
                DisplayModeEnum[DisplayModeEnum["fourth"] = 3] = "fourth";
                DisplayModeEnum[DisplayModeEnum["fifth"] = 4] = "fifth";
            })(DisplayModeEnum || (DisplayModeEnum = {}));
        }
    }
});

//# sourceMappingURL=app.loadWallet.js.map
