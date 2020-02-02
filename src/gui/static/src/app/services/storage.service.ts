import { Injectable } from '@angular/core';
import { ApiService } from './api.service';

export enum StorageType {
  CLIENT = 'client',
  NOTES = 'txid',
}

@Injectable()
export class StorageService {

  constructor(
    private apiService: ApiService,
  ) { }

  get(type: StorageType, key: string) {
    const params = <any> { type };

    if (key) {
      params.key = key;
    }

    return this.apiService.get('data', params, { useV2: true});
  }

  store(type: StorageType, key: string, value: string) {
    return this.apiService.post('data', {
      type: type,
      key: key,
      val: value,
    }, {useV2: true});
  }
}
