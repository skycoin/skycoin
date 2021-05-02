import { Pipe, PipeTransform } from '@angular/core';
import * as moment from 'moment';
import { TranslateService } from '@ngx-translate/core';

/**
 * Takes an Unix date (UTC seconds since the Epoch) and returns a translated string indicating
 * how much time has passed since that moment. It expects a number or a numeric string as
 * argument.
 */
@Pipe({
  name: 'dateFromNow',
  pure: false,
})
export class DateFromNowPipe implements PipeTransform {
  constructor(
    public translateService: TranslateService,
  ) { }

  transform(value: number) {
    const diff: number = moment().unix() - value;

    if (diff < 60) {
      return this.translateService.instant('time-from-now.few-seconds');
    } else if (diff < 120) {
      return this.translateService.instant('time-from-now.minute');
    } else if (diff < 3600) {
      return this.translateService.instant('time-from-now.minutes', { time: Math.floor(diff / 60) });
    } else if (diff < 7200) {
      return this.translateService.instant('time-from-now.hour');
    } else if (diff < 86400) {
      return this.translateService.instant('time-from-now.hours', { time: Math.floor(diff / 3600) });
    } else if (diff < 172800) {
      return this.translateService.instant('time-from-now.day');
    }

    return this.translateService.instant('time-from-now.days', { time: Math.floor(diff / 86400) });
  }
}
