import { Component, OnDestroy, Input, Output, EventEmitter } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';

import { SeedWordDialogComponent, WordAskedReasons } from '../../../../../layout/seed-word-dialog/seed-word-dialog.component';
import { MsgBarService } from '../../../../../../services/msg-bar.service';
import { ApiService } from '../../../../../../services/api.service';

/**
 * Component for entering a seed with the assisted mode, which uses a modal window for entering
 * each word independently.
 */
@Component({
  selector: 'app-assisted-seed-field',
  templateUrl: './assisted-seed-field.component.html',
  styleUrls: ['./assisted-seed-field.component.scss'],
})
export class AssistedSeedFieldComponent implements OnDestroy {
  // If the control is just for confirming a predefined seed (true) or not.
  @Input() showConfirmationOnly: boolean;
  // Predefined seed to confirm if showConfirmationOnly is true.
  @Input() seedToConfirm: string;
  // If the component is being shown on the wizard (true) or not (false).
  @Input() onboarding: boolean;
  // Allows to deactivate the component while the system is busy.
  @Input() busy = false;
  // Text to be shown if the user has not entered the seed.
  @Input() emptyText = '';
  // Text to be shown if showConfirmationOnly is true and the user already confirmed the seed.
  @Input() confirmedText = '';
  // How many words the seed must have.
  @Input() howManyWords: number;
  // Reason for asking the seed.
  @Input() reasonForAsking: WordAskedReasons;
  // Event emited when the seed entered by the user changes.
  @Output() seedChanged = new EventEmitter<void>();

  /**
   * Last seed entered by the user.
   */
  lastAssistedSeed = null;

  // Indicates if the component is checking if the seed is valid before accepting it.
  checkingSeed = false;

  // Saves the words the user enters while using the assisted mode.
  private partialSeed: string[];

  constructor(
    private apiService: ApiService,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
  ) { }

  ngOnDestroy() {
    this.seedChanged.complete();
  }

  // Starts the assisted procedure for entering the seed.
  enterSeed() {
    // Do not continue if trying to confirm an already confirmed seed.
    if (!this.showConfirmationOnly || !this.lastAssistedSeed) {
      this.askForWord(0);
      this.msgBarService.hide();
    }
  }

  /**
   * Recursively asks the user to enter the words of the seed.
   * @param wordIndex Index of the word which is going to be requested on this step. Must be
   * 0 when starting to ask for the words.
   */
  private askForWord(wordIndex: number) {
    if (wordIndex === 0) {
      this.partialSeed = [];
    }

    // Open the modal window for entering the seed word.
    return SeedWordDialogComponent.openDialog(this.dialog, {
      reason: this.reasonForAsking,
      wordNumber: wordIndex + 1,
    }).afterClosed().subscribe(word => {
      if (word) {
        // If confirming a seed, check if the user entered the requested word.
        if (this.showConfirmationOnly) {
          const seedToConfirmWords = this.seedToConfirm.split(' ');
          if (word !== seedToConfirmWords[wordIndex]) {
            this.msgBarService.showError('wallet.new.seed.incorrect-word-error');

            return;
          }
        }

        // Add the entered word to the list of words the user already entered.
        this.partialSeed[wordIndex] = word;
        wordIndex += 1;

        if (wordIndex < this.howManyWords) {
          // Ask for the next word.
          this.askForWord(wordIndex);
        } else {
          // Build the seed.
          const enteredSeed = this.partialSeed.join(' ');

          if (this.showConfirmationOnly) {
            this.lastAssistedSeed = enteredSeed;
            this.seedChanged.emit();
          } else {
            // Check the seed and use it only if it is valid.
            this.checkingSeed = true;

            this.apiService.post('wallet/seed/verify', {seed: enteredSeed}, {useV2: true}).subscribe(
              () => {
                this.lastAssistedSeed = enteredSeed;
                this.seedChanged.emit();
                this.checkingSeed = false;
              },
              () => {
                this.msgBarService.showError('wallet.new.seed.invalid-seed-error');
              this.checkingSeed = false;
              },
            );
          }
        }
      }
    });
  }
}
