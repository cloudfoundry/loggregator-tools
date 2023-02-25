package traverser

import (
	"hash/crc64"

	"code.cloudfoundry.org/go-pubsub"
)

func TraverserTraverse(data interface{}) pubsub.Paths {
	return _Added(data)
}

func done(data interface{}) pubsub.Paths {
	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		return 0, nil, false
	})
}

func hashBool(data bool) uint64 {
	if data {
		return 1
	}
	return 0
}

var tableECMA = crc64.MakeTable(crc64.ECMA)

func _Added(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(_Key), true
		case 1:

			return hashBool(data.(Notification).Added), pubsub.TreeTraverser(_Key), true
		default:
			return 0, nil, false
		}
	})
}

func _Key(data interface{}) pubsub.Paths {

	return pubsub.Paths(func(idx int, data interface{}) (path uint64, nextTraverser pubsub.TreeTraverser, ok bool) {
		switch idx {
		case 0:
			return 0, pubsub.TreeTraverser(done), true
		case 1:

			return crc64.Checksum([]byte(data.(Notification).Key), tableECMA), pubsub.TreeTraverser(done), true
		default:
			return 0, nil, false
		}
	})
}

type NotificationFilter struct {
	Added *bool
	Key   *string
}

func TraverserCreatePath(f *NotificationFilter) []uint64 {
	if f == nil {
		return nil
	}
	var path []uint64

	var count int
	if count > 1 {
		panic("Only one field can be set")
	}

	if f.Added != nil {

		path = append(path, hashBool(*f.Added))
	} else {
		path = append(path, 0)
	}

	if f.Key != nil {

		path = append(path, crc64.Checksum([]byte(*f.Key), tableECMA))
	} else {
		path = append(path, 0)
	}

	for i := len(path) - 1; i >= 1; i-- {
		if path[i] != 0 {
			break
		}
		path = path[:i]
	}

	return path
}
