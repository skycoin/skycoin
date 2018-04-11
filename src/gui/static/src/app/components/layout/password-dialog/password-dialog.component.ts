import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material';
import { FormControl, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-password-dialog',
  templateUrl: './password-dialog.component.html',
  styleUrls: ['./password-dialog.component.css']
})
export class PasswordDialogComponent implements OnInit {
  form: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<PasswordDialogComponent>,
    @Inject(MAT_DIALOG_DATA) private data: any,
  ) { }

  ngOnInit() {
    this.form = new FormGroup({});
    this.form.addControl('password', new FormControl('', [Validators.required]));
  }

  proceed() {
    this.dialogRef.close(this.form.get('password').value);
  }
}
