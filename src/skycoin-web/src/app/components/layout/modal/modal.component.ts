import { Component, Input } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.scss'],
})
export class ModalComponent {
  @Input() headline: string;
  @Input() loadingProgress: number;

  private dialogInternal: MatDialogRef<any>;
  private disableDismissInternal = false;

  @Input() set dialog(value: MatDialogRef<any>) {
    this.dialogInternal = value;
    this.dialogInternal.disableClose = this.disableDismissInternal;
  }

  @Input() set disableDismiss(value: boolean) {
    this.disableDismissInternal = value;
    if (this.dialogInternal) {
      this.dialogInternal.disableClose = value;
    }
  }

  get disableDismiss(): boolean {
    return this.disableDismissInternal;
  }

  closePopup() {
    if (!this.disableDismissInternal) {
      this.dialogInternal.close();
    }
  }
}
