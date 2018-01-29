import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'tellerStatus'
})
export class TellerStatusPipe implements PipeTransform {

  transform(value: any): any {
    switch (value) {
      case 'done':
        return 'Completed';
      case 'waiting_confirm':
        return 'Waiting for confirmation';
      case 'waiting_deposit':
        return 'Waiting for Bitcoin deposit';
      case 'waiting_send':
        return 'Waiting to send Skycoin';
      default:
        return 'Unknown';
    }
  }
}
