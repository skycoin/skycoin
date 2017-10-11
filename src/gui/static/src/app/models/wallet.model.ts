interface WalletModelMeta {
  filename: string;
  label: string;
  seed: string;
}

export interface WalletModel {
  meta: WalletModelMeta;
  entries: any[];
  balance?: number;
  hours?: number;
  visible?: boolean;
}
