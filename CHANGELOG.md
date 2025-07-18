<!--
All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
-->

## Version 1.3.0

### Features
- Add a new `force_reparse_old_blocks` option to the indexer config 
to force the indexer to reparse the blocks from the start height to the current node height

### Bug Fixes
- Fix error when unmarshalling the `Base64Bytes` type
- Fix missing initialization of the `globalObjects` field inside `IndexersBuilder`


## Version 1.2.1

### Features
- Add a hook to execute custom logic before the CLI start command executes
- Add a hook to execute custom logic after the configuration file have been read
- Allow objects to be shared across all modules

### Bug Fixes
- Fix missing hash in cosmos transactions

## Version 1.1.0

### Features
- Create a gRPC over RPC connection to communicate with a cosmos node
- Improve worker error messages

## Version 1.0.1

### Bug Fixes
- Fix typo in `modules/adapter` module

## Version 1.0.0

This is the first release of the project.

### Features
- Support for indexing Cosmos-SDK based chains
- Support for PostgreSQL as a database backend

