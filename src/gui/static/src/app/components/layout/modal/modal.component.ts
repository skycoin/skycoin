import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.scss']
})
export class ModalComponent implements OnChanges {
  @Input() dialog: MatDialogRef<any>;
  @Input() title: string;
  @Input() disableDismiss: boolean;

  closePopup() {
    if (!this.disableDismiss) {
      this.dialog.close();
    }
  }

  ngOnChanges(changes: SimpleChanges) {
    this.dialog.disableClose = changes.disableDismiss.currentValue;
  }
}
