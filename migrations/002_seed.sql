INSERT INTO teams (team_name) VALUES
    ('backend'),
    ('frontend'),
    ('mobile')
ON CONFLICT DO NOTHING;

INSERT INTO users (user_id, username, team_name, is_active) VALUES
    ('u1', 'Alice', 'backend', true),
    ('u2', 'Bob', 'backend', true),
    ('u3', 'Charlie', 'backend', false),

    ('u4', 'Diana', 'frontend', true),
    ('u5', 'Ivan', 'frontend', true),

    ('u6', 'Michael', 'mobile', true),
    ('u7', 'Lev', 'mobile', false)
ON CONFLICT DO NOTHING;