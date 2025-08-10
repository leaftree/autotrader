create table positions (
    id integer(20) not null,
    contract varchar(32) not null,
    size int,
    price number,
    status int comment '1: done',
) ;