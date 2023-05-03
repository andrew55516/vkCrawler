-- name: GetLikesOnlyUsers :many
select distinct(u.label)
from (select ln.label as label
      from likes_edges le
               join likes_nodes ln on le.source = ln.id
      where le.start < $1
      order by le.start) as u
where not exists(
        select *
        from comments
        where owner = u.label
          and created_at < $1
    );

-- name: GetAmountOfUsers :one
select count(label) from all_nodes
where start < $1;