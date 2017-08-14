import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'transactionsAmount'
})
export class TransactionsAmountPipe implements PipeTransform {

  transform(value: any): any {
    return value.reduce((a, b) => a + b.outputs.reduce((c, d) => c + parseInt(d.coins), 0), 0);
  }
}
