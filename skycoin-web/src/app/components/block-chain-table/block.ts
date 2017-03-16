export class Block {
  constructor(
      public header:BlockHeader,
      public body:BlockBody

  ){

  }
}

export class BlockHeader {
  constructor(
    public seq:number,
    public block_hash:string,
    public previous_block_hash:string,
    public timestamp:string,
    public fee:number,
    public version:number,
    public tx_body_hash:string,

  ){

  }
}

export class BlockBody {
  constructor(
    public txns:Transaction[]
  ){

  }
}

export class Output {
  constructor(
    public uxid:string,
    public dst:string,
    public coins:number,
    public hrs:number,
  ){

  }
}

export class BlockChainMetaDataHead{
  constructor(
    public seq:number,
    public block_hash:string,
    public previous_block_hash:string,
    public timestamp:number,
    public fee:number,
    public version:number,
    public tx_body_hash:string,
  ){

  }
}

export class BlockChainMetaData{
  constructor(
    public head:BlockChainMetaDataHead,
    public unspents:number,
    public unconfirmed:number,
  ){

  }
}


export class Transaction {
  constructor(public length:number,
              public type:number,
              public txid:string,
              public inner_hash:number,
              public sigs:string[],
              public inputs:string[],
              public outputs:Output[],
              ){

  }
}

export class CoinSupply {
  constructor(public coinSupply:number,
              public coinCap:number,
  ){

  }
}


export class BlockResponse{
  constructor(public blocks:Block[]){

  }
}
