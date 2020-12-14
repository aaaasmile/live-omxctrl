package omxstate

import "log"

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

type ActionTD int

///home/igors/projects/go/bin/stringer -type=actionTD

const (
	ActTerminate ActionTD = iota
	ActPlaying
	ActPause
	ActMute
	ActUnmute
)

type ActionDef struct {
	URI    string
	Action ActionTD
}

func ListenStateAction(actCh chan *ActionDef, chstatus chan *StateOmx) {
	log.Println("Waiting for action to change the state")
	var stateCurrent SPstateplaying
	stateCurrent = SPoff
	uriPlaying := ""
	for {
		st := <-actCh
		log.Println("New action in state: ", st.Action.String(), stateCurrent.String())
		stateNext := StateOmx{CurrURI: st.URI, StatePlayer: SPundef}
		switch stateCurrent {
		case SPoff:
			switch st.Action {
			case ActPlaying:
				stateNext.StatePlayer = SPplaying
				uriPlaying = st.URI
			}
		case SPplaying:
			switch st.Action {
			case ActPlaying:
				stateNext.StatePlayer = SPrestart
				uriPlaying = st.URI
			case ActPause:
				stateNext.StatePlayer = SPpause
			case ActTerminate:
				stateNext.StatePlayer = SPoff
				uriPlaying = ""
			}
		case SPpause:
			switch st.Action {
			case ActPlaying:
				stateNext.StatePlayer = SPplaying
			}
		case SPrestart:
			switch st.Action {
			case ActTerminate:
				stateNext.StatePlayer = SPplaying
			}
		}

		log.Println("Calculated next state ", stateNext.StatePlayer.String())
		if stateNext.StatePlayer != SPundef {
			log.Println("State trigger a change")
			stateCurrent = stateNext.StatePlayer
			stateNext.CurrURI = uriPlaying
			//op.mutex.Lock()
			//op.setState(&stateNext)
			//op.mutex.Unlock()
			chstatus <- &stateNext
		} else {
			log.Println("Ignored action ", st.Action.String())
		}
	}
}
