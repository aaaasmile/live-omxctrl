package omx

type SPstateplaying int

func (s *SPstateplaying) String() string {
	switch *s {
	case SPundef:
		return "undef"
	case SPoff:
		return "off"
	case SPplaying:
		return "playing"
	case SPpause:
		return "pause"
	case SPrestart:
		return "restart"
	}
	return ""
}

const (
	SPundef = iota
	SPoff
	SPplaying
	SPpause
	SPrestart
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
	CurrURI     string
	StatePlayer SPstateplaying
	StateMute   SMstatemute
	Info        string
	ItemType    string
	NextItem    string
	PrevItem    string
}
