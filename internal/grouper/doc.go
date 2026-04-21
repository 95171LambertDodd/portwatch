// Package grouper provides a Grouper type that partitions a slice of
// portscanner.PortEntry values into named buckets based on a configurable
// field.
//
// Supported grouping strategies:
//
//	"protocol"  — groups entries by their network protocol (tcp / udp).
//	"process"   — groups entries by the owning process name; entries with no
//	              resolved process are placed in the "unknown" bucket.
//	"portband"  — groups entries by IANA port-range band:
//	                well-known  (0–1023)
//	                registered  (1024–49151)
//	                dynamic     (49152–65535)
//
// Example:
//
//	g, err := grouper.New(grouper.GroupByProtocol)
//	if err != nil {
//		log.Fatal(err)
//	}
//	groups := g.Group(entries)
//	for _, grp := range groups {
//		fmt.Printf("%s: %d entries\n", grp.Key, len(grp.Entries))
//	}
package grouper
