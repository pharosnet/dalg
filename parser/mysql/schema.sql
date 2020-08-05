USE `ddd_test`;

-- name: users_domain_events
CREATE TABLE `users_domain_events`
(
    `id`             bigint       NOT NULL AUTO_INCREMENT,
    `aggregate_name` varchar(255) NOT NULL, -- name: AggName ref: github.com/foo/bar.SQLString
    `aggregate_id`   varchar(255) NOT NULL,
    `event_name`     varchar(255) NOT NULL,
    `event_id`       varchar(63)  NOT NULL,
    `event_time`     datetime(6) DEFAULT NULL,
    `event_data`     text,
    PRIMARY KEY (`id`),
    UNIQUE KEY `users_ix_event_id` (`event_id`),
    KEY `users_ix_aggname` (`aggregate_name`),
    KEY `users_ix_aggid` (`aggregate_id`),
    KEY `users_ix_event_name` (`event_name`),
    KEY `users_ix_event_time` (`event_time` DESC)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE `users_domain_snapshot`
(
    `id`             bigint NOT NULL AUTO_INCREMENT,
    `aggregate_name` varchar(255) DEFAULT NULL,
    `aggregate_id`   varchar(63)  DEFAULT NULL,
    `last_event_id`  varchar(63)  DEFAULT NULL,
    `snapshot_data`  text,
    PRIMARY KEY (`id`),
    KEY `users_ix_snapshot_last_eid` (`last_event_id`),
    KEY `users_ix_snapshot_aggname` (`aggregate_name`),
    KEY `users_ix_snapshot_aggid` (`aggregate_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;
