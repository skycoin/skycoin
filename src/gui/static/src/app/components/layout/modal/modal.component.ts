import { Component, Input, HostListener } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.scss'],
})
export class ModalComponent {
  @Input() useRedTitle = false;
  @Input() headline: string;
  @Input() disableDismiss: boolean;

  private dialogInternal: MatDialogRef<any>;
  @Input() set dialog(val: MatDialogRef<any>) {
    val.disableClose = true;
    this.dialogInternal = val;
  }

  @HostListener('window:keyup.esc') onKeyUp() {
    this.closePopup();
  }

  closePopup() {
    if (!this.disableDismiss) {
      this.dialogInternal.close();
    }
  }
}
