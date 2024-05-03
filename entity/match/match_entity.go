package match_entity

type Match struct {
	// Status status: 'requested', 'approved', 'rejected'
	Id, Message, CatIssuerId, CatReceiverId, Status, CreatedAt string
}
