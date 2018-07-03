
-- Sequence: public.vote_bot_id_seq

-- DROP SEQUENCE public.vote_bot_id_seq;

CREATE SEQUENCE public.vote_bot_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;
ALTER TABLE public.vote_bot_id_seq
  OWNER TO user;

-- Table: public.vote_bot

-- DROP TABLE public.vote_bot;

CREATE TABLE public.vote_bot
(
  id integer NOT NULL DEFAULT nextval('vote_bot_id_seq'::regclass),
  tg_id integer,
  user_name character varying(255),
  vote_variant integer,
  category character varying(255),
  updated_at integer,
  CONSTRAINT vote_bot_tg_id UNIQUE (tg_id),
  CONSTRAINT vote_bot_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.vote_bot OWNER TO user;
