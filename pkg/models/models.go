package models

type BlogEntry struct {
	Id                      int    `json:"id"`
	OriginalLocale          string `json:"originalLocale"`
	CreationTimeSeconds     int    `json:"creationTimeSeconds"`
	AuthorHandle            string `json:"authorHandle"`
	Title                   string `json:"title"`
	Content                 string `json:"content"`
	Locale                  string `json:"locale"`
	ModificationTimeSeconds int    `json:"modificationTimeSeconds"`
	AllowViewHistory        bool   `json:"allowViewHistory"`
	Tags                    string `json:"tags"`
	Rating                  int    `json:"rating"`
}

type Comment struct {
	Id                  int    `json:"id"`
	CreationTimeSeconds int    `json:"creationTimeSeconds"`
	CommentatorHandle   string `json:"commentatorHandle"`
	Locale              string `json:"locale"`
	Text                string `json:"text"`
	ParentCommentId     int    `json:"parentCommentId"`
	Rating              int    `json:"rating"`
}

type RecentAction struct {
	TimeSeconds int      `json:"timeSeconds"`
	BlogEntry   struct{} `json:"blogEntry"`
	Comment     struct{} `json:"comment"`
}
