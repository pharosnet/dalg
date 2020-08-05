-- name: 3
SELECT e.id, s.aggregate_id, s1.aggregate_id
FROM ddd_test.users_domain_events as e
         left join ddd_test.users_domain_snapshot as s on s.aggregate_id = e.aggregate_id
         left join ddd_test.users_domain_snapshot as s1
                   on s1.aggregate_id = e.aggregate_id;