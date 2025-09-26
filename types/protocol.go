package types

const (
	ProtocolIdXLaunch = iota + 1
	ProtocolIdPancakeV2
)

const (
	ProtocolNameXLaunch   = "XLaunch"
	ProtocolNamePancakeV2 = "PancakeV2"
)

func GetProtocolName(id int) string {
	switch id {
	case ProtocolIdXLaunch:
		return ProtocolNameXLaunch
	case ProtocolIdPancakeV2:
		return ProtocolNamePancakeV2
	default:
		panic("invalid id")
	}
}
