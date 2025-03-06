package DataModels

type CardDataModel struct {
	tableName       struct{} `sql:"card"`
	Id              int      `sql:"id, pk"`
	CardTypeId      int      `sql:"card_type_id"`
	Name            string   `sql:"name"`
	Description     string   `sql:"description"`
	PlayImmediately bool     `sql:"play_immediately"`
	Quantity        int      `sql:"quantity"`
	ExpansionId     int      `sql:"expansion_id"`
	Archivable      bool     `sql:"archivable"`
}
