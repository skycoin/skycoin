import { Pipe, PipeTransform } from '@angular/core';
import * as moment from 'moment';

@Pipe({
  name: 'dateFromNow'
})
export class DateFromNowPipe implements PipeTransform {

  transform(value: any) {
    return moment.unix(value).fromNow();
  }

}
