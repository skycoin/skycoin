import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'sky'
})
export class SkyPipe implements PipeTransform {
  transform(value: number) {
    if (value == null || value < 0) {
      return 'loading .. ';
    } else {
      return (value ? (value / 1000000) : 0) + ' SKY';
    }
  }
}
