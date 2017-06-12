export class AddressBalance{
  constructor(
    public confirmed:Coins,
    public predicted:Coins,
  ){

  }
}

export class Wallet{
  constructor(
    public entries:WalletEntry[],
    public meta:WalletMeta,
  ){

  }
}

export class WalletMeta{
  constructor(
    public coin:string,
    public filename:string,
    public label:string,
    public lastseed:string,
    public seed:string,
    public tm:number,
    public type:string,
    public version:string,

  ){

  }
}

export class WalletEntry{
  constructor(
    public address:string,
    public public_key:number,
    public secret_key:number,
  ){

  }
}


export class Coins{
  constructor(
    public coins:number,
    public hours:number,
  ){

  }
}
