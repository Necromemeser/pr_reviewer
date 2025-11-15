INSERT INTO teams (team_name) VALUES
    ('backend'),
    ('frontend'),
    ('mobile')
ON CONFLICT DO NOTHING;

INSERT INTO users (user_id, username, team_name, is_active) VALUES
    ('u1', 'Ivan', 'backend', true),
    ('u2', 'Olga', 'backend', true),
    ('u3', 'Sergey', 'backend', true),
    ('u4', 'Anna', 'backend', true),
    ('u5', 'Mikhail', 'backend', false),

    ('u6', 'Daria', 'frontend', true),
    ('u7', 'Nikolay', 'frontend', true),
    ('u8', 'Ekaterina', 'frontend', true),
    ('u9', 'Alexey', 'frontend', true),
    ('u10', 'Svetlana', 'frontend', false),

    ('u11', 'Lev', 'mobile', true),
    ('u12', 'Maria', 'mobile', true),
    ('u13', 'Pavel', 'mobile', true),
    ('u14', 'Yulia', 'mobile', true),
    ('u15', 'Dmitry', 'mobile', false)
ON CONFLICT DO NOTHING;
