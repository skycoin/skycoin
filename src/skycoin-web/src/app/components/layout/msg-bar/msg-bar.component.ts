import { ChangeDetectionStrategy, ChangeDetectorRef, Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';

export enum MsgBarIcons {
  Error = 'error',
  Done = 'done',
  Warning = 'warning',
}

export enum MsgBarColors {
  Red = 'red-background',
  Green = 'green-background',
  Yellow = 'yellow-background',
}

export class MsgBarConfig {
  title?: string;
  text!: string;
  icon?: MsgBarIcons;
  color?: MsgBarColors;
}

@Component({
    selector: 'app-msg-bar',
    templateUrl: './msg-bar.component.html',
    styleUrls: ['./msg-bar.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class MsgBarComponent implements OnInit {
  config = new MsgBarConfig();
  visible = false;
  UiIsShowingModalWindow = false;

  // Dialogs open in the top layer, above any z-index, so the bar joins it
  // each time it appears to be seen over an open dialog.
  @ViewChild('bar') set bar(el: ElementRef<HTMLElement> | undefined) {
    const e = el?.nativeElement as any;
    if (e && typeof e.showPopover === 'function' && !e.matches(':popover-open')) {
      e.showPopover();
    }
  }

  constructor(private dialog: CustomMatDialogService,
  private changeDetectorRef: ChangeDetectorRef,) { }

  ngOnInit() {
    this.dialog.showingDialog.subscribe(value => {
      this.UiIsShowingModalWindow = value;
      this.changeDetectorRef.markForCheck();
    });
  }

  show() {
    if (this.visible) {
      this.visible = false;
      setTimeout(() => {
        this.visible = true;
        this.changeDetectorRef.markForCheck();
      }, 32);
    } else {
      this.visible = true;
      this.changeDetectorRef.markForCheck();
    }
  }

  hide() {
    this.visible = false;
    this.changeDetectorRef.markForCheck();
  }
}
