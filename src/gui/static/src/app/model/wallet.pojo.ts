/**
 * Created by nakul.pandey@gmail.com on 01/01/17.
 */
export interface WalletMeta {
    coin:string;
    filename:string;
    label:string;
    lastseed:string;
    seed:string;
    tm:number;
    type:string;
    version:number;
}
export interface WalletAddress{
    address:string;
    public_key:string;
    secret_key:string;
}

export interface Wallet{
    meta:WalletMeta;
    entries:WalletAddress[]
}