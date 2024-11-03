CREATE OR REPLACE FUNCTION tasks.tasks_create(
    _title text,
    _description text,
    _status text,
    _due_date timestamp without time zone
)
returns bigint
language plpgsql
as
$$
    DECLARE id_task bigint;
begin
    insert into tasks.tasks (title, description, status, created_at, due_date, updated_at) 
        values (_title, _description, _status::task_status, NOW(), _due_date, NOW()) 
        returning id into id_task;
    
    return id_task;
end;
$$;

CREATE OR REPLACE FUNCTION tasks.tasks_list()
returns table (
    id bigint,
    title text,
    description text,
    status text,
    created_at timestamp without time zone,
    due_date timestamp without time zone,
    updated_at timestamp without time zone
)
language plpgsql
as
$$
begin
    return query
        SELECT t.id, t.title, t.description, t.status::text, t.created_at, t.due_date, t.updated_at from tasks.tasks t;
end;
$$;

CREATE OR REPLACE PROCEDURE tasks.tasks_update(
    _id bigint,
    _title text,
    _description text,
    _status text,
    _due_date timestamp without time zone
)
language plpgsql
as
$$
begin
    update tasks.tasks set (title, description, status, due_date, updated_at) =
        (_title, _description, _status::task_status, _due_date, NOW()) 
        where id=_id;
end;
$$;

CREATE OR REPLACE PROCEDURE tasks.tasks_delete(
    _ids bigint[]
)
language plpgsql
as
$$
begin
    delete from tasks.tasks where id=any(_ids);
end;
$$;

--------------------------------

CREATE OR REPLACE FUNCTION users.auth(
    _login text
)
returns table (
    password text,
    role text
)
language plpgsql
as
$$
begin
    return query
        SELECT u.password, u.role from users.users u where u.login=_login;
end;
$$;

CREATE OR REPLACE FUNCTION users.register(
    _login text,
    _password text,
    _role text,
    _fname text,
    _lname text
)
returns bigint
language plpgsql
as
$$
    DECLARE id_user bigint;
begin
    insert into users.users (login, password, role, f_name, l_name, date_registration) values (_login, _password, _role, _fname, _lname, NOW()) returning id into id_user;
    return id_user;
end;
$$;

--------------------------------

CREATE OR REPLACE FUNCTION tasks.comments_list(
    _id_task bigint
)
returns table (
    id bigint,
    id_user bigint,
    id_task bigint,
    content text,
    created_at timestamp without time zone
)
language plpgsql
as
$$
begin
    return query
        SELECT c.id, c.id_user, c.id_task, c.content, c.created_at from tasks.comments c where c.id_task=_id_task;
end;
$$;

CREATE OR REPLACE FUNCTION tasks.comment_create(
    _id_task bigint,
    _login text,
    _content text
)
returns bigint
language plpgsql
as
$$
    DECLARE _id_user bigint;
    DECLARE _id_comment bigint;
begin

    select u.id from users.users u where u.login=_login into _id_user;

    if _id_user is null then 
        raise exception 'not found user with given login';
    end if;

    insert into tasks.comments (id_user, id_task, content, created_at)
        values (_id_user, _id_task, _content, NOW())
        returning id into _id_comment;
    
    return _id_comment;
end;
$$;
