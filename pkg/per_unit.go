package pkg

type PerUnit struct {
	Sbase float64
}

func (p *PerUnit) Zbase(v float64) float64 {
	return v * v / p.Sbase
}

func (p *PerUnit) R(siR, v float64) float64 {
	return siR / p.Zbase(v)
}

func (p *PerUnit) X(siX, v float64) float64 {
	return siX / p.Zbase(v)
}
