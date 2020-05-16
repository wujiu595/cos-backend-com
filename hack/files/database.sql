--
-- PostgreSQL database dump
--

-- Dumped from database version 10.3
-- Dumped by pg_dump version 10.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: comunion; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA comunion;


--
-- Name: fake_id(text); Type: FUNCTION; Schema: comunion; Owner: -
--

CREATE FUNCTION comunion.fake_id(id text) RETURNS bigint
    LANGUAGE sql IMMUTABLE
    AS $_$
SELECT hex_to_int(concat(left(md5(substring(id from '(.*?)(?:-\d+)?$')), 10), coalesce(
    lpad(to_hex(substring(id from '.*-(\d+)$')::BIGINT), 4, '0')
)))
$_$;


--
-- Name: generate_date_series(timestamp with time zone, timestamp with time zone, text, text); Type: FUNCTION; Schema: comunion; Owner: -
--

CREATE FUNCTION comunion.generate_date_series(from_date timestamp with time zone, to_date timestamp with time zone, tz_time text, intval text) RETURNS SETOF timestamp with time zone
    LANGUAGE plpgsql IMMUTABLE
    AS $$
DECLARE
    r TIMESTAMPTZ;
BEGIN
    PERFORM set_config('timezone', 'UTC', FALSE);
    PERFORM set_config('timezone', (tz_time::timestamptz - tz_time::timestamp)::text, FALSE);
    FOR r IN SELECT
        date_trunc(intval, series)
    FROM generate_series(from_date , to_date-'1s'::INTERVAL, concat('1 ', intval)::interval) AS series LOOP
        RETURN NEXT r;
    END LOOP;
    RETURN;
END;
$$;


--
-- Name: id_generator(); Type: FUNCTION; Schema: comunion; Owner: -
--

CREATE FUNCTION comunion.id_generator() RETURNS bigint
    LANGUAGE sql
    AS $$
		SELECT
			(((EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT - 946684800000) << 22) |
			(1 << 12) |
			(nextval('global_id_sequence') % 4096)
$$;


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: categories; Type: TABLE; Schema: comunion; Owner: -
--

CREATE TABLE comunion.categories (
    id bigint DEFAULT comunion.id_generator() NOT NULL,
    name text NOT NULL,
    code text NOT NULL,
    source text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


--
-- Name: COLUMN categories.source; Type: COMMENT; Schema: comunion; Owner: -
--

COMMENT ON COLUMN comunion.categories.source IS 'startup';


--
-- Name: global_id_sequence; Type: SEQUENCE; Schema: comunion; Owner: -
--

CREATE SEQUENCE comunion.global_id_sequence
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: startups; Type: TABLE; Schema: comunion; Owner: -
--

CREATE TABLE comunion.startups (
    id bigint DEFAULT comunion.id_generator() NOT NULL,
    name text NOT NULL,
    uid bigint NOT NULL,
    mission text,
    logo text NOT NULL,
    tx_id text NOT NULL,
    block_num bigint,
    description_addr text NOT NULL,
    category_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    state integer DEFAULT 0 NOT NULL,
    is_iro boolean DEFAULT false NOT NULL
);


--
-- Name: COLUMN startups.state; Type: COMMENT; Schema: comunion; Owner: -
--

COMMENT ON COLUMN comunion.startups.state IS '0 创建中,1 已创建,2 未确认到tx产生,3 上链失败，4 已设置';


--
-- Name: users; Type: TABLE; Schema: comunion; Owner: -
--

CREATE TABLE comunion.users (
    id bigint DEFAULT comunion.id_generator() NOT NULL,
    wallet_addr text NOT NULL,
    public_secret text,
    private_secret text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_hunter boolean DEFAULT false NOT NULL
);


--
-- Name: categories categories_id_pk; Type: CONSTRAINT; Schema: comunion; Owner: -
--

ALTER TABLE ONLY comunion.categories
    ADD CONSTRAINT categories_id_pk PRIMARY KEY (id);


--
-- Name: startups startups_id_pk; Type: CONSTRAINT; Schema: comunion; Owner: -
--

ALTER TABLE ONLY comunion.startups
    ADD CONSTRAINT startups_id_pk PRIMARY KEY (id);


--
-- Name: users users_id_pk; Type: CONSTRAINT; Schema: comunion; Owner: -
--

ALTER TABLE ONLY comunion.users
    ADD CONSTRAINT users_id_pk PRIMARY KEY (id);


--
-- Name: categories_code; Type: INDEX; Schema: comunion; Owner: -
--

CREATE UNIQUE INDEX categories_code ON comunion.categories USING btree (code);


--
-- Name: categories_name; Type: INDEX; Schema: comunion; Owner: -
--

CREATE UNIQUE INDEX categories_name ON comunion.categories USING btree (name);


--
-- Name: start_ups_tx_id; Type: INDEX; Schema: comunion; Owner: -
--

CREATE UNIQUE INDEX start_ups_tx_id ON comunion.startups USING btree (tx_id);


--
-- Name: users_wallet_addr; Type: INDEX; Schema: comunion; Owner: -
--

CREATE UNIQUE INDEX users_wallet_addr ON comunion.users USING btree (wallet_addr);


--
-- PostgreSQL database dump complete
--

