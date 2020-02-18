import { Component, Input, HostListener } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

/**
 * Parent component for the content of all modal windows. It provides the title, scroll
 * bars and more. It is designed to be added in the HTML of modal window components,
 * wrapping the content.
 */
@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.scss'],
})
export class ModalComponent {
  @Input() useRedTitle = false;
  @Input() headline: string;
  // This disables all the ways provided by default by the UI for closing the modal window.
  @Input() disableDismiss: boolean;

  // MatDialogRef of the modal window component which is using this component for wrapping
  // the contents.
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
