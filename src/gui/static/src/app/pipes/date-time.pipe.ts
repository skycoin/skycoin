import { Pipe, PipeTransform } from '@angular/core';
import * as moment from 'moment';

/**
 * Takes an Unix date (UTC seconds since the Epoch) and converts it to a readable date-time
 * string. Even as Unix dates are UTC-based, the returned string is in local time. It expects
 * a number or a numeric string as argument.
 */
@Pipe({
  name: 'dateTime',
})
export class DateTimePipe implements PipeTransform {
  transform(value: number) {
    return moment.unix(value).format('YYYY-MM-DD HH:mm');
  }
}
