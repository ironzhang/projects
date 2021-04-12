package types

type Part struct {
	Part  string   `json:",omitempty"`
	Means []string `json:",omitempty"`
}

type Sound struct {
	EN  string `json:",omitempty"`
	AM  string `json:",omitempty"`
	TTS string `json:",omitempty"`
}

type Word struct {
	Word    string `json:",omitempty"`
	PhEN    string `json:",omitempty"`
	PhAM    string `json:",omitempty"`
	PhOther string `json:",omitempty"`
	Parts   []Part `json:",omitempty"`
	Sound   Sound  `json:",omitempty"`
}
