import { ErrorHandler, Injectable} from '@angular/core';

@Injectable()
export class AppErrorHandler implements ErrorHandler {
  handleError(error: Error) { }
}