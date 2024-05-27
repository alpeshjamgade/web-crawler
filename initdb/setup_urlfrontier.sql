-- set password
ALTER USER postgres WITH PASSWORD 'postgres';

-- create database urlfrontier
CREATE DATABASE urlfrontier;

-- create table seeds
create table seeds
(
    url varchar not null constraint table_name_pk primary key,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);

