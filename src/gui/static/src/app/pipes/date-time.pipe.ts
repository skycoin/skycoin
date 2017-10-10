import { Pipe, PipeTransform } from '@angular/core';
import * as moment from 'moment';

@Pipe({
  name: 'dateTime'
})
export class DateTimePipe implements PipeTransform {

  transform(value: any) {
    return moment.unix(value).format('YYYY-MM-DD HH:mm');
  }
}
