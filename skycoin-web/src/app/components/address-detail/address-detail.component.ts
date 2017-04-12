import {Component, OnInit, AfterViewInit} from '@angular/core';
import {UxOutputsService} from "./UxOutputs.service";
import {Router, ActivatedRoute, Params} from "@angular/router";
import {Observable} from "rxjs";

declare var QRCode:any;

@Component({
  selector: 'app-address-detail',
  templateUrl: './address-detail.component.html',
  styleUrls: ['./address-detail.component.css'],
})
export class AddressDetailComponent implements OnInit,  AfterViewInit{

  private UxOutputs:Observable<any>;

  private showUxID:boolean;

  private transactions:any[];//

  private currentAddress:string;

  private currentBalance:string;

  constructor(   private service:UxOutputsService,
                 private route: ActivatedRoute,
                 private router: Router) {
    this.UxOutputs=null;
    this.currentBalance="0";
    this.transactions =[];
    this.currentAddress = null;
    this.showUxID = false;
  }

  ngOnInit() {

  }

  ngAfterViewInit(){
    this.UxOutputs= this.route.params
      .switchMap((params: Params) => {
        let address = params['address'];
        this.currentAddress = address;
        let qrcode = new QRCode("qr-code");
        qrcode.makeCode(this.currentAddress);
        return this.service.getUxOutputsForAddress(address);
      });

    this.UxOutputs.subscribe((uxoutputs)=>{
      this.transactions = uxoutputs;
      console.log(uxoutputs);
    });

    this.route.params
    .switchMap((params: Params) => {
      let address = params['address'];
      return this.service.getCurrentBalanceOfAddress(address);
    }).subscribe((addressDetails)=>{
        if(addressDetails.head_outputs.length>0){
          this.currentBalance = addressDetails.head_outputs[0].coins;
        }
    });
  }

  showUxId(){
    this.showUxID = true;
    return false;
  }

  hideUxId(){
    this.showUxID = false;
    return false;
  }

}
