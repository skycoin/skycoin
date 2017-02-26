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
