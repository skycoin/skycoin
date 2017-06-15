export class PendingTxn{
  constructor(
      public transaction: Transaction,
      public received: string,
      public checked: string,
      public announced: string
  ){}
}
export class Transaction {
  constructor(public length:number,
              public type:number,
              public txid:string,
              public inner_hash:number,
              public sigs:string[],
              public inputs:string[],
              public outputs:Output[]
  ){

  }
}
export class Output {
  constructor(
      public uxid:string,
      public dst:string,
      public coins:number,
      public hrs:number
  ){

  }
}

