export function convertAsciiToHexa(str): string {
  const arr1: string[] = [];
  for (let n = 0, l = str.length; n < l; n ++) {
    const hex = Number(str.charCodeAt(n)).toString(16);
    arr1.push(hex.length !== 1 ? hex : '0' + hex);
  }
  return arr1.join('');
}
