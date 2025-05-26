package types

import (
	"context"

	"github.com/rs/zerolog"
)

type IndexerContextKey string

const ContextKey = IndexerContextKey("indexer.context")

type IndexerContext struct {
	Config        *Config
	IndexerConfig *IndexerConfig
	Logger        zerolog.Logger
}

func NewIndexerContext(cfg *Config, indexerCfg *IndexerConfig, logger zerolog.Logger) IndexerContext {
	return IndexerContext{
		Config:        cfg,
		IndexerConfig: indexerCfg,
		Logger:        logger.With().Str("indexer", indexerCfg.Name).Logger(),
	}
}

func InjectIndexerContext(ctx context.Context, indexerCtx IndexerContext) context.Context {
	return context.WithValue(ctx, ContextKey, indexerCtx)
}

func GetIndexerContext(ctx context.Context) IndexerContext {
	indexerCtx, ok := ctx.Value(ContextKey).(IndexerContext)
	if !ok {
		panic("can't get IndexerContext from the provided Context")
	}

	return indexerCtx
}
