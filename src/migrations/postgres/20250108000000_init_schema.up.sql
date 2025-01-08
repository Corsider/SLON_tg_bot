create schema app;

create table if not exists app.targets(
    creator_id bigint not null,
    target text not null,
    schedule int not null,
    tags text default null
);