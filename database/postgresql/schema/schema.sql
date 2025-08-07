CREATE TABLE blocks
(
    -- Name of the indexer that has indexed the block.
    indexer     TEXT NOT NULL,
    -- ID of the chain from which the block has been fetched.
    chain_id    TEXT NOT NULL,
    -- Height of the indexed block.
    height      BIGINT,
    -- Time at which the indexed block has been produced by the chain.
    timestamp   TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT unique_chain_block UNIQUE (indexer, chain_id, height)
);

