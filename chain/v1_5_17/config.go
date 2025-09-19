package v1_5_17

import (
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

func newUint64(val uint64) *uint64 { return &val }

// BSCChainConfig defines the chain configuration for Binance Smart Chain (BSC)
// https://github.com/bnb-chain/bsc/blob/v1.5.17/params/config.go#L156
var BSCChainConfig = &params.ChainConfig{
	ChainID:             big.NewInt(56),
	HomesteadBlock:      big.NewInt(0),
	EIP150Block:         big.NewInt(0),
	EIP155Block:         big.NewInt(0),
	EIP158Block:         big.NewInt(0),
	ByzantiumBlock:      big.NewInt(0),
	ConstantinopleBlock: big.NewInt(0),
	PetersburgBlock:     big.NewInt(0),
	IstanbulBlock:       big.NewInt(0),
	MuirGlacierBlock:    big.NewInt(0),
	RamanujanBlock:      big.NewInt(0),
	NielsBlock:          big.NewInt(0),
	MirrorSyncBlock:     big.NewInt(5184000),
	BrunoBlock:          big.NewInt(13082000),
	EulerBlock:          big.NewInt(18907621),
	NanoBlock:           big.NewInt(21962149),
	MoranBlock:          big.NewInt(22107423),
	GibbsBlock:          big.NewInt(23846001),
	PlanckBlock:         big.NewInt(27281024),
	LubanBlock:          big.NewInt(29020050),
	PlatoBlock:          big.NewInt(30720096),
	BerlinBlock:         big.NewInt(31302048),
	LondonBlock:         big.NewInt(31302048),
	HertzBlock:          big.NewInt(31302048),
	HertzfixBlock:       big.NewInt(34140700),
	ShanghaiTime:        newUint64(1705996800), // 2024-01-23 08:00:00 AM UTC
	KeplerTime:          newUint64(1705996800), // 2024-01-23 08:00:00 AM UTC
	FeynmanTime:         newUint64(1713419340), // 2024-04-18 05:49:00 AM UTC
	FeynmanFixTime:      newUint64(1713419340), // 2024-04-18 05:49:00 AM UTC
	CancunTime:          newUint64(1718863500), // 2024-06-20 06:05:00 AM UTC
	HaberTime:           newUint64(1718863500), // 2024-06-20 06:05:00 AM UTC
	HaberFixTime:        newUint64(1727316120), // 2024-09-26 02:02:00 AM UTC
	BohrTime:            newUint64(1727317200), // 2024-09-26 02:20:00 AM UTC
	PascalTime:          newUint64(1742436600), // 2025-03-20 02:10:00 AM UTC
	PragueTime:          newUint64(1742436600), // 2025-03-20 02:10:00 AM UTC
	LorentzTime:         newUint64(1745903100), // 2025-04-29 05:05:00 AM UTC
	MaxwellTime:         newUint64(1751250600), // 2025-06-30 02:30:00 AM UTC
	FermiTime:           nil,

	Parlia: &params.ParliaConfig{},
	BlobScheduleConfig: &params.BlobScheduleConfig{
		Cancun: params.DefaultCancunBlobConfig,
		Prague: params.DefaultPragueBlobConfigBSC,
	},
}

var ChapelChainConfig = &params.ChainConfig{
	ChainID:             big.NewInt(97),
	HomesteadBlock:      big.NewInt(0),
	EIP150Block:         big.NewInt(0),
	EIP155Block:         big.NewInt(0),
	EIP158Block:         big.NewInt(0),
	ByzantiumBlock:      big.NewInt(0),
	ConstantinopleBlock: big.NewInt(0),
	PetersburgBlock:     big.NewInt(0),
	IstanbulBlock:       big.NewInt(0),
	MuirGlacierBlock:    big.NewInt(0),
	RamanujanBlock:      big.NewInt(1010000),
	NielsBlock:          big.NewInt(1014369),
	MirrorSyncBlock:     big.NewInt(5582500),
	BrunoBlock:          big.NewInt(13837000),
	EulerBlock:          big.NewInt(19203503),
	GibbsBlock:          big.NewInt(22800220),
	NanoBlock:           big.NewInt(23482428),
	MoranBlock:          big.NewInt(23603940),
	PlanckBlock:         big.NewInt(28196022),
	LubanBlock:          big.NewInt(29295050),
	PlatoBlock:          big.NewInt(29861024),
	BerlinBlock:         big.NewInt(31103030),
	LondonBlock:         big.NewInt(31103030),
	HertzBlock:          big.NewInt(31103030),
	HertzfixBlock:       big.NewInt(35682300),
	ShanghaiTime:        newUint64(1702972800), // 2023-12-19 8:00:00 AM UTC
	KeplerTime:          newUint64(1702972800),
	FeynmanTime:         newUint64(1710136800), // 2024-03-11 6:00:00 AM UTC
	FeynmanFixTime:      newUint64(1711342800), // 2024-03-25 5:00:00 AM UTC
	CancunTime:          newUint64(1713330442), // 2024-04-17 05:07:22 AM UTC
	HaberTime:           newUint64(1716962820), // 2024-05-29 06:07:00 AM UTC
	HaberFixTime:        newUint64(1719986788), // 2024-07-03 06:06:28 AM UTC
	BohrTime:            newUint64(1724116996), // 2024-08-20 01:23:16 AM UTC
	PascalTime:          newUint64(1740452880), // 2025-02-25 03:08:00 AM UTC
	PragueTime:          newUint64(1740452880), // 2025-02-25 03:08:00 AM UTC
	LorentzTime:         newUint64(1744097580), // 2025-04-08 07:33:00 AM UTC
	MaxwellTime:         newUint64(1748243100), // 2025-05-26 07:05:00 AM UTC
	FermiTime:           nil,

	Parlia: &params.ParliaConfig{},
	BlobScheduleConfig: &params.BlobScheduleConfig{
		Cancun: params.DefaultCancunBlobConfig,
		Prague: params.DefaultPragueBlobConfigBSC,
	},
}
