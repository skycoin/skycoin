import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-onboarding-disclaimer',
  templateUrl: './onboarding-disclaimer.component.html',
  styleUrls: ['./onboarding-disclaimer.component.scss'],
})
export class OnboardingDisclaimerComponent {

  acceptTerms = false;

  constructor(
    public dialogRef: MatDialogRef<OnboardingDisclaimerComponent>,
  ) { }

  closePopup() {
    this.dialogRef.close(this.acceptTerms);
  }

  setAccept(event) {
    this.acceptTerms = event.checked;
  }

}
