package DataLayer

import (
	"Architorture-Backend/DataLayer/DataModels"
	"Architorture-Backend/DataLayer/DataModels/CardTypeEnum"
	"github.com/go-pg/pg/orm"
)

func (this *DatabaseConnection) GetCardTypes() []DataModels.CardTypeDataModel {
	var models []DataModels.CardTypeDataModel
	err := this.db.Model(&models).Select()
	if err != nil {
		panic(err)
	}

	return models
}

func (this *DatabaseConnection) GetCardsWithoutArchitortureOrSaveCards(expansions []int) []DataModels.CardDataModel {
	var models []DataModels.CardDataModel
	err := this.db.Model(&models).
		Where("card_type_id <> ?", CardTypeEnum.Architorture).
		Where("id <> ?", 1).
		Where("id <> ?", 2).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			for _, expansion := range expansions {
				q = q.WhereOr("expansion_id = ?", expansion)
			}

			return q, nil
		}).Select()
	if err != nil {
		panic(err)
	}

	return models
}

func (this *DatabaseConnection) GetSaveCards(expansions []int) []DataModels.CardDataModel {
	var models []DataModels.CardDataModel
	err := this.db.Model(&models).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr("id = ?", 1)
			q = q.WhereOr("id = ?", 2)
			return q, nil
		}).
		Where("card_type_id = ?", CardTypeEnum.Save).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			for _, expansion := range expansions {
				q = q.WhereOr("expansion_id = ?", expansion)
			}

			return q, nil
		}).Select()
	if err != nil {
		panic(err)
	}

	return models
}

func (this *DatabaseConnection) GetArchitortureCards(expansions []int) []DataModels.CardDataModel {
	var models []DataModels.CardDataModel
	err := this.db.Model(&models).
		Where("card_type_id = ?", CardTypeEnum.Architorture).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			for _, expansion := range expansions {
				q = q.WhereOr("expansion_id = ?", expansion)
			}

			return q, nil
		}).Select()
	if err != nil {
		panic(err)
	}

	return models
}

func (this *DatabaseConnection) GetAvailableCardNames(expansions []int) []string {
	var models []DataModels.CardDataModel
	err := this.db.Model(&models).
		ColumnExpr("DISTINCT(name)").
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			for _, expansion := range expansions {
				q = q.WhereOr("expansion_id = ?", expansion)
			}

			return q, nil
		}).Select()
	if err != nil {
		panic(err)
	}

	cardNames := make([]string, len(models))
	for i, card := range models {
		cardNames[i] = card.Name
	}

	return cardNames
}

func (this *DatabaseConnection) GetAvailableCards(expansions []int) []DataModels.CardDataModel {
	var models []DataModels.CardDataModel
	err := this.db.Model(&models).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			for _, expansion := range expansions {
				q = q.WhereOr("expansion_id = ?", expansion)
			}

			return q, nil
		}).Select()
	if err != nil {
		panic(err)
	}

	return models
}
