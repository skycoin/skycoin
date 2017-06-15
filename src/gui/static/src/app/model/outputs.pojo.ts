/**
 * Created by napandey on 1/13/17.
 */

export interface Output {
    hash:string;
    src_tx:string;
    address:string;
    coins:string;
    seed:string;
    hours:number;
}

export interface OutputsResponse{
    head_outputs:Output[]
}
