CREATE TABLE blocks
(
    chain_id    TEXT NOT NULL,
    height      BIGINT,
    timestamp   TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT unique_chain_block UNIQUE (chain_id, height)
);

