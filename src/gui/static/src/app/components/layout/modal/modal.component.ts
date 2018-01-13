import { Component, Input } from '@angular/core';
import { MdDialogRef } from '@angular/material';

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.scss']
})
export class ModalComponent {
  @Input() dialog: MdDialogRef<any>;
  @Input() title: string;

  closePopup() {
    this.dialog.close();
  }
}
