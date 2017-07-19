export interface TransactionInfo {
  actualTransaction?: TransitionContent,
  confirmed?: string,
  transactionInputs?: Array<TransactionInput>,
  transactionOutputs?: Array<TransactionOutput>,
  type?: string
}

export interface TransitionContent {
  inner_hash?: string;
  inputs?: Array<TransactionInput>;
  outputs?: Array<TransactionOutput>;
  length?: number;
  sigs?: Array<string>;
  status?: TransactionStatus;
  timestamp?: number;
  txid?: string;
  type?: number;
}
export interface TransactionStatus {
  block_seq?: string;
  confirmed?: boolean;
  height?: number;
  unconfirmed?: boolean;
  unknown?: boolean;
}
export interface TransactionInput {
  owner?: string;
  uxid?: string;
}

export interface TransactionOutput {
  coins?: string;
  dst?: string;
  hours?: number;
  uxid?: string;
}