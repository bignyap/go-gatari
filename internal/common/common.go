package common

type PubSubChannel string

const (
	EndpointCreated PubSubChannel = "endpoint:created"
	EndpointDeleted PubSubChannel = "endpoint:deleted"
)

type EndpointCreatedEvent struct {
	Code   string
	Path   string
	Method string
}

type EndpointDeletedEvent struct {
	Code string
}

func FetchAll[T any](fetchFunc func(offset, batchsize int32) ([]T, error), batchsize int32) ([]T, error) {

	var results []T
	offset := int32(0)

	for {
		items, err := fetchFunc(offset, batchsize)
		if err != nil {
			return nil, err
		}

		results = append(results, items...)

		if int32(len(items)) < batchsize {
			break
		}

		offset += batchsize
	}

	return results, nil
}
