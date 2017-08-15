import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-button',
  templateUrl: 'button.component.html',
  styleUrls: ['button.component.css']
})

export class ButtonComponent {
  @Input() form: any;
  @Input() placeholder: string;

  error: string;
  state: number;

  setLoading() {
    this.state = 0;
  }

  setSuccess() {
    this.state = 1;
    setTimeout(() => this.state = null, 3000);
  }

  setError(error: any) {
    this.error = error['_body'];
    this.state = 2;
  }

  private disabled() {
    return this.state === 0  || (!(this.form === undefined) && !(this.form && this.form.valid));
  }
}
