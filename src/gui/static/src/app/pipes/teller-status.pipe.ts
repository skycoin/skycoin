import { Pipe, PipeTransform } from '@angular/core';

/**
 * Returns the variable, on the translation file, to display the name of one of the states
 * returned by the teller service, using the translate pipe. It expects the value to be a
 * string with one of the known states, or a variable for indicating that the state is unknown
 * is returned.
 */
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
