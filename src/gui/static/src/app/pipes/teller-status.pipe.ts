import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'tellerStatus'
})
export class TellerStatusPipe implements PipeTransform {

  transform(value: any): any {
    switch (value) {
      case 'waiting_deposit':
        return 'Waiting for deposit';
      case 'waiting_send':
        return 'Waiting for sending out';
      case 'waiting_confirm':
        return 'Waiting for confirmation';
      case 'done':
        return 'Completed';
      default:
        return 'Unknown';
    }
  }
}
