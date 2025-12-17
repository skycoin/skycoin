import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'tellerStatus'
})
export class TellerStatusPipe implements PipeTransform {

  transform(value: any): any {
    switch (value) {
      case 'waiting_deposit':
        return 'teller.waiting-deposit';
      case 'waiting_send':
        return 'teller.waiting-send';
      case 'waiting_confirm':
        return 'teller.waiting-confirm';
      case 'done':
        return 'teller.done';
      default:
        return 'teller.unknown';
    }
  }
}
