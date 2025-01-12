create schema app;

create table if not exists app.targets(
    creator_id bigint not null,
    target text not null,
    schedule int not null,
    tags text default null
);

update app.targets
    set tags = '{}';

alter table app.targets
alter column tags TYPE int[] USING '{}'::INT[];