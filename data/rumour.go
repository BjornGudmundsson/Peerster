package data

//RumourMessage keeps holds of the content, original sender
//and the corresponding ID for a rumour in a struct
type RumourMessage struct {
	Origin string
	ID     uint32
	Text   string
}
