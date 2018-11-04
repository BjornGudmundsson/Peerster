package data

//PrivateMessageStorage is a type that keeps track
//of all private message that this node has received.
//The underlying type is a map with the value being a
//a slice of strings. The keys are the origins nodes
//of the message and the values are the messages themselves
type PrivateMessageStorage map[string][]string

//GetMessagesFromOrigin is a function that allows the node to retrieve all
//the messages from a given origin.
func (priv *PrivateMessageStorage) GetMessagesFromOrigin(og string) []string {
	return (*priv)[og]
}

//PutMessageFromOrigin adds a message to the slice of messages known from a given origin
func (priv *PrivateMessageStorage) PutMessageFromOrigin(og string, msg string) {
	messages := (*priv)[og]
	messages = append(messages, msg)
	(*priv)[og] = messages
}
