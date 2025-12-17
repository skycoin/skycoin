export function convertAsciiToHexa(str): string {
  const arr1: string[] = [];
  for (let n = 0, l = str.length; n < l; n ++) {
    const hex = Number(str.charCodeAt(n)).toString(16);
    arr1.push(hex.length != 1 ? hex : '0' + hex);
  }
  return arr1.join('');
}

export function testCases(values: any[], func: (value: any) => void) {
  for (let i = 0, count = values.length; i < count; i++) {
    func.apply(this, [values[i]]);
  }
}

export interface Address {
  address: string;
  next_seed?: string;
  secret_key?: string;
  public_key?: string;
  balance?: number;
  hours?: number;
  outputs?: GetOutputsRequestOutput[];
}

export interface GetOutputsRequestOutput {
  hash: string;
  src_tx: string;
  address: string;
  coins: string;
  calculated_hours: number;
}