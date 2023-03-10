package models

type BlogEntry struct {
	Id                      int      `json:"id"`
	OriginalLocale          string   `json:"originalLocale"`
	CreationTimeSeconds     int      `json:"creationTimeSeconds"`
	AuthorHandle            string   `json:"authorHandle"`
	Title                   string   `json:"title"`
	Content                 string   `json:"content"`
	Locale                  string   `json:"locale"`
	ModificationTimeSeconds int      `json:"modificationTimeSeconds"`
	AllowViewHistory        bool     `json:"allowViewHistory"`
	Tags                    []string `json:"tags"`
	Rating                  int      `json:"rating"`
}

type Comment struct {
	Id                  int    `json:"id" bson:"id"`
	CreationTimeSeconds int    `json:"creationTimeSeconds" bson:"creationTimeSeconds"`
	CommentatorHandle   string `json:"commentatorHandle" bson:"commentatorHandle"`
	Locale              string `json:"locale" bson:"locale"`
	Text                string `json:"text" bson:"text"`
	ParentCommentId     int    `json:"parentCommentId" bson:"parentCommentId"`
	Rating              int    `json:"rating" bson:"rating"`
}

type RecentAction struct {
	TimeSeconds int64      `json:"timeSeconds" bson:"timeSeconds"`
	BlogEntry   *BlogEntry `json:"blogEntry" bson:"blogEntry"`
	Comment     *Comment   `json:"comment" bson:"comment"`
}

type User struct {
	Username         string `json:"username" bson:"username"`
	Email            string `json:"email" bson:"email"`
	CodeforcesHandle string `json:"codeforcesHandle" bson:"codeforcesHandle"'`
	SubscribedBlogs  []int  `json:"subscribedBlogs" bson:"subscribedBlogs"`
}
