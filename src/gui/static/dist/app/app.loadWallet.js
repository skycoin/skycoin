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
    var PagerService, loadWalletComponent, DisplayModeEnum;
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
            class PagerService {
                getPager(totalItems, currentPage = 1, pageSize = 5) {
                    // calculate total pages
                    var totalPages = Math.ceil(totalItems / pageSize);
                    var startPage, endPage;
                    if (totalPages <= 10) {
                        // less than 10 total pages so show all
                        startPage = 1;
                        endPage = totalPages;
                    }
                    else {
                        // more than 10 total pages so calculate start and end pages
                        if (currentPage <= 6) {
                            startPage = 1;
                            endPage = 10;
                        }
                        else if (currentPage + 4 >= totalPages) {
                            startPage = totalPages - 9;
                            endPage = totalPages;
                        }
                        else {
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
            exports_1("PagerService", PagerService);
            let loadWalletComponent = class loadWalletComponent {
                //Constructor method for load HTTP object
                constructor(http, pagerService) {
                    this.http = http;
                    this.pagerService = pagerService;
                    this.displayModeEnum = DisplayModeEnum;
                    // pager object
                    this.historyPager = {};
                    this.blockPager = {};
                }
                //Init function for load default value
                ngOnInit() {
                    this.displayMode = DisplayModeEnum.first;
                    this.totalSky = 0;
                    this.selectedWallet = {};
                    this.loadWallet();
                    this.loadConnections();
                    this.loadDefaultConnections();
                    this.loadBlockChain();
                    this.loadProgress();
                    this.loadOutputs();
                    this.loadTransactions();
                    this.isValidAddress = false;
                    //Set interval function for load wallet every 15 seconds
                    setInterval(() => {
                        this.loadWallet();
                        //console.log("Refreshing balance");
                    }, 30000);
                    setInterval(() => {
                        this.loadConnections();
                        //this.loadBlockChain();
                        //console.log("Refreshing connections");
                    }, 15000);
                    //Enable Send tab "textbox" and "Ready" button by default
                    this.sendDisable = true;
                    this.readyDisable = false;
                    this.pendingTable = [];
                    this.selectedMenu = "Wallets";
                    this.sortDir = { time: 0, amount: 0, address: 0 };
                    this.filterAddressVal = '';
                    this.historySearchKey = '';
                    if (localStorage.getItem('historyAddresses') != null) {
                        this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
                    }
                    else {
                        localStorage.setItem('historyAddresses', JSON.stringify([]));
                        this.addresses = JSON.parse(localStorage.getItem('historyAddresses'));
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
                    this.totalSky = 0;
                    this.http.post('/wallets', '')
                        .map((res) => res.json())
                        .subscribe(data => {
                        if (this.wallets == null || this.wallets.length == 0) {
                            _.each(data, (o) => {
                                o.showChild = false;
                            });
                            this.wallets = data;
                            if (this.wallets.length > 0) {
                                this.onSelectWallet(this.wallets[0].meta.filename);
                            }
                        }
                        else {
                            data.map((w) => {
                                var old = _.find(this.wallets, (o) => {
                                    return o.meta.filename === w.meta.filename;
                                });
                                if (old) {
                                    _.extend(old, w);
                                }
                                else {
                                    w.showChild = false;
                                    this.wallets.push(w);
                                }
                            });
                        }
                        //console.log("this.wallets", this.wallets);
                        //Load Balance for each wallet
                        var inc = 0;
                        for (var item in data) {
                            var filename = data[inc].meta.filename;
                            this.loadWalletItem(filename, inc);
                            inc++;
                        }
                        //Load Balance for each wallet end
                    }, err => console.log("Error on load wallet: " + err), () => {
                        //console.log('Wallet load done')
                    });
                }
                checkValidAddress(address) {
                    if (address === "")
                        this.isValidAddress = false;
                    else {
                        var headers = new http_2.Headers();
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
                        });
                    }
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
                        this.totalSky += this.wallets[inc].balance;
                    }, err => console.log("Error on load balance: " + err), () => {
                        //console.log('Balance load done')
                    });
                    //get address balances
                    this.wallets[inc].entries.map((entry) => {
                        this.http.get('/balance?addrs=' + entry.address, { headers: headers })
                            .map((res) => res.json())
                            .subscribe(
                        //Response from API
                        response => {
                            //console.log('balance:' + entry.address, response);
                            entry.balance = response.confirmed.coins / 1000000;
                        }, err => console.log("Error on load balance: " + err), () => {
                            //console.log('Balance load done')
                        });
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
                    _.each(transaction.txn.outputs, function (o) {
                        ret += Number(o.coins);
                    });
                    return ret;
                }
                GetTransactionAmount2(transaction) {
                    var ret = 0;
                    _.each(transaction.outputs, function (o) {
                        ret += Number(o.coins);
                    });
                    return ret;
                }
                GetBlockAmount(block) {
                    var ret = [];
                    _.each(block.body.txns[0].outputs, function (o) {
                        ret.push(o.coins);
                    });
                    return ret.join(",");
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
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    this.http.get('/last_blocks?num=10', { headers: headers })
                        .map((res) => res.json())
                        .subscribe(data => {
                        console.log("blockchain", data);
                        this.blockChain = _.sortBy(data.blocks, function (o) {
                            return o.header.seq * (-1);
                        });
                        this.setBlockPage(1);
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
                toggleShowChild(wallet) {
                    wallet.showChild = !wallet.showChild;
                }
                //Switch tab function
                switchTab(mode, wallet) {
                    //"Textbox" and "Ready" button enable in Send tab while switching tabs
                    this.sendDisable = true;
                    this.readyDisable = false;
                    this.displayMode = mode;
                    if (wallet) {
                        this.spendid = wallet.meta.filename;
                        this.selectedWallet = _.find(this.wallets, function (o) {
                            return o.meta.filename === wallet.meta.filename;
                        });
                    }
                }
                selectMenu(menu, event) {
                    this.displayMode = this.displayModeEnum.fifth;
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
                    this.randomWords = this.getRandomWords();
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
                createNewWallet(label, seed) {
                    /*if(addressCount < 1) {
                      alert("Please input correct address count");
                      return;
                    }*/
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    //Post method executed
                    var stringConvert = 'label=' + label + '&seed=' + seed;
                    this.http.post('/wallet/create', stringConvert, { headers: headers })
                        .map((res) => res.json())
                        .subscribe(response => {
                        console.log(response);
                        //Hide new wallet popup
                        this.NewWalletIsVisible = false;
                        alert("New wallet created successfully");
                        //Load wallet for refresh list
                        this.loadWallet();
                    }, err => {
                        console.log(err);
                    }, () => { });
                }
                //Edit existing wallet function
                editWallet(wallet) {
                    this.EditWalletIsVisible = true;
                    this.walletId = wallet.meta.filename;
                }
                addNewAddress(wallet) {
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    //Post method executed
                    var stringConvert = 'id=' + wallet.meta.filename;
                    this.http.post('/wallet/newAddress', stringConvert, { headers: headers })
                        .map((res) => res.json())
                        .subscribe(response => {
                        console.log(response);
                        alert("New address created successfully");
                        //Load wallet for refresh list
                        this.loadWallet();
                    }, err => {
                        console.log(err);
                    }, () => { });
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
                    //this.historyTable.push({address:spendaddress, amount:spendamount, time:Date.now()/1000});
                    //localStorage.setItem('historyTable',JSON.stringify(this.historyTable));
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
                        console.log(response);
                        response.txn.time = Date.now() / 1000;
                        response.txn.address = spendaddress;
                        response.txn.amount = spendamount;
                        this.pendingTable.push(response);
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
                        //this.pendingTable.push({complete: 'Pending', address: spendaddress, amount: spendamount});
                    }, () => {
                        //console.log('Spend successfully')
                        $("#send_pay_to").val("");
                        $("#send_amount").val(0);
                    });
                }
                setHistoryPage(page) {
                    this.historyPager.totalPages = this.historyTable.length;
                    if (page < 1 || page > this.historyPager.totalPages) {
                        return;
                    }
                    // get pager object from service
                    this.historyPager = this.pagerService.getPager(this.historyTable.length, page);
                    // get current page of items
                    this.historyPagedItems = this.historyTable.slice(this.historyPager.startIndex, this.historyPager.endIndex + 1);
                    //console.log('this.pagedItems', this.historyTable, this.pagedItems);
                }
                setBlockPage(page) {
                    this.blockPager.totalPages = this.blockChain.length;
                    if (page < 1 || page > this.blockPager.totalPages) {
                        return;
                    }
                    // get pager object from service
                    this.blockPager = this.pagerService.getPager(this.blockChain.length, page);
                    // get current page of items
                    this.blockPagedItems = this.blockChain.slice(this.blockPager.startIndex, this.blockPager.endIndex + 1);
                    //console.log("this.blockPagedItems", this.blockPagedItems);
                }
                searchHistory(searchKey) {
                    console.log(searchKey);
                }
                getRandomWords() {
                    var ret = [];
                    for (var i = 0; i < 11; i++) {
                        var length = Math.round(Math.random() * 10);
                        length = Math.max(length, 3);
                        ret.push(this.createRandomWord(length));
                    }
                    return ret.join(" ");
                }
                createRandomWord(length) {
                    var consonants = 'bcdfghjklmnpqrstvwxyz', vowels = 'aeiou', rand = function (limit) {
                        return Math.floor(Math.random() * limit);
                    }, i, word = '', consonants2 = consonants.split(''), vowels2 = vowels.split('');
                    for (i = 0; i < length / 2; i++) {
                        var randConsonant = consonants2[rand(consonants.length)], randVowel = vowels2[rand(vowels.length)];
                        word += (i === 0) ? randConsonant.toUpperCase() : randConsonant;
                        word += i * 2 < length - 1 ? randVowel : '';
                    }
                    return word;
                }
                onSelectWallet(val) {
                    console.log("onSelectWallet", val);
                    //this.selectedWallet = val;
                    this.spendid = val;
                    this.selectedWallet = _.find(this.wallets, function (o) {
                        return o.meta.filename === val;
                    });
                }
            };
            loadWalletComponent = __decorate([
                core_1.Component({
                    selector: 'load-wallet',
                    directives: [router_1.ROUTER_DIRECTIVES, ng2_qrcode_ts_1.QRCodeComponent],
                    providers: [PagerService],
                    templateUrl: 'app/templates/wallet.html'
                }), 
                __metadata('design:paramtypes', [http_1.Http, PagerService])
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
