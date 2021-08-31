package loadavg

func Get() (*Stats, error) {
	return get()
}

// Stats represents load average values
type Stats struct {
	Load1, Load5, Load15 float64
}
