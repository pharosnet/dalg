-- name: 1
INSERT INTO `ddd_test`.`users_domain_events`
(`id`,
 `aggregate_name`,
 `aggregate_id`,
 `event_name`,
 `event_id`,
 `event_time`,
 `event_data`)
VALUES
(?,?,?,?,?,now(),'');