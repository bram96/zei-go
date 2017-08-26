package zei

//Hook interface for ZEI changes
type Hook interface {
	OnPositionChanged(Position)
}
