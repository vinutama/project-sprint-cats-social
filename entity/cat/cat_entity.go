package cat_entity

type Cat struct {
	Id          string
	Name        string
	Race        string
	Sex         string
	AgeInMonth  int
	UserId      string
	Description string
	ImageURLs   []string
	HasMatched  bool
	IsDeleted   bool
	CreatedAt   string
}
