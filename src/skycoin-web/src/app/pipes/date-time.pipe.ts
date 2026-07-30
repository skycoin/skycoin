import { Pipe, PipeTransform } from '@angular/core';
import moment from 'moment';

@Pipe({
    name: 'dateTime',
    standalone: false
})
export class DateTimePipe implements PipeTransform {

  transform(value: any) {
    return moment.unix(value).format('YYYY-MM-DD HH:mm');
  }
}
