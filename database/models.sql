CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL,
    username VARCHAR(150) NOT NULL UNIQUE,
    password varchar(256) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    year INT NOT NULL,
    admin BOOLEAN NOT NULL,
    picture VARCHAR(256) DEFAULT 'https://placekitten.com/g/300/300',
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    CONSTRAINT pk_users PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS subjects (
    id serial NOT NULL,
    name VARCHAR(150) NOT NULL,
    year INT NOT NULL,
    CONSTRAINT pk_subjects PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS posts (
    id serial NOT NULL,
    user_id int NOT NULL,
    subject_id int NOT NULL,
    title VARCHAR(150) NOT NULL,
    category VARCHAR(150) NOT NULL,
    body text NOT NULL,
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    CONSTRAINT pk_posts PRIMARY KEY(id),
    CONSTRAINT fk_posts_users FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_posts_subjects FOREIGN KEY(subject_id) REFERENCES subjects(id)
);

CREATE TABLE IF NOT EXISTS replies (
    id serial NOT NULL,
    user_id int NOT NULL,
    post_id int NOT NULL,
    body text NOT NULL,
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    CONSTRAINT pk_replies PRIMARY KEY(id),
    CONSTRAINT fk_replies_users FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_replies_posts FOREIGN KEY(post_id) REFERENCES posts(id)
);
