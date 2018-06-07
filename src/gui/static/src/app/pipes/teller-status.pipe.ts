import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'tellerStatus',
})
export class TellerStatusPipe implements PipeTransform {
  private statuses = ['done', 'waiting_confirm', 'waiting_deposit', 'waiting_send'];

  transform(value: any): any {
    return this.statuses.find(status => status === value)
      ? 'teller.' + value.replace('_', '-')
      : 'teller.unknown';
  }
}
