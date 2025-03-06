package DataModels

type ExpansionDataModel struct {
	tableName struct{} `sql:"expansion"`
	Id        int      `sql:"id"`
	Name      string   `sql:"name"`
}
