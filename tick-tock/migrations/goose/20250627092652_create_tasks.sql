-- +goose Up
-- +goose StatementBegin
create table if not exists task
(
    id          bigint auto_increment
        primary key,
    app         varchar(64) default ''                not null,
    tid         varchar(36)                           not null,
    output      text                                  not null comment '执行结果',
    run_time    timestamp                             not null comment '执行时间',
    cost_time   bigint                                not null comment '执行耗时（毫秒）',
    status      tinyint                               not null comment '当前状态 未执行 0 执行中 1 执行成功 2 执行失败 3',
    create_time timestamp   default CURRENT_TIMESTAMP not null,
    update_time timestamp   default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table task;
-- +goose StatementEnd
