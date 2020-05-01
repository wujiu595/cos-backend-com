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


--
-- PostgreSQL database dump complete
--

