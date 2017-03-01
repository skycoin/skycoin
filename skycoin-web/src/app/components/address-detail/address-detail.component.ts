import { Component, OnInit } from '@angular/core';
import {UxOutputsService} from "./UxOutputs.service";
import {Router, ActivatedRoute, Params} from "@angular/router";
import {Observable} from "rxjs";
import {QRCodeComponent} from "ng2-qrcode";

@Component({
  selector: 'app-address-detail',
  templateUrl: './address-detail.component.html',
  styleUrls: ['./address-detail.component.css'],
})
export class AddressDetailComponent implements OnInit {

  private UxOutputs:Observable<any>;

  private transactions:any[];

  private currentAddress:string;

  constructor(   private service:UxOutputsService,
                 private route: ActivatedRoute,
                 private router: Router) {
    this.UxOutputs=null;
    this.transactions =[];
    this.currentAddress = null;
  }

  ngOnInit() {
    this.UxOutputs= this.route.params
      .switchMap((params: Params) => {
        let address = params['address'];
        this.currentAddress = address;
        return this.service.getUxOutputsForAddress(address);
      });

    this.UxOutputs.subscribe((uxoutputs)=>{
      this.transactions = uxoutputs;
      console.log(uxoutputs);
    })
  }

  getCurrentBalance():string{
    let outputs = this.transactions[this.transactions.length-1].outputs;
    if(this.currentAddress){
      for(var i=0;i<outputs.length;i++){
        let currentAddress = outputs[i].dst;
        if(currentAddress == this.currentAddress){
          return outputs[i].coins;
        }
      }
    }

    return "0";
  }

}
