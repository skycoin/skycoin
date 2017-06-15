import { SeedService } from '../services/seed.service';
import { Seed } from '../model/seed.pojo';
import {OnInit, Component} from "@angular/core";


@Component({
  selector: 'seed-mnemonic',
  template: `
                 <textarea rows="4"  placeholder="Wallet Seed" cols="46" class="form-control" [(ngModel)]="seedValue"></textarea>
              `
  ,
  providers:[SeedService]
})

export class SeedComponent implements OnInit {

  constructor(private _seedService:SeedService){}

  seedValue: string = ''

  currentSeed:Seed = new Seed('');

  ngOnInit(): any {
    this._seedService.getMnemonicSeed().subscribe(seedReceived=>
        {
          this.currentSeed=seedReceived;
          this.seedValue = seedReceived.seed;
        },
        err => {
          console.log(err);
        }
    );
  }

  getCurrentSeed():string {
    return this.seedValue;
  }
}