--
-- PostgreSQL database dump
--

\restrict jbjVyI78SeXuqOASmBDDQaeU0XsezdSXaxiMwaxLu2AAJcgmFv2jHzum2cX8IN6

-- Dumped from database version 15.17 (Debian 15.17-1.pgdg13+1)
-- Dumped by pg_dump version 18.3 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: auth_identities; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.auth_identities (
    id integer NOT NULL,
    user_id integer,
    provider text NOT NULL,
    provider_id text NOT NULL,
    created_id timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.auth_identities OWNER TO admin;

--
-- Name: auth_identities_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.auth_identities_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.auth_identities_id_seq OWNER TO admin;

--
-- Name: auth_identities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.auth_identities_id_seq OWNED BY public.auth_identities.id;


--
-- Name: expenses; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.expenses (
    id integer NOT NULL,
    group_id integer,
    paid_by integer,
    amount numeric(15,2),
    description text,
    category text,
    receipt_image text,
    split_type text DEFAULT 'equal'::text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.expenses OWNER TO admin;

--
-- Name: expenses_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.expenses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.expenses_id_seq OWNER TO admin;

--
-- Name: expenses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.expenses_id_seq OWNED BY public.expenses.id;


--
-- Name: group_members; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.group_members (
    group_id integer NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.group_members OWNER TO admin;

--
-- Name: groups; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.groups (
    id integer NOT NULL,
    name text NOT NULL,
    created_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.groups OWNER TO admin;

--
-- Name: groups_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.groups_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.groups_id_seq OWNER TO admin;

--
-- Name: groups_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.groups_id_seq OWNED BY public.groups.id;


--
-- Name: participants; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.participants (
    expense_id integer,
    user_id integer,
    share_amount numeric(15,2)
);


ALTER TABLE public.participants OWNER TO admin;

--
-- Name: users; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name text,
    password text,
    email text,
    profile_pic text
);


ALTER TABLE public.users OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: auth_identities id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.auth_identities ALTER COLUMN id SET DEFAULT nextval('public.auth_identities_id_seq'::regclass);


--
-- Name: expenses id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.expenses ALTER COLUMN id SET DEFAULT nextval('public.expenses_id_seq'::regclass);


--
-- Name: groups id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.groups ALTER COLUMN id SET DEFAULT nextval('public.groups_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: auth_identities; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.auth_identities (id, user_id, provider, provider_id, created_id) FROM stdin;
1	1	local	alice@test.com	2026-04-28 00:09:36.205626
2	2	local	bob@test.com	2026-04-28 00:09:36.205626
3	3	local	charlie@test.com	2026-04-28 00:09:36.205626
4	6	google	102587030244501572972	2026-04-27 18:41:48.865114
5	7	google	112446934825321314394	2026-04-27 20:24:30.042754
7	8	google	s	2026-04-28 17:14:43.458959
8	7	github	97462628	2026-04-28 17:15:59.125858
\.


--
-- Data for Name: expenses; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.expenses (id, group_id, paid_by, amount, description, category, receipt_image, split_type, created_at) FROM stdin;
1	1	1	900.00	Airbnb Payment	Accommodation	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	equal	2026-04-25 22:20:18.809521
2	2	1	150.00	Thai Food	Food	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	equal	2026-04-26 03:52:00.96504
3	2	2	3000.00	Monthly Rent	Housing	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	manual	2026-04-26 03:52:00.96504
4	3	3	45.00	Trail Snacks	Food	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	equal	2026-04-26 03:52:00.96504
5	4	1	60.00	Tickets	Entertainment	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	equal	2026-04-26 03:52:00.96504
6	5	2	120.00	Electricity	Bills	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	equal	2026-04-26 03:52:00.96504
7	6	3	90.00	Protein Powder	Health	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	manual	2026-04-26 03:52:00.96504
8	7	1	500.00	Fuel & Car Wash	Transport	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	equal	2026-04-26 03:52:00.96504
9	8	2	150.00	Gift Box	Gifts	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	equal	2026-04-26 03:52:00.96504
10	9	3	210.00	Wine & Cheese	Food	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	equal	2026-04-26 03:52:00.96504
11	10	1	85.00	Weekly Veggies	Food	https://cdn-icons-png.flaticon.com/512/3135/3135679.png	equal	2026-04-26 03:52:00.96504
\.


--
-- Data for Name: group_members; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.group_members (group_id, user_id) FROM stdin;
1	1
1	2
1	3
2	1
2	2
2	3
3	1
3	2
3	3
4	1
4	2
4	3
5	1
5	2
5	3
6	1
6	2
6	3
7	1
7	2
7	3
8	1
8	2
8	3
9	1
9	2
9	3
10	1
10	2
10	3
11	1
11	2
11	3
\.


--
-- Data for Name: groups; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.groups (id, name, created_by, created_at) FROM stdin;
1	Goa Trip	1	2026-04-25 22:19:46.032925
2	Office Lunch	1	2026-04-26 03:52:00.96504
3	House Rent	2	2026-04-26 03:52:00.96504
4	Weekend Hike	3	2026-04-26 03:52:00.96504
5	Movie Night	1	2026-04-26 03:52:00.96504
6	Electricity Bill	2	2026-04-26 03:52:00.96504
7	Gym Membership	3	2026-04-26 03:52:00.96504
8	Road Trip	1	2026-04-26 03:52:00.96504
9	Wedding Gift	2	2026-04-26 03:52:00.96504
10	Dinner Party	3	2026-04-26 03:52:00.96504
11	Groceries	1	2026-04-26 03:52:00.96504
\.


--
-- Data for Name: participants; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.participants (expense_id, user_id, share_amount) FROM stdin;
1	1	300.00
1	2	300.00
1	3	300.00
2	1	50.00
2	2	50.00
2	3	50.00
3	1	2000.00
3	3	1000.00
4	1	15.00
4	2	15.00
4	3	15.00
5	1	20.00
5	2	20.00
5	3	20.00
6	1	40.00
6	2	40.00
6	3	40.00
7	2	90.00
8	1	166.66
8	2	166.67
8	3	166.67
9	1	50.00
9	2	50.00
9	3	50.00
10	1	70.00
10	2	70.00
10	3	70.00
11	1	28.33
11	2	28.33
11	3	28.34
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.users (id, name, password, email, profile_pic) FROM stdin;
1	Alice	$2a$10$0qRKLe7yI3LF0opb2QUkku8laGDJnKwYttvlnAFd3qetBHtiI.dH6	alice@test.com	https://cdn-icons-png.flaticon.com/512/4140/4140048.png
2	Bob	$2a$10$hbupj2TYkRSNCA3Jqu.tw.zYMzD8DuJJvLVDYgriQhu8w3LGGooQe	bob@test.com	https://cdn-icons-png.flaticon.com/512/4140/4140048.png
3	Charlie	$2a$10$8pqMdCExIu3Y1NIdnqnR7uz7uquO4Iq4h.vy878hBfNiO/ZcuoFHa	charlie@test.com	https://cdn-icons-png.flaticon.com/512/4140/4140048.png
6	littlelaughterlabs	\N	littlelaughterlab@gmail.com	https://lh3.googleusercontent.com/a/ACg8ocLaMCo-Shvg2FQaQMzUaKoKMIZwaXt0boYOMBo5aUPzdgBC188=s96-c
7	Gaurav Sahay	\N	gauravsahay2468@gmail.com	https://lh3.googleusercontent.com/a/ACg8ocL1qvX4rhW6lCs1vkvUhEjLl9dqCMBT6DrUSYqJ5UUOEW28uA=s96-c
8	Gaurav Sahay	\N	gauravsahay54321@gmail.com	https://lh3.googleusercontent.com/a/ACg8ocJH-fqQ2NhA-TNFmnwZetbZ0XYT06rE7Glzd-Tx8nlIn30nxuM=s96-c
\.


--
-- Name: auth_identities_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.auth_identities_id_seq', 8, true);


--
-- Name: expenses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.expenses_id_seq', 11, true);


--
-- Name: groups_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.groups_id_seq', 11, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.users_id_seq', 8, true);


--
-- Name: auth_identities auth_identities_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.auth_identities
    ADD CONSTRAINT auth_identities_pkey PRIMARY KEY (id);


--
-- Name: auth_identities auth_identities_provider_provider_id_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.auth_identities
    ADD CONSTRAINT auth_identities_provider_provider_id_key UNIQUE (provider, provider_id);


--
-- Name: expenses expenses_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.expenses
    ADD CONSTRAINT expenses_pkey PRIMARY KEY (id);


--
-- Name: group_members group_members_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.group_members
    ADD CONSTRAINT group_members_pkey PRIMARY KEY (group_id, user_id);


--
-- Name: groups groups_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT groups_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_expenses_paid_by; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_expenses_paid_by ON public.expenses USING btree (paid_by);


--
-- Name: idx_group_members_user_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_group_members_user_id ON public.group_members USING btree (user_id);


--
-- Name: idx_participants_expense_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_participants_expense_id ON public.participants USING btree (expense_id);


--
-- Name: auth_identities auth_identities_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.auth_identities
    ADD CONSTRAINT auth_identities_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: expenses expenses_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.expenses
    ADD CONSTRAINT expenses_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.groups(id);


--
-- Name: expenses expenses_paid_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.expenses
    ADD CONSTRAINT expenses_paid_by_fkey FOREIGN KEY (paid_by) REFERENCES public.users(id);


--
-- Name: group_members group_members_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.group_members
    ADD CONSTRAINT group_members_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.groups(id);


--
-- Name: group_members group_members_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.group_members
    ADD CONSTRAINT group_members_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: groups groups_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT groups_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: participants participants_expense_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.participants
    ADD CONSTRAINT participants_expense_id_fkey FOREIGN KEY (expense_id) REFERENCES public.expenses(id) ON DELETE CASCADE;


--
-- Name: participants participants_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.participants
    ADD CONSTRAINT participants_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict jbjVyI78SeXuqOASmBDDQaeU0XsezdSXaxiMwaxLu2AAJcgmFv2jHzum2cX8IN6

