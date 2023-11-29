package troner

func Trx2Sun(amt float64) uint64 { return uint64(amt * 1000000) }
func Sun2Trx(amt uint64) float64 { return float64(amt) / 1000000 }

func TxSize2Fee(size int) float64 { return float64(size) * 0.001 }
