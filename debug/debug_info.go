package debug

type Debugs struct {
	DebugInfo *Debug // debug info
	CurNum    int    // Current() transfer times
	NextNum   int    // Next() transfer times
	GetNum    int    // GetGE() transfer times
}

func NewDebugs(debug *Debug) *Debugs {
	return &Debugs{
		DebugInfo: debug,
	}
}
