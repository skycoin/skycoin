export class Wallet {
  label: string;
  filename: string;
  seed: string;
  coins: number;
  hours: number;
  addresses: Address[];
  visible?: boolean;
  hideEmpty?: boolean;
}

export class Address {
  address: string;
  coins: number;
  hours: number;
}

export class GetWalletsResponseWallet {
  meta: GetWalletsResponseMeta;
  entries: GetWalletsResponseEntry[];
}

export class GetWalletsResponseMeta {
  label: string;
  filename: string;
  seed: string;
}

export class GetWalletsResponseEntry {
  address: string;
}
