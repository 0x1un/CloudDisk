create table if not exists filemetas (
    id serial not null primary key,
    file_md5 varchar(32) not null unique default '',
    file_name varchar(256) not null default '',
    file_size numeric(20) default '0',
    location varchar(1024) not null default '',
    upload_at timestamp default 'now',
    status smallint not null default '0',
    ext1 int default '0',
    ext2 text
);

create unique index file_idx_file_md5 on filemetas using btree(file_md5);
create index file_idx_status on filemetas using btree(status);


create table if not exists users (
    id serial not null primary key,
    user_name varchar(20) not null default '',
    user_pwd varchar(256) not null default '',
    email varchar(64) not null default '',
    phone varchar(11) default '' unique,
    email_validate boolean not null default false,
    phone_validate boolean default false,
    signup_at timestamp default 'now',
    last_active timestamp default 'now',
    profile varchar default '',
    status smallint not null default '0'
);

create index user_idx_status on users using btree(status);
create unique index user_idx_phone on users using btree(phone);
