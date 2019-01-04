create table  if not exists users (
    id int unsigned not null primary key auto_increment,
    no varchar(15) not null default '',
    name nvarchar(255) not null default '',
    group_id varchar(15) not null default '',
    create_at timestamp not null default now(),
    create_by varchar(15) not null default '',
    constraint uc_users_no unique(no)
);

create table  if not exists  groups (
    id int unsigned not null primary key auto_increment,
    no varchar(15) not null default '',
    name nvarchar(255) not null default '',
    create_at timestamp not null default now(),
    create_by varchar(15) not null default '',
    constraint uc_users_no unique(no) 
);
    
create table if not exists goods_images(
    
    id  int unsigned not null  primary key auto_increment,
    goods_id int unsigned not null default 0,
    uuid varchar(50) not null default '',
    size int unsigned not null default 0,
    url varchar(255) not null default '',
    create_at timestamp not null default now(),
    create_by varchar(15) not null default '',
    constraint uc_goods_images_id unique(goods_id, uuid)
);
    

