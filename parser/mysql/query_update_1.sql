-- name: 1
UPDATE `ddd_test`.`users_domain_events`
SET
    `aggregate_name` = ?,
    `aggregate_id` = ?,
    `event_name` = '1',
    `event_id` = ?,
    `event_time` = now(),
    `event_data` = null
WHERE `id` = ?;