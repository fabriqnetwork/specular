package eth

type BlockTag string

// Block tag values (used both for L1 and L2)
// https://ethereum.github.io/execution-apis/api-documentation/
const (
	// - L1: The most recent block in the canonical chain observed by the client,
	// 	     this block may be re-orged out of the canonical chain even under healthy/normal conditions
	// - L2: Local chain head (unconfirmed on L1)
	Latest = "latest"
	// - L1: The most recent block that is safe from re-orgs under honest majority and certain synchronicity assumptions
	// - L2: Derived chain tip from L1 data
	Safe = "safe"
	// - L1: The most recent crypto-economically secure block,
	//       cannot be re-orged outside of manual intervention driven by community coordination
	// - L2: Derived chain tip from finalized L1 data
	Finalized = "finalized"
)
