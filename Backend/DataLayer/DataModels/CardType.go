package DataModels

type CardTypeDataModel struct {
	tableName struct{} `sql:"card_type"`
	Id        int      `sql:"id"`
	Name      string   `sql:"name"`
}
