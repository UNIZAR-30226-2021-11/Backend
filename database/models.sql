CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL,
    username VARCHAR(150) NOT NULL UNIQUE,
    password varchar(256) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    location varchar(150),
    games_won INT DEFAULT 0,
    games_lost INT DEFAULT 0,
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    CONSTRAINT pk_users PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS games (
    id serial NOT NULL,
    name VARCHAR(150) NOT NULL UNIQUE,
    public BOOLEAN NOT NULL,
    tournament BOOLEAN NOT NULL DEFAULT false,
    creation_date timestamp DEFAULT now(),
    end_date timestamp,
    CONSTRAINT pk_games PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS pairs (
    id serial NOT NULL,
    winned BOOLEAN DEFAULT false,
    game_points INT DEFAULT 0,
    game_id INT NOT NULL,
    CONSTRAINT pk_pairs PRIMARY KEY(id),
    CONSTRAINT fk_pairs_games FOREIGN KEY(game_id) REFERENCES games(id)
);

CREATE TABLE IF NOT EXISTS players (
    id serial NOT NULL,
    user_id INT NOT NULL,
    pair_id INT NOT NULL,
    CONSTRAINT pk_players PRIMARY KEY(id),
    CONSTRAINT fk_players_users FOREIGN KEY(user_id) REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_players_pairs FOREIGN KEY(pair_id) REFERENCES pairs(id)
);