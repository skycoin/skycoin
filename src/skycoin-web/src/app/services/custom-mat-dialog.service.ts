import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';
import { TemplateRef, Injectable } from '@angular/core';
import { ComponentType } from '@angular/cdk/overlay';
import { Observable, BehaviorSubject, map } from 'rxjs';
import { withWalletScope } from '../wallet-scope';

@Injectable()
export class CustomMatDialogService extends MatDialog {

  get showingDialog(): Observable<boolean> {
    return this.dialogsDisplayed.asObservable().pipe(map(value => value !== 0));
  }

  private dialogsDisplayed: BehaviorSubject<number> = new BehaviorSubject<number>(0);

  open<T, D>(componentOrTemplateRef: ComponentType<T> | TemplateRef<T>, config?: MatDialogConfig<D>, ignoreAppStyle?: boolean): MatDialogRef<T, any> {
    if (!ignoreAppStyle) {
      if (!config) {
        config = new MatDialogConfig();
      }
      config.panelClass = withWalletScope('default-dialog-style');
    }
    // A MatDialogConfig carries an empty backdropClass, which overrides the
    // module's default and would leave the backdrop out of the wallet's scope.
    if (config && !config.backdropClass) {
      config.backdropClass = withWalletScope('cdk-overlay-dark-backdrop');
    }

    this.dialogsDisplayed.next(this.dialogsDisplayed.value + 1);

    const result = super.open(componentOrTemplateRef, config);

    result.afterClosed().subscribe(() => this.dialogsDisplayed.next(this.dialogsDisplayed.value - 1));

    return result;
  }
}
