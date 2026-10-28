package main

import (
	"fmt"
	"log"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	_ "gorm.io/driver/postgres"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&entity.TenantEntity{},
		&entity.PlanEntity{},
		&entity.SubscriptionEntity{},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(stmts)
}
