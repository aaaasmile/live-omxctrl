package omx

type SPstateplaying int

func (s *SPstateplaying) String() string {
	switch *s {
	case SPoff:
		return "off"
	case SPplaying:
		return "playing"
	case SPpause:
		return "pause"
	}
	return ""
}

const (
	SPoff = iota
	SPplaying
	SPpause
)

type SMstatemute int

func (s *SMstatemute) String() string {
	switch *s {
	case SMnormal:
		return "normal"
	case SMmuted:
		return "muted"
	}
	return ""
}

const (
	SMnormal = iota
	SMmuted
)

type StateOmx struct {
	CurrURI      string
	StatePlaying SPstateplaying
	StateMute    SMstatemute
}
