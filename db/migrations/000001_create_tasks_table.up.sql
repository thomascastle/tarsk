CREATE TYPE priority AS ENUM ('none', 'low', 'medium', 'high');

CREATE TABLE tasks (
    description VARCHAR NOT NULL,
    done BOOLEAN DEFAULT FALSE NOT NULL,
    due_at TIMESTAMP WITHOUT TIME ZONE,
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    priority priority DEFAULT 'none' NOT NULL,
    started_at TIMESTAMP WITHOUT TIME ZONE
);
