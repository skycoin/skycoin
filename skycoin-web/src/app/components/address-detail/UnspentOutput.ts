export class UnspentOutput {
  constructor(
    public uxid:string,
    public time:number,
    public src_block_seq:number,
    public src_tx:string,
    public owner_address:string,
    public coins:number,
    public hours:number,
    public spent_block_seq:number,
    public spent_tx:number
  ){

  }
}

export class AddressBalanceResponse {
  constructor(
    public head_outputs:HeadOutput[],
    public outgoing_outputs:any[],
    public incoming_outputs:any[],
  ){

  }
}

export class HeadOutput {
  constructor(
    public hash:string,
    public src_tx:string,
    public address:string,
    public coins:string,
    public hours:number,
  ){
  }
}

