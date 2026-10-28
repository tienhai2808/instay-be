package translator

import "gorm.io/gorm"

const (
	gormPostgresTranslatorCallbackName       = "instay:postgres_error_translator"
	gormPostgresTranslatorDeleteCallbackName = "instay:postgres_error_translator_delete"
)

func RegisterGORMPostgresTranslator(db *gorm.DB) error {
	if err := db.Callback().Create().After("gorm:create").Register(gormPostgresTranslatorCallbackName, gormTranslateError); err != nil {
		return err
	}
	if err := db.Callback().Query().After("gorm:query").Register(gormPostgresTranslatorCallbackName, gormTranslateError); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:update").Register(gormPostgresTranslatorCallbackName, gormTranslateError); err != nil {
		return err
	}
	if err := db.Callback().Delete().After("gorm:delete").Register(gormPostgresTranslatorDeleteCallbackName, gormTranslateDeleteError); err != nil {
		return err
	}
	if err := db.Callback().Row().After("gorm:row").Register(gormPostgresTranslatorCallbackName, gormTranslateError); err != nil {
		return err
	}
	if err := db.Callback().Raw().After("gorm:raw").Register(gormPostgresTranslatorCallbackName, gormTranslateError); err != nil {
		return err
	}
	return nil
}

func gormTranslateError(db *gorm.DB) {
	db.Error = PostgresTranslate(db.Error)
}

func gormTranslateDeleteError(db *gorm.DB) {
	db.Error = PostgresTranslateDelete(db.Error)
}
