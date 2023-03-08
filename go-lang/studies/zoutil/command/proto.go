package command

type TraceObject struct {
	Type    string
	Service string
}

type SetTraceVerboseArgs struct {
	TraceObject
	Verbose int
}
