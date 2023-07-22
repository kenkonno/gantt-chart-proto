CREATE TABLE public.stores
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    mail_address VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    store_name   VARCHAR(255) NOT NULL,
    domain       VARCHAR(255) NOT NULL,
    status       varchar(255) NOT NULL,
    retired_at   timestamp(0) with time zone,
    created_at   timestamp(0) with time zone,
    updated_at   timestamp(0) with time zone
);


-- Name        string `jaFieldName:"名前" json:"name" binding:"required"`
-- 	PhoneNumber string `jaFieldName:"電話番号" json:"phone_number" binding:"required"`
-- 	MailAddress string `jaFieldName:"メールアドレス" json:"email" binding:"required"`
-- 	CardNumber  int    `jaFieldName:"クレジットカード番号" json:"card_number" binding:"required"`
-- 	Domain      string `jaFieldName:"URL" json:"domain" binding:"required"`
-- 	CVS         int    `jaFieldName:"CVS" json:"cvs" binding:"required"`
-- 	Recapcha    string `jaFieldName:"Recapcha" json:"recapcha" binding:"required"`
-- 	Month       int    `jaFieldName:"月" json:"month" binding:"required"`
-- 	Year        int    `jaFieldName:"年" json:"year" binding:"required"`
-- 	ShopName    string `jaFieldName:"店舗名" json:"shop_name" binding:"required"`

-- -- ユーザーマスタ
-- CREATE TABLE public.users
-- (
--     id              SERIAL PRIMARY KEY,
--     phone_number    VARCHAR(20)  NOT NULL,
--     password        VARCHAR(255) NOT NULL,
--     name            VARCHAR(255) NOT NULL,
--     interview_sheet text,
--     created_at      timestamp(0) with time zone,
--     updated_at      timestamp(0) with time zone
-- );
--
-- -- 祝日
-- CREATE TABLE public.holidays
-- (
--     id         SERIAL PRIMARY KEY,
--     date       date,
--     memo       text,
--     created_at timestamp(0) with time zone,
--     updated_at timestamp(0) with time zone
-- );
--
-- -- ユーザーマスタ
-- CREATE TABLE public.admins
-- (
--     id         SERIAL PRIMARY KEY,
--     user_id    VARCHAR(20)  NOT NULL,
--     password   VARCHAR(255) NOT NULL,
--     name       VARCHAR(255) NOT NULL,
--     created_at timestamp(0) with time zone,
--     updated_at timestamp(0) with time zone
-- );
--
-- -- ユーザーマスタ
-- CREATE TABLE public.rooms
-- (
--     id           SERIAL PRIMARY KEY,
--     schema_name  VARCHAR(255) NOT NULL,
--     name         VARCHAR(255) NOT NULL,
--     start_hour   int          NOT NULL,
--     start_minute int          NOT NULL,
--     end_hour     int          NOT NULL,
--     end_minute   int          NOT NULL,
--     duration     int          NOT NULL,
--     created_at   timestamp(0) with time zone,
--     updated_at   timestamp(0) with time zone
-- );
--
