-- Modify "tenants" table
ALTER TABLE "public"."tenants" ADD COLUMN "name" character varying(150) NOT NULL, ADD COLUMN "logo" character varying(150) NULL, ADD COLUMN "phone" character varying(20) NOT NULL, ADD COLUMN "city" character varying(100) NOT NULL;
-- Create index "tenants_logo_key" to table: "tenants"
CREATE UNIQUE INDEX "tenants_logo_key" ON "public"."tenants" ("logo");
