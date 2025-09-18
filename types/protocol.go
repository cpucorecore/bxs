package types

const (
	ProtocolIdNewSwap = iota + 1
	ProtocolIdUniswapV3
	ProtocolIdXLaunch
)

const (
	ProtocolNameNewSwap   = "NewSwap"
	ProtocolNameUniswapV3 = "UniswapV3"
	ProtocolNameXLaunch   = "XLaunch"
)

func GetProtocolName(protocolId int) string {
	switch protocolId {
	case ProtocolIdNewSwap:
		return ProtocolNameNewSwap
	case ProtocolIdUniswapV3:
		return ProtocolNameUniswapV3
	case ProtocolIdXLaunch:
		return ProtocolNameXLaunch
	default:
		return "Unknown"
	}
}
