package storage

// Post - публикация.
type Post struct {
	ID          int    `bson:"_id"`
	Title       string `bson:"title"`
	Content     string `bson:"content"`
	AuthorID    int    `bson:"authorid"`
	AuthorName  string `bson:"authorname"`
	CreatedAt   int64  `bson:"createdat"`
	PublishedAt int64  `bson:"publishedat"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
}
