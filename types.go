package squads

type Permission uint8

const (
	Initiate Permission = 1 << 0
	Vote     Permission = 1 << 1
	Execute  Permission = 1 << 2
)

func (p Permission) Has(permission Permission) bool {
	return p&permission != 0
}
