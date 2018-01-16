import { Component, Input } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.scss']
})
export class ModalComponent {
  @Input() dialog: MatDialogRef<any>;
  @Input() title: string;

  closePopup() {
    this.dialog.close();
  }
}
