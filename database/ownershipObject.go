package database

type Ownership struct {
	StockId    int
	Name       string
	Amount     int
	CostPerOne float64
}

func (osh Ownership) SumPrices() float64 {
	return osh.CostPerOne * float64(osh.Amount)
}
