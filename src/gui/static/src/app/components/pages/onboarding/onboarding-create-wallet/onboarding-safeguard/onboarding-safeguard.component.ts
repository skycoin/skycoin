import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-onboarding-safeguard',
  templateUrl: './onboarding-safeguard.component.html',
  styleUrls: ['./onboarding-safeguard.component.scss'],
})
export class OnboardingSafeguardComponent {
  public acceptSafe = false;

  constructor(
    public dialogRef: MatDialogRef<OnboardingSafeguardComponent>,
  ) { }

  closePopup() {
    this.dialogRef.close(this.acceptSafe);
  }

  setAccept(event) {
    this.acceptSafe = event.checked;
  }
}
