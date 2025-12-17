import { Component, Input } from '@angular/core';

import { ExplorerService } from '../../../services/explorer/explorer.service';

/**
 * Different states of the list of inputs and outputs.
 */
enum ShowMoreStatus {
  /**
   * Only some elements are being shown.
   */
  showMore = 0,
  /**
   * Preparing to show all elements.
   */
  loading = 1,
  /**
   * All elements are being shown.
   */
  dontShowMore = 2,
}

/**
 * Generic control for showing the details of a transaction (ID, date, inputs, etc).
 */
@Component({
  selector: 'app-transaction-info',
  templateUrl: './transaction-info.component.html',
  styleUrls: ['./transaction-info.component.scss']
})
export class TransactionInfoComponent {
  /**
   * If set to true, makes the control stop responding to mouse clicks, to avoid problems with
   * slow operations.
   */
  disableClicks = false;
  /**
   * Allows to access the ShowMoreStatus enum in the HTML file.
   */
  showMoreStatus = ShowMoreStatus;
  /**
   * How many inputs the current transaction has.
   */
  totalInputs = 0;
  /**
   * Current state of the inputs list.
   */
  showMoreInputs: ShowMoreStatus = ShowMoreStatus.dontShowMore;
  /**
   * Inputs that will be shown in the UI.
   */
  inputsToShow: any[] = [];
  /**
   * How many outputs the current transaction has.
   */
  totalOutputs = 0;
  /**
   * Current state of the outputs list.
   */
  showMoreOutputs: ShowMoreStatus = ShowMoreStatus.dontShowMore;
  /**
   * Outputs that will be shown in the UI.
   */
  outputsToShow: any[] = [];

  /**
   * Indicates how many inputs/outputs can be shown initially. If the transaction has more inputs
   * or outputs, the used will have to click a link for showing the rest. This helps to mantain
   * good performance.
   */
  @Input() maxInitialElements = 10;

  /**
   * Indicates if the transaction has already been confirmed or not.
   */
  @Input() unconfirmed = false;

  /**
   * Internal variable for storing the transaction that is being shown.
   */
  private transactionInternal: any;
  @Input()
  set transaction(value: any) {
    this.transactionInternal = value;

    // Save how many inputs and outputs the transaction has.
    this.totalInputs = this.transaction.inputs.length;
    this.totalOutputs = this.transaction.outputs.length;

    // Show the inputs and outputs.
    this.showInitialInputs();
    this.showInitialOutputs();
  }
  get transaction(): any {
    return this.transactionInternal;
  }

  constructor(public explorer: ExplorerService) { }

  /**
   * Creates the list of inputs to be shown, but limiting it to the number of elements set
   * in maxInitialElements.
   */
  showInitialInputs() {
    if (this.totalInputs > this.maxInitialElements) {
      // Add to this.inputsToShow a max of maxInitialElements items.
      this.inputsToShow = [];
      for (let i = 0; i < this.maxInitialElements; i++) {
        this.inputsToShow.push(this.transaction.inputs[i]);
      }
      // Indicate that there are additional inputs to show.
      this.showMoreInputs = ShowMoreStatus.showMore;
    } else {
      // Show all inputs.
      this.showAllInputs();
    }
  }

  /**
   * Creates the list of outputs to be shown, but limiting it to the number of elements set
   * in maxInitialElements.
   */
  showInitialOutputs() {
    // Works similar to showInitialInputs();
    if (this.totalOutputs > this.maxInitialElements) {
      this.outputsToShow = [];
      for (let i = 0; i < this.maxInitialElements; i++) {
        this.outputsToShow.push(this.transaction.outputs[i]);
      }
      this.showMoreOutputs = ShowMoreStatus.showMore;
    } else {
      this.showAllOutputs();
    }
  }

  /**
   * Makes the control start showing all the inputs and prepares the UI in case the process
   * takes too much time to complete.
   */
  startShowingAllInputs() {
    // If the process takes too long to be completed, the UI may end blocked and all mouse
    // clicks would be processed after finishing, which could make users think that the
    // application behaves erratically (specially if the user starts clicking
    // indiscriminately trying to make the app respond and that causes a navigation after
    // finishing the operation). To avoid this, disableClicks is set to true to make the
    // control ignore all mouse clicks temporarily.
    this.disableClicks = true;
    // Indicate that the elements are being loaded.
    this.showMoreInputs = ShowMoreStatus.loading;
    // Load all the elements after 2 frames, to give the application time for updating the
    // UI, in case it gets blocked.
    setTimeout(() => this.showAllInputs(), 32);
  }

  /**
   * Makes the control start showing all the outputs and prepares the UI in case the process
   * takes too much time to complete.
   */
  startShowingAllOutputs() {
    // Works similar to startShowingAllInputs().
    this.disableClicks = true;
    this.showMoreOutputs = ShowMoreStatus.loading;
    setTimeout(() => this.showAllOutputs(), 32);
  }

  /**
   * Makes the control show all the inputs.
   */
  private showAllInputs() {
    // Updates the list with the elements.
    this.inputsToShow = this.transaction.inputs;
    // Updates the UI after a frame.
    setTimeout(() => {
      // Idicate that all elements are being shown.
      this.showMoreInputs = ShowMoreStatus.dontShowMore;
      // Accept mouse click.
      this.disableClicks = false;
    });
  }

  /**
   * Makes the control show all the outputs.
   */
  private showAllOutputs() {
    // Works similar to showAllInputs().
    this.outputsToShow = this.transaction.outputs;
    setTimeout(() => {
      this.showMoreOutputs = ShowMoreStatus.dontShowMore;
      this.disableClicks = false;
    });
  }
}
