import { Component, forwardRef, Output, EventEmitter, Input, Renderer2 } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import { Overlay } from '@angular/cdk/overlay';

import { BaseCoin } from '../../../coins/basecoin';
import { openChangeCoinModal } from '../../../utils';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';

@Component({
  selector: 'app-select-coin',
  templateUrl: 'select-coin.component.html',
  styleUrls: ['select-coin.component.scss'],
  providers: [{
    provide: NG_VALUE_ACCESSOR,
    useExisting: forwardRef(() => SelectCoinComponent),
    multi: true
  }]
})
export class SelectCoinComponent implements ControlValueAccessor {
  @Output() onCoinChanged = new EventEmitter<BaseCoin>();
  @Input() selectedCoin: BaseCoin;

  constructor(private dialog: CustomMatDialogService,
    private overlay: Overlay,
    private renderer: Renderer2) {}

  onInputClick() {
    openChangeCoinModal(this.dialog, this.renderer, this.overlay)
      .subscribe(response => {
        if (response) {
          this.selectedCoin = response;
          this.onChangeCallback(this.selectedCoin);
          this.onCoinChanged.emit(this.selectedCoin);
        }
      });
  }

  writeValue(coin: any) {
    this.selectedCoin = coin;
  }

  registerOnChange(fn: any) {
    this.onChangeCallback = fn;
  }

  registerOnTouched(fn: any) {
  }

  private onChangeCallback: any = () => {};
}
