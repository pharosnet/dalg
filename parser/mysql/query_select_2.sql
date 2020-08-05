-- name: 2
SELECT `users_domain_snapshot`.`id`,
       `users_domain_snapshot`.`aggregate_name`, /* dd */
       `users_domain_snapshot`.`aggregate_id`,
       `users_domain_snapshot`.`last_event_id`,
       `users_domain_snapshot`.`snapshot_data`,
       (`users_domain_snapshot`.`id` > 1)                                                          as `over`,
       (select count(`id`)
        from `ddd_test`.`users_domain_events`
        where `users_domain_events`.`aggregate_id` = `users_domain_snapshot`.`aggregate_id`)       as `count`,
       (select sum(`id`)
        from `ddd_test`.`users_domain_events`
        where `users_domain_events`.`aggregate_id` = `users_domain_snapshot`.`aggregate_id`)       as `sum`,
       exists(select `id`
              from `ddd_test`.`users_domain_events`
              where `users_domain_events`.`aggregate_id` = `users_domain_snapshot`.`aggregate_id`) as "x"
FROM `ddd_test`.`users_domain_snapshot`;