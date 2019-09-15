create table filemetas (
    id serial not null primary key,
    file_md5 varchar(32) not null unique default '',
    file_name varchar(256) not null default '',
    file_size numeric(20) default '0',
    location varchar(1024) not null default '',
    upload_at timestamp default 'now',
    status smallint not null default '0',
    ext1 smallint default '0',
    ext2 text
);
