import { SeedService } from '../services/seed.service';
import { Seed } from '../model/seed.pojo';
import {OnInit, Component} from "@angular/core";


@Component({
    selector: 'seed-mnemonic',
    template: `
                 <textarea rows="4"  placeholder="Wallet Seed" cols="46" class="form-control" value="{{currentSeed.seed}}"></textarea>
              `
    ,
    providers:[SeedService]
})

export class SeedComponent implements OnInit {

    constructor(private _seedService:SeedService){}

    currentSeed:Seed = new Seed('');

    ngOnInit(): any {
       this._seedService.getMnemonicSeed().subscribe(seedReceived=>
       {
           this.currentSeed=seedReceived;
       },
        err => {
            console.log(err);
           }
       );
    }
}