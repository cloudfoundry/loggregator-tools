package logcacheutil

import (
	"fmt"

	"code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
)

func DrainTypeToEnvelopeTypes(drainType string) ([]logcache_v1.EnvelopeType, error) {
	switch drainType {
	case "logs":
		return []logcache_v1.EnvelopeType{logcache_v1.EnvelopeType_LOG}, nil
	case "metrics":
		return []logcache_v1.EnvelopeType{
			logcache_v1.EnvelopeType_COUNTER,
			logcache_v1.EnvelopeType_GAUGE,
		}, nil
	case "all":
		return []logcache_v1.EnvelopeType{
			logcache_v1.EnvelopeType_LOG,
			logcache_v1.EnvelopeType_COUNTER,
			logcache_v1.EnvelopeType_GAUGE,
		}, nil
	default:
		return nil, fmt.Errorf("unknown drain type: %s", drainType)
	}
}
