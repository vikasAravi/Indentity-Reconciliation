CREATE TABLE contact (
                         id SERIAL PRIMARY KEY,
                         phone_number TEXT,
                         email TEXT,
                         linked_id INTEGER REFERENCES contact(id),
                         link_precedence TEXT,
                         created_at TIMESTAMPTZ DEFAULT current_timestamp,
                         updated_at TIMESTAMPTZ DEFAULT current_timestamp,
                         deleted_at TIMESTAMPTZ DEFAULT NULL
);
