CREATE TABLE IF NOT EXISTS registros (
                                         id SERIAL PRIMARY KEY,
                                         trace_id UUID NOT NULL,
                                         payload JSONB NOT NULL,
                                         byte_size INTEGER NOT NULL,
                                         total_characters INTEGER NOT NULL,
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE registros
    ADD COLUMN published_to_queue BOOLEAN DEFAULT FALSE;