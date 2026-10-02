import { Pipe, PipeTransform } from '@angular/core';

/** A unix time in seconds as local 'YYYY-MM-DD HH:mm'. */
@Pipe({
    name: 'dateTime',
    standalone: false
})
export class DateTimePipe implements PipeTransform {

  transform(value: any) {
    const d = new Date(Number(value) * 1000);
    if (isNaN(d.getTime())) {
      return 'Invalid date';
    }
    const two = (n: number) => String(n).padStart(2, '0');
    return `${d.getFullYear()}-${two(d.getMonth() + 1)}-${two(d.getDate())} ${two(d.getHours())}:${two(d.getMinutes())}`;
  }
}
