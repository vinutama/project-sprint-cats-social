package cat_entity

import "time"

type Cat struct {
	Id          string
	Name        string
	Race        string
	Sex         string
	AgeInMonth  int
	UserId      string
	Description string
	ImageURLs   []string `db:"image_urls"`
	HasMatched  bool
	IsDeleted   bool
	CreatedAt   time.Time
}

type CatSearch struct {
	Id           string
	Name         string
	Race         string
	Sex          string
	AgeInMonth   int
	AgeCondition string
	HasMatched   string
	Owned        string
	UserId       string
	Limit        int
	Offset       int
}
