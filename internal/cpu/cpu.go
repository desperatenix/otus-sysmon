package cpu

type Stats struct {
	Total, User, System, Idle float64
}

func Get() (*Stats, error) {
	return get()
}
