/* Create the hangman database, schema and tables.*/

CREATE DATABASE hangman;

\connect hangman;

CREATE SCHEMA IF NOT EXISTS hangman;

CREATE TABLE hangman.games (
    uuid    varchar(100) CONSTRAINT uuid PRIMARY KEY,
    turns_left integer NOT NULL,
    word    varchar(100) NOT NULL,
    used    varchar(100),
    available_hints integer NOT NULL
);
