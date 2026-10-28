variable "db_url" {
  type    = string
  default = getenv("DATABASE_URL")
}

variable "db_dev_url" {
  type    = string
  default = getenv("DATABASE_DEV_URL")
}

data "external_schema" "gorm" {
  program = ["./tmp/migrate"]
}

env "local" {
  src = data.external_schema.gorm.url
  url = var.db_url
  dev =  var.db_dev_url

  migration {
    dir = "file://migrations"
    format = atlas
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "prod" {
  url = var.db_url

  migration {
    dir    = "file:///app/migrations"
    format = atlas
  }
}