export function testCases(this: any, values: any[], func: (value: any) => void) {
  for (let i = 0, count = values.length; i < count; i++) {
    func.apply(this, [values[i]]);
  }
}
