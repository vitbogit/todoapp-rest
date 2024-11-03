CREATE TYPE IF NOT EXISTS task_status AS ENUM ('frozen', 'pending', 'in-progress', 'completed');

CREATE TABLE IF NOT EXISTS tasks.tasks (
    id bigserial primary key,
    title text not null,
    description text null,
    status  task_status null,
    created_at timestamp without time zone not null,
    due_date timestamp without time zone null,
    updated_at timestamp without time zone not null
)

CREATE TABLE IF NOT EXISTS users.users (
    id bigserial primary key,
    login text unique not null,
    password text not null,
    f_name text not null,
    l_name text not null,
    role text not null,
    date_registration timestamp without time zone not null
)

CREATE TABLE IF NOT EXISTS tasks.comments (
    id bigserial primary key,
    id_user bigint not null,
    id_task bigint not null,
    content text not null,
    created_at timestamp without time zone not null,
    FOREIGN KEY (id_user) REFERENCES users.users(id),
    FOREIGN KEY (id_task) REFERENCES tasks.tasks(id)
)