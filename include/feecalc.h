#ifndef CALLFEECALCULATOR
#define CALLFEECALCULATOR
static inline GoUint32_ callFeeCalculator(FeeCalc feeCalc, Transaction__Handle handle, GoUint64_* pFee){
  return feeCalc(handle, pFee);
}
#endif
