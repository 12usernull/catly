package engine

type IndexEntry struct {
	Segment int
	Offset  int64
	Length  int
	User    string
	Time    int64
}

var Index []IndexEntry
