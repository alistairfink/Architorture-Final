CREATE TABLE public.expansion
(
    id integer NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT expansion_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.expansion
    OWNER to postgres;

CREATE TABLE public.card_type
(
    id integer NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT card_type_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.card_type
    OWNER to postgres;

CREATE TABLE public.card
(
    id integer NOT NULL,
    card_type_id integer NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    play_immediately boolean NOT NULL,
    quantity integer NOT NULL,
    expansion_id integer NOT NULL,
    archivable boolean NOT NULL,
    CONSTRAINT card_pkey PRIMARY KEY (id),
    CONSTRAINT card_type_id_fk FOREIGN KEY (card_type_id)
        REFERENCES public.card_type (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT expansion_id_fk FOREIGN KEY (expansion_id)
        REFERENCES public.expansion (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE public.card
    OWNER to postgres;