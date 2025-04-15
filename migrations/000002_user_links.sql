-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_links (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    short_link_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_user_links_user
        FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE CASCADE,
        
    CONSTRAINT fk_user_links_short_link
        FOREIGN KEY (short_link_id) 
        REFERENCES short_links(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_user_links_user_id ON user_links(user_id);
CREATE INDEX idx_user_links_short_link_id ON user_links(short_link_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_links;
-- +goose StatementEnd


