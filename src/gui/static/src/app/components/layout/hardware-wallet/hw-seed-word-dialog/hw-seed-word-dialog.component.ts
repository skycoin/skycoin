import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';

import { MessageIcons } from '../hw-message/hw-message.component';

@Component({
  selector: 'app-hw-seed-word-dialog',
  templateUrl: './hw-seed-word-dialog.component.html',
  styleUrls: ['./hw-seed-word-dialog.component.scss'],
})
export class HwSeedWordDialogComponent implements OnInit {
  msgIcons = MessageIcons;
  form: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<HwSeedWordDialogComponent>,
    private formBuilder: FormBuilder,
  ) {}

  ngOnInit() {
    this.form = this.formBuilder.group({
      seed: ['', Validators.required],
    });
  }

  sendWord() {
    if (this.form.valid) {
      this.dialogRef.close((this.form.value.seed as string).trim().toLowerCase());
    }
  }
}
