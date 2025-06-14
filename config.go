package koboi

type Config struct {
	Type string
}

const (
	ROUND_ROBIN = iota
	WEIGHTED
	LEAST_CONNECTION
)
