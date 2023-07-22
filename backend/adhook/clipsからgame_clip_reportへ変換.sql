INSERT INTO game_clip_reports
SELECT
    c.game_id, c.broadcaster_id, c.broadcast_name, s.profile_image_url, c.video_id, c.id,
    thumbnail_url, c.title, c.view_count, c.vod_offset, c.duration, '2023-03-31', '30', -100 + RANK() OVER( ORDER BY c.view_count DESC, c.url ),
        c.url, c.created_at, c.updated_at
FROM clips c INNER JOIN streamers s on c.broadcaster_id = s.id
WHERE broadcaster_id = '144295271'
  AND view_count > 100;

-- 日付は柔軟に