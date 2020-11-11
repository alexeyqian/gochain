package config

// system
//const IsIntBigEndian = true

// application
const Storage = "BoltStorage" //"MemoryStorage"
//const Storage = "MemoryStorage"
const BoltDbFileName = "data/status.db"

const GenesisTime = 1605065492         // must be dividable by 3
const DefaultBlockSize = 1024 * 10     // 10K
const DefaultAccountCreationFee = 1000 // 1 coin

const MaxWitnesses = 21
const BlockInterval = 3
