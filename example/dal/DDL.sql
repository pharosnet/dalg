-- Table: public."USER"

DROP TABLE IF EXISTS public."USER";

CREATE TABLE public."USER"
(
	"ID" varchar(64)  COLLATE pg_catalog."default", 
	"NAME" varchar(255)  COLLATE pg_catalog."default", 
	"AGE" numeric(18) , 
	"SEX" boolean , 
	"MONEY" numeric(18,2) , 
	"INFO" jsonb , 
	"CREATE_BY" varchar(255)  COLLATE pg_catalog."default", 
	"CREATE_TIME" timestamp , 
	"MODIFY_BY" varchar(255)  COLLATE pg_catalog."default", 
	"MODIFY_TIME" timestamp , 
	"VERSION" bigint , 
	CONSTRAINT "USER_PK" PRIMARY KEY ("ID")
)

WITH (
	OIDS = FALSE
)

TABLESPACE pg_default;

ALTER TABLE public."USER" OWNER to postgres;

-- Index: USER_IDX_NAME_AGE

CREATE  INDEX "USER_IDX_NAME_AGE" ON public."USER" USING btree ("NAME" COLLATE pg_catalog."default"  desc, "AGE" desc) TABLESPACE pg_default;

-- Index: USER_IDX_CREATE_TIME

CREATE  INDEX "USER_IDX_CREATE_TIME" ON public."USER" USING btree ("CREATE_TIME" desc) TABLESPACE pg_default;

