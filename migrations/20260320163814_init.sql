-- Create "plans" table
CREATE TABLE "public"."plans" (
  "id" bigserial NOT NULL,
  "name" character varying(150) NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "currency" character varying(20) NOT NULL,
  "features" jsonb NULL,
  "description" text NOT NULL,
  "duration" character varying(50) NOT NULL,
  "duration_days" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_plans_currency" CHECK ((currency)::text = ANY ((ARRAY['usd'::character varying, 'vnd'::character varying])::text[]))
);
-- Create index "plans_name_key" to table: "plans"
CREATE UNIQUE INDEX "plans_name_key" ON "public"."plans" ("name");
-- Create "tenants" table
CREATE TABLE "public"."tenants" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "domain" character varying(255) NOT NULL,
  "schema_name" character varying(50) NOT NULL,
  "settings" jsonb NULL,
  "is_active" boolean NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "tenants_domain_key" to table: "tenants"
CREATE UNIQUE INDEX "tenants_domain_key" ON "public"."tenants" ("domain");
-- Create index "tenants_schema_name_key" to table: "tenants"
CREATE UNIQUE INDEX "tenants_schema_name_key" ON "public"."tenants" ("schema_name");
-- Create "subscriptions" table
CREATE TABLE "public"."subscriptions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "tenant_id" bigint NOT NULL,
  "plan_id" bigint NOT NULL,
  "start_at" timestamptz NOT NULL,
  "end_at" timestamptz NOT NULL,
  "is_using" boolean NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_subscriptions_plan" FOREIGN KEY ("plan_id") REFERENCES "public"."plans" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "fk_subscriptions_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE CASCADE ON DELETE RESTRICT
);
