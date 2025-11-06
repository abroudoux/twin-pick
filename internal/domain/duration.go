package domain

type Duration int

const (
	Short Duration = iota
	Medium
	Long
)

func (d Duration) String() string {
	switch d {
	case Short:
		return "short"
	case Medium:
		return "medium"
	case Long:
		return "long"
	default:
		return "unknown"
	}
}
