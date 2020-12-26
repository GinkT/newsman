package db

import "github.com/mailru/dbr"

type SubscribeModel struct {
	Id string `db:"id"`
	UserId string `db:"user_id"`
	Tag string `db:"tag"`
	ReadenArticles []string `db:"readen_articles"`
}

func InsertNewSubscribeToSubscribes(s *dbr.Session, model SubscribeModel) (err error) {
	_, err = s.InsertInto("subscribes").
		Columns("user_id", "tag", "readen_articles").
		Record(&model).
		Exec()
	return
}

func DeleteSubscribeFromSubscribes(s *dbr.Session, userId string, tag string) (err error) {
	_, err = s.DeleteFrom("subscribes").
		Where("user_id = ? AND tag = ?", userId, tag).
		Exec()
	return
}

func SelectAllSubscribes(s *dbr.Session) (models []SubscribeModel, err error) {
	_, err = s.Select("id", "user_id", "tag", "readen_articles").
		From("subscribes").
		LoadStructs(&models)
	return
}

func AddReadenArticleToSubscribe(s *dbr.Session, userId, article string) (err error) {
	_, err = s.Update("subscribes").
		Set("readen_articles", "readen_articles || " + article).
		Where("user_id = ?", userId).
		Exec()
	return
}