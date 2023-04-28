-- name: FillAllNodes :exec
INSERT INTO all_nodes (label, start)
SELECT owner, min(created_at) as min_created
from (select owner, created_at
      from comments
      union
      select thread_owner as owner, created_at
      from comments
      union
      select owner, created_at
      from posts
      union
      select l.owner, p.created_at
      from likes l
               join posts p on l.post_id = p.id) as owners
group by owner;

-- name: FillAllEdges :exec
INSERT INTO all_edges (source, target, start)
select source, target, u.start
from (select t.id as source, s.id as target, start_time as start
      from all_nodes s
               join (select l.owner as like_owner, p.created_at as start_time, p.owner as post_owner
                     from likes l
                              join posts p on l.post_id = p.id) l on post_owner = s.label
               join all_nodes t on t.label = like_owner
      union all
      select t.id as source, s.id as target, c.created_at as start
      from all_nodes s
               join comments c on s.label = c.thread_owner
               join all_nodes t on t.label = c.owner) as u;

-- name: FillLikesNodes :exec
INSERT INTO likes_nodes (label, start)
SELECT owner, min(created_at) as min_created
from (select owner, created_at
      from posts
      union
      select l.owner, p.created_at
      from likes l
               join posts p on l.post_id = p.id) as owners
group by owner;

-- name: FillLikesEdges :exec
INSERT INTO likes_edges (source, target, start)
select source, target, u.start
from (select t.id as source, s.id as target, start_time as start
      from likes_nodes s
               join (select l.owner as like_owner, p.created_at as start_time, p.owner as post_owner
                     from likes l
                              join posts p on l.post_id = p.id) l on post_owner = s.label
               join likes_nodes t on t.label = like_owner) as u;

-- name: FillCommentsNodes :exec
INSERT INTO comments_nodes (label, start)
SELECT owner, min(created_at) as min_created
from (select owner, created_at
      from comments
      union
      select thread_owner as owner, created_at
      from comments
      union
      select owner, created_at
      from posts) as owners
group by owner;

-- name: FillCommentsEdges :exec
INSERT INTO comments_edges (source, target, start)
select source, target, u.start
from (select t.id as source, s.id as target, c.created_at as start
      from comments_nodes s
               join comments c on s.label = c.thread_owner
               join comments_nodes t on t.label = c.owner) as u;

-- name: FillWeightedAllEdges :exec
insert into all_edges_weighted (source, target, start, weight)
Select source, target, min(start) as start, count(*) as weight
from all_edges
group by source, target;

-- name: FillWeightedLikesEdges :exec
insert into likes_edges_weighted (source, target, start, weight)
Select source, target, min(start) as start, count(*) as weight
from likes_edges
group by source, target;

-- name: FillWeightedCommentsEdges :exec
insert into comments_edges_weighted (source, target, start, weight)
Select source, target, min(start) as start, count(*) as weight
from comments_edges
group by source, target;


-- name: CreateLikeNode :one
INSERT INTO likes_nodes (label, start)
VALUES ($1, $2)
RETURNING *;
--
-- name: CreateCommentNode :one
INSERT INTO comments_nodes (label, start)
VALUES ($1, $2)
RETURNING *;
--
-- name: CreateAllNode :one
INSERT INTO all_nodes (label, start)
VALUES ($1, $2)
RETURNING *;
--
-- name: CreateLikeEdge :one
INSERT INTO likes_edges (source, target, start)
VALUES ($1, $2, $3)
RETURNING *;
--
-- name: CreateCommentEdge :one
INSERT INTO comments_edges (source, target, start)
VALUES ($1, $2, $3)
RETURNING *;
--
-- name: CreateAllEdge :one
INSERT INTO all_edges (source, target, start)
VALUES ($1, $2, $3)
RETURNING *;