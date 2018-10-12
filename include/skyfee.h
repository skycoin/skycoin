#ifndef CALLFEECALCULATOR
#define CALLFEECALCULATOR
static inline GoUint32_ callFeeCalculator(FeeCalculator* feeCalc, Transaction__Handle handle, GoUint64_* pFee){
  return feeCalc->callback(handle, pFee, feeCalc->context);
}
#endif
