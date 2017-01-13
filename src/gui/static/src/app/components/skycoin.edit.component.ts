/**
 * Created by nakul.pandey@gmail.com on 01/01/17.
 */

import {Component, Input, EventEmitter, ElementRef, Output, OnInit} from '@angular/core';
import {Wallet} from '../model/wallet.pojo';
import {WalletService} from '../services/wallet.service';
declare var toastr: any;
declare var _: any;

@Component({
    selector: 'skycoin-edit',
    styles: [
        ` #skycoin-edit-ic {
        margin-left: 10px;
        color: #d9d9d9;
        }
        .skycoin-edit-comp {
            padding:6px;
            border-radius: 3px;
        }
        .active-skycoin-edit {
            background-color: #f0f0f0;
            border: 1px solid #d9d9d9;
        }
        input {
            border-radius: 5px;
            box-shadow: none;
            border: 1px solid #dedede;
            min-width: 5px;
        }
        .skycoin-edit-buttons {
            background-color: #f0f0f0;
            border: 1px solid #ccc;
            border-top: none;
            border-radius: 0 0 3px 3px;
            box-shadow: 0 3px 6px rgba(111,111,111,0.2);
            outline: none;
            padding: 3px;
            position: absolute;
            margin-left: 6px;
            z-index: 1;
        }
        .skycoin-edit-comp:hover {
            border: 1px solid grey;
        }
        .skycoin-edit-comp:hover > skycoin-edit-ic {
            display:block;
        }
        .skycoin-edit-save {
            margin-right:3px;
        }
        .skycoin-edit-active {
            background-color: #f0f0f0;
            border: 1px solid #d9d9d9;
        }
        .ng-invalid {
                background: #ffb8b8;
            }
        .err-bubble {
            position: absolute;
            margin: 16px 100px;
            border: 1px solid red;
            font-size: 14px;
            background: #ffb8b8;
            padding: 10px;
            border-radius: 7px;
        }
       `
    ],
    template: `<span class='skycoin-edit-comp' [ngClass]="{'skycoin-edit-active':show}">
<input *ngIf='show' [ngClass]="{'ng-invalid': invalid}" (ngModelChange)="validate($event)" type='text' [(ngModel)]='text' />
<div class='err-bubble' *ngIf="invalid">{{error || " must contain " + min + " to -" + max +" chars."}}</div>
<i class="fa fa-edit" (click)='makeEditable()' id='skycoin-edit-ic' *ngIf='!show'></i>
<span *ngIf='!show' (click)='makeEditable()'>{{text || '-Empty Field-'}}</span>
</span>
<div class='skycoin-edit-buttons' *ngIf='show'>
<button class='btn-x-sm' (click)='callSave()'><i class="fa fa-check"></i></button>
<button class='btn-x-sm' (click)='cancelEditable()'><i class="fa fa-times"></i></button>
</div>`,

    host: {
        "(document: click)": "compareEvent($event)",
        "(click)": "trackEvent($event)"
    },
    providers: [WalletService]
})

export class SkyCoinEditComponent implements OnInit {

    @Input('text') text;
    @Input('wallets') wallets: Wallet[];
    @Input('walletId') walletId;

    @Output() onWalletChanged = new EventEmitter();

    originalText;
    tracker;
    el: ElementRef;
    show = false;
    m: Number = 3;
    min = 0;
    max = 100;
    error;
    regex;
    invalid = false;

    constructor(el: ElementRef, private _walletService: WalletService) {
        this.el = el;
    }

    ngOnInit() {
        this.originalText = this.text;
    }

    validate(text) {
        if (this.regex) {
            var re = new RegExp('' + this.regex, "ig");
            if (re.test(text)) {
                this.invalid = false;
            }
            else {
                this.invalid = true;
            }
        }
        else {
            if ((text.length <= this.max) && (text.length >= this.min)) {
                this.invalid = false;
            }
            else {
                this.invalid = true;
            }
        }
    }

    makeEditable() {
        if (this.show == false) {
            this.show = true;
        }
    }

    compareEvent(globalEvent) {
        if (this.tracker != globalEvent && this.show) {
            this.cancelEditable();
        }
    }

    trackEvent(newHostEvent) {
        this.tracker = newHostEvent;
    }

    cancelEditable() {
        this.show = false;
        this.invalid = false;
        this.text = this.originalText;
    }

    callSave() {
        if (!this.invalid && !this.isDuplicate(this.text)) {
            var data = {};
            data["newText"] = this.text;
            data["walletId"] = this.walletId;

            this._walletService.updateWallet(data).subscribe(response => {
                    this.onWalletChanged.emit(data);
                    this.show = false;
                    toastr.info("Wallet name updated");
                },
                err => {
                    this.cancelEditable();
                    toastr.error("Unable to update the name. Please try after some time");
                }
            );
        }
    }

    isDuplicate(text) {
        var old = _.find(this.wallets, function (o) {
            return (o.meta.label == text)
        })

        if (old) {
            toastr.error("This wallet label is used already");
            this.cancelEditable();
            return true;
        }
        return false;
    }
}
