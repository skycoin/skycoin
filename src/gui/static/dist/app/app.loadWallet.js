System.register(['@angular/core', '@angular/router', '@angular/http', 'rxjs/Rx', 'rxjs/add/operator/map', 'rxjs/add/operator/catch', './ng2-qrcode', './components/skycoin.edit.component', './components/seed.component', './components/outputs.component', "./components/pending.transactions.component", "./components/wallet.backup.page.component"], function(exports_1, context_1) {
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
    var core_1, router_1, http_1, http_2, Rx_1, ng2_qrcode_1, skycoin_edit_component_1, seed_component_1, outputs_component_1, pending_transactions_component_1, wallet_backup_page_component_1;
    var PagerService, LoadWalletComponent, DisplayModeEnum;
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
            function (Rx_1_1) {
                Rx_1 = Rx_1_1;
            },
            function (_1) {},
            function (_2) {},
            function (ng2_qrcode_1_1) {
                ng2_qrcode_1 = ng2_qrcode_1_1;
            },
            function (skycoin_edit_component_1_1) {
                skycoin_edit_component_1 = skycoin_edit_component_1_1;
            },
            function (seed_component_1_1) {
                seed_component_1 = seed_component_1_1;
            },
            function (outputs_component_1_1) {
                outputs_component_1 = outputs_component_1_1;
            },
            function (pending_transactions_component_1_1) {
                pending_transactions_component_1 = pending_transactions_component_1_1;
            },
            function (wallet_backup_page_component_1_1) {
                wallet_backup_page_component_1 = wallet_backup_page_component_1_1;
            }],
        execute: function() {
            PagerService = (function () {
                function PagerService() {
                }
                PagerService.prototype.getPager = function (totalItems, currentPage, pageSize) {
                    if (currentPage === void 0) { currentPage = 1; }
                    if (pageSize === void 0) { pageSize = 5; }
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
                };
                return PagerService;
            }());
            exports_1("PagerService", PagerService);
            LoadWalletComponent = (function () {
                //Constructor method for load HTTP object
                function LoadWalletComponent(http, pagerService) {
                    this.http = http;
                    this.pagerService = pagerService;
                    this.displayModeEnum = DisplayModeEnum;
                    this.selectedBlock = {};
                    this.selectedBlockTransaction = {};
                    this.selectedBlockAddressBalance = 0;
                    this.selectedBlackAddressTxList = [];
                    // pager object
                    this.historyPager = {};
                    this.blockPager = {};
                }
                //Init function for load default value
                LoadWalletComponent.prototype.ngOnInit = function () {
                    var _this = this;
                    this.displayMode = DisplayModeEnum.first;
                    this.totalSky = 0;
                    this.selectedWallet = {};
                    this.userTransactions = [];
                    this.loadWallet();
                    this.loadConnections();
                    this.loadDefaultConnections();
                    this.loadBlockChain();
                    this.loadNumberOfBlocks();
                    this.loadProgress();
                    this.isValidAddress = false;
                    this.blockViewMode = 'recentBlocks';
                    //Set interval function for load wallet every 15 seconds
                    setInterval(function () {
                        _this.loadWallet();
                    }, 30000);
                    setInterval(function () {
                        _this.loadConnections();
                        _this.loadBlockChain();
                        _this.loadNumberOfBlocks();
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
                };
                //Ready button function for disable "textbox" and enable "Send" button for ready to send coin
                LoadWalletComponent.prototype.ready = function (spendId, spendaddress, spendamount) {
                    if (!spendId) {
                        toastr.error("Please select from id");
                        return false;
                    }
                    if (!spendaddress) {
                        toastr.error("Please enter pay to");
                        return false;
                    }
                    if (!spendamount) {
                        toastr.error("Please enter amount");
                        return false;
                    }
                    this.readyDisable = true;
                    this.sendDisable = false;
                };
                LoadWalletComponent.prototype.loadNumberOfBlocks = function () {
                    var _this = this;
                    this.numberOfBlocks = 0;
                    this.http.get('/blockchain/metadata')
                        .map(function (res) { return res.json(); })
                        .subscribe(function (data) {
                        _this.numberOfBlocks = data.head.seq;
                    });
                };
                LoadWalletComponent.prototype.loadTransactionsForWallet = function () {
                    var _this = this;
                    var addresses = [];
                    this.userTransactions = [];
                    _.each(this.wallets, function (wallet) {
                        _.each(wallet.entries, function (entry) {
                            addresses.push(entry.address);
                        });
                    });
                    _.each(addresses, function (address) {
                        _this.http.get('/explorer/address?address=' + address, {})
                            .map(function (res) { return res.json(); })
                            .subscribe(function (transactions) {
                            _.each(transactions, function (transaction) {
                                _this.userTransactions.push({ 'type': 'confirmed', 'transactionInputs': transaction.inputs, 'transactionOutputs': transaction.outputs,
                                    'actualTransaction': transaction
                                });
                            });
                        });
                    });
                };
                //Load wallet function
                LoadWalletComponent.prototype.loadWallet = function () {
                    var _this = this;
                    this.totalSky = 0;
                    this.http.post('/wallets', '')
                        .map(function (res) { return res.json(); })
                        .subscribe(function (data) {
                        if (_this.wallets == null || _this.wallets.length == 0) {
                            _.each(data, function (o) {
                                o.showChild = false;
                            });
                            _this.wallets = data;
                            if (_this.wallets.length > 0) {
                                _this.onSelectWallet(_this.wallets[0].meta.filename);
                            }
                        }
                        else {
                            data.map(function (w) {
                                var old = _.find(_this.wallets, function (o) {
                                    return o.meta.filename === w.meta.filename;
                                });
                                if (old) {
                                    _.extend(old, w);
                                }
                                else {
                                    w.showChild = false;
                                    _this.wallets.push(w);
                                }
                            });
                        }
                        //console.log("this.wallets", this.wallets);
                        //Load Balance for each wallet
                        //var inc = 0;
                        //console.log("data", data);
                        _.map(data, function (item, idx) {
                            var filename = item.meta.filename;
                            _this.loadWalletItem(filename, idx);
                        });
                        _this.walletsWithAddress = [];
                        _.map(_this.wallets, function (o, idx) {
                            _this.walletsWithAddress.push({
                                wallet: o,
                                type: 'wallet'
                            });
                            _.map(o.entries, function (_o, idx) {
                                _this.walletsWithAddress.push({
                                    entry: _o,
                                    type: 'address',
                                    wallet: o,
                                    idx: idx == 0 ? '' : '(' + idx + ')'
                                });
                            });
                        });
                        _this.loadTransactionsForWallet();
                    }, function (err) { return console.log(err); }, function () {
                        //console.log('Wallet load done')
                    });
                };
                LoadWalletComponent.prototype.checkValidAddress = function (address) {
                    var _this = this;
                    if (address === "") {
                        this.isValidAddress = false;
                    }
                    else {
                        var headers = new http_2.Headers();
                        headers.append('Content-Type', 'application/x-www-form-urlencoded');
                        this.http.get('/balance?addrs=' + address, { headers: headers })
                            .map(function (res) { return res.json(); })
                            .subscribe(
                        //Response from API
                        function (response) {
                            _this.isValidAddress = true;
                        }, function (err) {
                            //console.log("Error on load balance: " + err)
                            _this.isValidAddress = false;
                        }, function () {
                        });
                    }
                };
                LoadWalletComponent.prototype.loadWalletItem = function (address, inc) {
                    var _this = this;
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    this.http.get('/wallet/balance?id=' + address, { headers: headers })
                        .map(function (res) { return res.json(); })
                        .subscribe(
                    //Response from API
                    function (response) {
                        //console.log('load done: ' + inc, response);
                        _this.wallets[inc].balance = response.confirmed.coins / 1000000;
                        _this.totalSky += _this.wallets[inc].balance;
                    }, function (err) { return console.log("Error on load balance: " + err); }, function () {
                        //console.log('Balance load done')
                    });
                    //get address balances
                    this.wallets[inc].entries.map(function (entry) {
                        _this.http.get('/balance?addrs=' + entry.address, { headers: headers })
                            .map(function (res) { return res.json(); })
                            .subscribe(
                        //Response from API
                        function (response) {
                            //console.log('balance:' + entry.address, response);
                            entry.balance = response.confirmed.coins / 1000000;
                        }, function (err) { return console.log("Error on load balance: " + err); }, function () {
                            //console.log('Balance load done')
                        });
                    });
                };
                LoadWalletComponent.prototype.loadConnections = function () {
                    var _this = this;
                    this.http.post('/network/connections', '')
                        .map(function (res) { return res.json(); })
                        .subscribe(function (data) {
                        //console.log("connections", data);
                        _this.connections = data.connections;
                    }, function (err) { return console.log("Error on load connection: " + err); }, function () {
                        //console.log('Connection load done')
                    });
                };
                LoadWalletComponent.prototype.loadTransactions = function () {
                    var _this = this;
                    this.historyTable = [];
                    this.http.get('/lastTxs', {})
                        .map(function (res) { return res.json(); })
                        .subscribe(function (data) {
                        console.log("transactions", data);
                        _this.historyTable = _this.historyTable.concat(data);
                        _this.setHistoryPage(1);
                    }, function (err) { return console.log("Error on load transactions: " + err); }, function () {
                        //console.log('Connection load done')
                    });
                    this.http.get('/pendingTxs', {})
                        .map(function (res) { return res.json(); })
                        .subscribe(function (data) {
                        console.log("pending transactions", data);
                        _this.historyTable = _this.historyTable.concat(data);
                        _this.setHistoryPage(1);
                    }, function (err) { return console.log("Error on pending transactions: " + err); }, function () {
                    });
                };
                LoadWalletComponent.prototype.GetTransactionAmount = function (transaction) {
                    var ret = 0;
                    _.each(transaction.outputs, function (o) {
                        ret += Number(o.coins);
                    });
                    return ret;
                };
                LoadWalletComponent.prototype.GetTransactionAmount2 = function (transaction) {
                    var ret = 0;
                    _.each(transaction.outputs, function (o) {
                        ret += Number(o.coins);
                    });
                    return ret;
                };
                LoadWalletComponent.prototype.GetBlockAmount = function (block) {
                    var ret = [];
                    _.each(block.body.txns, function (o) {
                        _.each(o.outputs, function (_o) {
                            ret.push(_o.coins);
                        });
                    });
                    return ret.join(",");
                };
                LoadWalletComponent.prototype.GetBlockTotalAmount = function (block) {
                    var ret = 0;
                    _.each(block.body.txns, function (o) {
                        _.each(o.outputs, function (_o) {
                            ret += Number(_o.coins);
                        });
                    });
                    return ret;
                };
                LoadWalletComponent.prototype.loadDefaultConnections = function () {
                    var _this = this;
                    this.http.post('/network/defaultConnections', '')
                        .map(function (res) { return res.json(); })
                        .subscribe(function (data) {
                        //console.log("default connections", data);
                        _this.defaultConnections = data;
                    }, function (err) { return console.log("Error on load default connection: " + err); }, function () {
                        //console.log('Default connections load done')
                    });
                };
                LoadWalletComponent.prototype.loadBlockChain = function () {
                    var _this = this;
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    this.http.get('/last_blocks?num=10', { headers: headers })
                        .map(function (res) { return res.json(); })
                        .subscribe(function (data) {
                        //console.log("blockchain", data);
                        _this.blockChain = _.sortBy(data.blocks, function (o) {
                            return o.header.seq * (-1);
                        });
                        _this.setBlockPage(1);
                    }, function (err) { return console.log("Error on load blockchain: " + err); }, function () {
                        //console.log('blockchain load done');
                    });
                }; //Load progress function for Skycoin
                LoadWalletComponent.prototype.loadProgress = function () {
                    var _this = this;
                    //Post method executed
                    this.http.post('/blockchain/progress', '')
                        .map(function (res) { return res.json(); })
                        .subscribe(
                    //Response from API
                    function (response) { _this.progress = (parseInt(response.current, 10) + 1) / parseInt(response.highest, 10) * 100; }, function (err) { return console.log("Error on load progress: " + err); }, function () {
                        //console.log('Progress load done:' + this.progress)
                    });
                };
                LoadWalletComponent.prototype.toggleShowChild = function (wallet) {
                    wallet.showChild = !wallet.showChild;
                };
                //Switch tab function
                LoadWalletComponent.prototype.switchTab = function (mode, wallet) {
                    //"Textbox" and "Ready" button enable in Send tab while switching tabs
                    this.sendDisable = true;
                    this.readyDisable = false;
                    this.displayMode = mode;
                    if (wallet) {
                        this.spendid = wallet.meta.filename;
                        this.selectedWallet = _.find(this.wallets, function (o) {
                            return o.meta.filename === wallet.meta.filename;
                        });
                        console.log("selected wallet", this.spendid, this.selectedWallet);
                    }
                };
                LoadWalletComponent.prototype.selectMenu = function (menu, event) {
                    this.displayMode = this.displayModeEnum.fifth;
                    event.preventDefault();
                    this.selectedMenu = menu;
                    if (menu == 'Outputs') {
                        if (this.outputComponent) {
                            this.outputComponent.refreshOutputs();
                        }
                    }
                    if (menu == 'PendingTxns') {
                        if (this.pendingTxnComponent) {
                            this.pendingTxnComponent.refreshPendingTxns();
                        }
                    }
                };
                LoadWalletComponent.prototype.getDateTimeString = function (ts) {
                    return moment.unix(ts).format("YYYY-MM-DD HH:mm");
                };
                LoadWalletComponent.prototype.getElapsedTime = function (ts) {
                    return moment().unix() - ts;
                };
                //Show QR code function for show QR popup
                LoadWalletComponent.prototype.showQR = function (address) {
                    this.QrAddress = address;
                    this.QrIsVisible = true;
                };
                //Hide QR code function for hide QR popup
                LoadWalletComponent.prototype.hideQrPopup = function () {
                    this.QrIsVisible = false;
                };
                //Show wallet function for view New wallet popup
                LoadWalletComponent.prototype.showNewWalletDialog = function () {
                    this.NewWalletIsVisible = true;
                };
                //Hide wallet function for hide New wallet popup
                LoadWalletComponent.prototype.hideWalletPopup = function () {
                    this.NewWalletIsVisible = false;
                };
                LoadWalletComponent.prototype.showNewDefaultConnectionDialog = function () {
                    this.NewDefaultConnectionIsVisible = true;
                };
                LoadWalletComponent.prototype.hideNewDefaultConnectionDialog = function () {
                    this.NewDefaultConnectionIsVisible = false;
                };
                LoadWalletComponent.prototype.showEditDefaultConnectionDialog = function (item) {
                    this.oldConnection = item;
                    this.EditDefaultConnectionIsVisible = true;
                };
                LoadWalletComponent.prototype.hideEditDefaultConnectionDialog = function () {
                    this.EditDefaultConnectionIsVisible = false;
                };
                LoadWalletComponent.prototype.createDefaultConnection = function (connectionValue) {
                    //console.log("new value", connectionValue);
                    this.defaultConnections.push(connectionValue);
                    this.NewDefaultConnectionIsVisible = false;
                };
                LoadWalletComponent.prototype.updateDefaultConnection = function (connectionValue) {
                    //console.log("old/new value", this.oldConnection, connectionValue);
                    var idx = this.defaultConnections.indexOf(this.oldConnection);
                    this.defaultConnections.splice(idx, 1);
                    this.defaultConnections.splice(idx, 0, connectionValue);
                    this.EditDefaultConnectionIsVisible = false;
                };
                LoadWalletComponent.prototype.deleteDefaultConnection = function (item) {
                    var idx = this.defaultConnections.indexOf(item);
                    this.defaultConnections.splice(idx, 1);
                };
                //Add new wallet function for generate new wallet in Skycoin
                LoadWalletComponent.prototype.createNewWallet = function (label, seed, addressCount) {
                    var _this = this;
                    if (addressCount < 1) {
                        //alert("Please input correct address count");
                        toastr.error('Please input correct address count');
                        return;
                    }
                    //check if label is duplicated
                    var old = _.find(this.wallets, function (o) {
                        return (o.meta.label == label);
                    });
                    if (old) {
                        toastr.error('This wallet label is used already.');
                        //alert("This wallet label is used already");
                        return;
                    }
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    if (this.seedComponent) {
                        seed = this.seedComponent.getCurrentSeed();
                    }
                    //Post method executed
                    var stringConvert = 'label=' + label + '&seed=' + seed;
                    this.http.post('/wallet/create', stringConvert, { headers: headers })
                        .map(function (res) { return res.json(); })
                        .subscribe(function (response) {
                        console.log(response);
                        if (addressCount > 1) {
                            var repeats = [];
                            for (var i = 0; i < addressCount - 1; i++) {
                                repeats.push(i);
                            }
                            async.map(repeats, function (idx, callback) {
                                var stringConvert = 'id=' + response.meta.filename;
                                _this.http.post('/wallet/newAddress', stringConvert, { headers: headers })
                                    .map(function (res) { return res.json(); })
                                    .subscribe(function (response) {
                                    console.log(response);
                                    callback(null, null);
                                }, function (err) {
                                    callback(err, null);
                                }, function () { });
                            }, function (err, ret) {
                                if (err) {
                                    console.log(err);
                                    return;
                                }
                                //Hide new wallet popup
                                _this.NewWalletIsVisible = false;
                                toastr.info("New wallet created successfully");
                                //Load wallet for refresh list
                                _this.loadWallet();
                            });
                        }
                        else {
                            //Hide new wallet popup
                            _this.NewWalletIsVisible = false;
                            toastr.info("New wallet created successfully");
                            //Load wallet for refresh list
                            _this.loadWallet();
                        }
                    }, function (err) {
                        console.log(err);
                    }, function () { });
                };
                LoadWalletComponent.prototype.addNewAddress = function (wallet) {
                    var _this = this;
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    //Post method executed
                    var stringConvert = 'id=' + wallet.meta.filename;
                    this.http.post('/wallet/newAddress', stringConvert, { headers: headers })
                        .map(function (res) { return res.json(); })
                        .subscribe(function (response) {
                        console.log(response);
                        toastr.info("New address created successfully");
                        //Load wallet for refresh list
                        _this.loadWallet();
                    }, function (err) {
                        console.log(err);
                    }, function () { });
                };
                //Load wallet seed function
                LoadWalletComponent.prototype.openLoadWallet = function (walletName, seed) {
                    this.loadSeedIsVisible = true;
                };
                //Hide load wallet seed function
                LoadWalletComponent.prototype.hideLoadSeedWalletPopup = function () {
                    this.loadSeedIsVisible = false;
                };
                //Load wallet seed function for create new wallet with name and seed
                LoadWalletComponent.prototype.createWalletSeed = function (walletName, seed) {
                    var _this = this;
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    var stringConvert = 'name=' + walletName + '&seed=' + seed;
                    //Post method executed
                    this.http.post('/wallet/create', stringConvert, { headers: headers })
                        .map(function (res) { return res.json(); })
                        .subscribe(function (response) {
                        //Hide load wallet seed popup
                        _this.loadSeedIsVisible = false;
                        //Load wallet for refresh list
                        _this.loadWallet();
                    }, function (err) { return console.log("Error on create load wallet seed: " + JSON.stringify(err)); }, function () {
                        //console.log('Load wallet seed done')
                    });
                };
                LoadWalletComponent.prototype.sortHistory = function (key) {
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
                    if (key == 'time') {
                        this.historyTable = _.sortBy(this.historyTable, function (o) {
                            return o.txn.timestamp;
                        });
                    }
                    else if (key == 'amount') {
                        this.historyTable = _.sortBy(this.historyTable, function (o) {
                            return Number(o[key]);
                        });
                    }
                    else if (key == 'address') {
                        this.historyTable = _.sortBy(this.historyTable, function (o) {
                            return o[key];
                        });
                    }
                    ;
                    if (this.sortDir[key] == -1) {
                        this.historyTable = this.historyTable.reverse();
                    }
                    this.setHistoryPage(this.historyPager.currentPage);
                };
                LoadWalletComponent.prototype.filterHistory = function (address) {
                    console.log("filterHistory", address);
                    this.filterAddressVal = address;
                };
                LoadWalletComponent.prototype.updateStatusOfTransaction = function (txid, metaData) {
                    var _this = this;
                    var self = this;
                    var transactionConfirmed = false;
                    Rx_1.Observable.timer(0, 1000).map(function (i) {
                        if (transactionConfirmed) {
                            throw new Error("Transaction confirmed");
                        }
                        var headers = new http_2.Headers();
                        headers.append('Content-Type', 'application/x-www-form-urlencoded');
                        _this.http.get('/transaction?txid=' + txid, { headers: headers })
                            .map(function (res) { return res.json(); })
                            .subscribe(function (res) {
                            transactionConfirmed = res.status.confirmed;
                            self.pendingTable = [];
                            self.pendingTable.push({ 'time': res.txn.timestamp, 'status': res.status.confirmed ? 'Completed' : 'Unconfirmed', 'amount': metaData.amount, 'txId': txid, 'address': metaData.address });
                            self.loadWallet();
                        }, function (err) {
                            console.log("Error on load transaction: " + err);
                        }, function () {
                        });
                    }).subscribe(function () {
                    }, function (err) {
                        console.log("Transaction confirmed");
                    });
                };
                LoadWalletComponent.prototype.spend = function (spendid, spendaddress, spendamount) {
                    var _this = this;
                    var amount = Number(spendamount);
                    if (amount < 1) {
                        toastr.error('Cannot send values less than 1.');
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
                    var self = this;
                    this.http.post('/wallet/spend', stringConvert, { headers: headers })
                        .map(function (res) { return res.json(); })
                        .subscribe(function (response) {
                        console.log(response);
                        _this.updateStatusOfTransaction(response.txn.txid, { address: spendaddress, amount: amount });
                        _this.readyDisable = false;
                        _this.sendDisable = true;
                        self.spendAddress.nativeElement.value = '';
                        self.spendAmount.nativeElement.value = 0;
                        self.transactionNote.nativeElement.value = '';
                        self.isValidAddress = false;
                    }, function (err) {
                        _this.readyDisable = false;
                        _this.sendDisable = true;
                        var logBody = err._body;
                        if (logBody == 'Invalid "coins" value') {
                            toastr.error('Incorrect amount value.');
                            return;
                        }
                        else if (logBody == 'Invalid connection') {
                            toastr.error(logBody);
                            return;
                        }
                        else {
                            var logContent = JSON.parse(logBody.substring(logBody.indexOf("{")));
                            toastr.error(logContent.error);
                        }
                        //this.pendingTable.push({complete: 'Pending', address: spendaddress, amount: spendamount});
                    }, function () {
                        self.spendAddress.nativeElement.value = '';
                        self.spendAmount.nativeElement.value = 0;
                        self.transactionNote.nativeElement.value = '';
                        self.isValidAddress = false;
                        $("#send_pay_to").val("");
                        $("#send_amount").val(0);
                    });
                };
                LoadWalletComponent.prototype.setHistoryPage = function (page) {
                    this.historyPager.totalPages = this.historyTable.length;
                    if (page < 1 || page > this.historyPager.totalPages) {
                        return;
                    }
                    // get pager object from service
                    this.historyPager = this.pagerService.getPager(this.historyTable.length, page);
                    console.log("this.historyPager", this.historyPager);
                    // get current page of items
                    this.historyPagedItems = this.historyTable.slice(this.historyPager.startIndex, this.historyPager.endIndex + 1);
                    //console.log('this.pagedItems', this.historyTable, this.pagedItems);
                };
                LoadWalletComponent.prototype.setBlockPage = function (page) {
                    this.blockPager.totalPages = this.blockChain.length;
                    if (page < 1 || page > this.blockPager.totalPages) {
                        return;
                    }
                    // get pager object from service
                    this.blockPager = this.pagerService.getPager(this.blockChain.length, page);
                    // get current page of items
                    this.blockPagedItems = this.blockChain.slice(this.blockPager.startIndex, this.blockPager.endIndex + 1);
                    //console.log("this.blockPagedItems", this.blockPagedItems);
                };
                LoadWalletComponent.prototype.searchHistory = function (searchKey) {
                    console.log(searchKey);
                };
                LoadWalletComponent.prototype.searchBlockHistory = function (searchKey) {
                    console.log(searchKey);
                };
                LoadWalletComponent.prototype.onSelectWallet = function (val) {
                    console.log("onSelectWallet", val);
                    //this.selectedWallet = val;
                    this.spendid = val;
                    this.selectedWallet = _.find(this.wallets, function (o) {
                        return o.meta.filename === val;
                    });
                };
                LoadWalletComponent.prototype.showBlockDetail = function (block) {
                    //change viewMode as blockDetail
                    this.blockViewMode = 'blockDetail';
                    this.selectedBlock = block;
                };
                LoadWalletComponent.prototype.showRecentBlock = function () {
                    this.blockViewMode = 'recentBlocks';
                };
                LoadWalletComponent.prototype.showBlockTransactionDetail = function (txns) {
                    this.blockViewMode = 'blockTransactionDetail';
                    this.selectedBlockTransaction = txns;
                };
                LoadWalletComponent.prototype.showTransactionDetail = function (txId) {
                    var _this = this;
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    this.http.get('/transaction?txid=' + txId, { headers: headers })
                        .map(function (res) { return res.json(); })
                        .subscribe(
                    //Response from API
                    function (response) {
                        console.log(response);
                        _this.blockViewMode = 'blockTransactionDetail';
                        _this.selectedBlockTransaction = response.txn;
                    }, function (err) {
                        console.log("Error on load transaction: " + err);
                    }, function () {
                    });
                };
                LoadWalletComponent.prototype.showBlockAddressDetail = function (address) {
                    var _this = this;
                    this.blockViewMode = 'blockAddressDetail';
                    this.selectedBlockAddress = address;
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    var txList = [];
                    async.parallel([
                        function (callback) {
                            _this.http.get('/balance?addrs=' + address, { headers: headers })
                                .map(function (res) { return res.json(); })
                                .subscribe(
                            //Response from API
                            function (response) {
                                //console.log(response);
                                _this.selectedBlockAddressBalance = response.confirmed.coins / 1000000;
                                callback(null, null);
                            }, function (err) {
                                callback(err, null);
                                //console.log("Error on load balance: " + err)
                            }, function () {
                            });
                        },
                        function (callback) {
                            _this.http.get('/address_in_uxouts?address=' + address, { headers: headers })
                                .map(function (res) { return res.json(); })
                                .subscribe(
                            //Response from API
                            function (response) {
                                console.log("address_in_uxouts", response);
                                _.map(response, function (o) {
                                    o.type = 'in';
                                    txList.push(o);
                                });
                                callback(null, null);
                            }, function (err) {
                                callback(err, null);
                                //console.log("Error on load balance: " + err)
                            }, function () {
                            });
                        },
                        function (callback) {
                            _this.http.get('/address_out_uxouts?address=' + address, { headers: headers })
                                .map(function (res) { return res.json(); })
                                .subscribe(
                            //Response from API
                            function (response) {
                                console.log("address_out_uxouts", response);
                                _.map(response, function (o) {
                                    o.type = 'out';
                                    txList.push(o);
                                });
                                callback(null, null);
                            }, function (err) {
                                callback(err, null);
                                //console.log("Error on load balance: " + err)
                            }, function () {
                            });
                        }
                    ], function (err, rets) {
                        console.log(err, rets);
                        _this.selectedBlackAddressTxList = _.sortBy(txList, function (o) {
                            return o.time;
                        });
                    });
                };
                __decorate([
                    core_1.ViewChild(outputs_component_1.SkyCoinOutputComponent), 
                    __metadata('design:type', outputs_component_1.SkyCoinOutputComponent)
                ], LoadWalletComponent.prototype, "outputComponent", void 0);
                __decorate([
                    core_1.ViewChild(pending_transactions_component_1.PendingTxnsComponent), 
                    __metadata('design:type', pending_transactions_component_1.PendingTxnsComponent)
                ], LoadWalletComponent.prototype, "pendingTxnComponent", void 0);
                __decorate([
                    core_1.ViewChild('spendaddress'), 
                    __metadata('design:type', Object)
                ], LoadWalletComponent.prototype, "spendAddress", void 0);
                __decorate([
                    core_1.ViewChild('spendamount'), 
                    __metadata('design:type', Object)
                ], LoadWalletComponent.prototype, "spendAmount", void 0);
                __decorate([
                    core_1.ViewChild('transactionNote'), 
                    __metadata('design:type', Object)
                ], LoadWalletComponent.prototype, "transactionNote", void 0);
                __decorate([
                    core_1.ViewChild(seed_component_1.SeedComponent), 
                    __metadata('design:type', seed_component_1.SeedComponent)
                ], LoadWalletComponent.prototype, "seedComponent", void 0);
                LoadWalletComponent = __decorate([
                    core_1.Component({
                        selector: 'load-wallet',
                        directives: [router_1.ROUTER_DIRECTIVES, ng2_qrcode_1.QRCodeComponent, seed_component_1.SeedComponent, skycoin_edit_component_1.SkyCoinEditComponent, outputs_component_1.SkyCoinOutputComponent, pending_transactions_component_1.PendingTxnsComponent, wallet_backup_page_component_1.WalletBackupPageComponent],
                        providers: [PagerService],
                        templateUrl: 'app/templates/wallet.html'
                    }), 
                    __metadata('design:paramtypes', [http_1.Http, PagerService])
                ], LoadWalletComponent);
                return LoadWalletComponent;
            }());
            exports_1("LoadWalletComponent", LoadWalletComponent);
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
