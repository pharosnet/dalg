-- name: 1
SELECT `ee`.`id` as `xxxx`,
       `ee`.`aggregate_name`,
       `users_domain_events`.`aggregate_id`,
       `ee`.`event_name`,
       `ee`.`event_id`
FROM `ddd_test`.`users_domain_events` as `ee`
WHERE `ee`.`aggregate_id` = ?
  AND `ee`.`aggregate_name` = 'DD'
  AND `ee`.`event_id` IN ('#xxxx#')
  and `ee`.`event_name` between ? and ?
ORDER BY `id` DESC
LIMIT ? OFFSET ?;