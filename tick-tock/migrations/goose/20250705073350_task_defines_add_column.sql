-- +goose Up
-- +goose StatementBegin
alter table task_defines
    add last_migrate_time timestamp default CURRENT_TIMESTAMP not null after notify_http_param;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE task_defines
DROP
COLUMN last_migrate_time;
-- +goose StatementEnd
