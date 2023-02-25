package logcache

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/golang/protobuf/jsonpb"
)

// GroupReaderClient reads and interacts from LogCache via the RESTful or gRPC
// Group API.
type GroupReaderClient struct {
	addr string

	httpClient HTTPClient
	grpcClient logcache_v1.GroupReaderClient
}

// NewGroupReaderClient creates a GroupReaderClient.
func NewGroupReaderClient(addr string, opts ...ClientOption) *GroupReaderClient {
	c := &GroupReaderClient{
		addr: addr,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	for _, o := range opts {
		o.configure(c)
	}

	return c
}

// BuildReader is used to create a Reader (useful for things like Walk) with a
// RequesterID. It simply wraps the GroupReaderClient.Read method.
func (c *GroupReaderClient) BuildReader(requesterID uint64) Reader {
	return Reader(func(
		ctx context.Context,
		name string,
		start time.Time,
		opts ...ReadOption,
	) ([]*loggregator_v2.Envelope, error) {
		return c.Read(ctx, name, start, requesterID, opts...)
	})
}

// Read queries the LogCache and returns the given envelopes. To override any
// query defaults (e.g., end time), use the according option.
func (c *GroupReaderClient) Read(
	ctx context.Context,
	name string,
	start time.Time,
	requesterID uint64,
	opts ...ReadOption,
) ([]*loggregator_v2.Envelope, error) {
	if c.grpcClient != nil {
		return c.grpcRead(ctx, name, start, requesterID, opts)
	}

	u, err := url.Parse(c.addr)
	if err != nil {
		return nil, err
	}
	u.Path = "v1/group/" + name
	q := u.Query()
	q.Set("start_time", strconv.FormatInt(start.UnixNano(), 10))
	q.Set("requester_id", strconv.FormatUint(requesterID, 10))

	// allow the given options to configure the URL.
	for _, o := range opts {
		o(u, q)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var r logcache_v1.ReadResponse
	if err := jsonpb.Unmarshal(resp.Body, &r); err != nil {
		return nil, err
	}

	return r.Envelopes.Batch, nil
}

func (c *GroupReaderClient) grpcRead(
	ctx context.Context,
	name string,
	start time.Time,
	requesterID uint64,
	opts []ReadOption,
) ([]*loggregator_v2.Envelope, error) {
	u := &url.URL{}
	q := u.Query()
	// allow the given options to configure the URL.
	for _, o := range opts {
		o(u, q)
	}

	req := &logcache_v1.GroupReadRequest{
		Name:        name,
		RequesterId: requesterID,
		StartTime:   start.UnixNano(),
	}

	if v, ok := q["limit"]; ok {
		req.Limit, _ = strconv.ParseInt(v[0], 10, 64)
	}

	if v, ok := q["end_time"]; ok {
		req.EndTime, _ = strconv.ParseInt(v[0], 10, 64)
	}

	if v, ok := q["envelope_types"]; ok {
		req.EnvelopeTypes = []logcache_v1.EnvelopeType{
			logcache_v1.EnvelopeType(logcache_v1.EnvelopeType_value[v[0]]),
		}
	}

	resp, err := c.grpcClient.Read(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Envelopes.Batch, nil
}

// AddToGroup adds a sourceID to the given group. If the group doesn't exist,
// then it is created. If the group already has the given sourceID, then it is
// a NOP.
func (c *GroupReaderClient) AddToGroup(ctx context.Context, name, sourceID string) error {
	if c.grpcClient != nil {
		return c.grpcAddToGroup(ctx, name, sourceID)
	}

	u, err := url.Parse(c.addr)
	if err != nil {
		return err
	}
	u.Path = fmt.Sprintf("v1/group/%s/%s", name, sourceID)

	req, err := http.NewRequest("PUT", u.String(), nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}

func (c *GroupReaderClient) grpcAddToGroup(ctx context.Context, name, sourceID string) error {
	_, err := c.grpcClient.AddToGroup(ctx, &logcache_v1.AddToGroupRequest{
		Name:     name,
		SourceId: sourceID,
	})
	return err
}

// RemoveFromGroup removes a sourceID from the given group. If the given
// sourceID was the last one, then the grou is removed. If the group does not
// have the given sourceID, then it is a NOP.
func (c *GroupReaderClient) RemoveFromGroup(ctx context.Context, name, sourceID string) error {
	if c.grpcClient != nil {
		return c.grpcRemoveFromGroup(ctx, name, sourceID)
	}

	u, err := url.Parse(c.addr)
	if err != nil {
		return err
	}
	u.Path = fmt.Sprintf("v1/group/%s/%s", name, sourceID)

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}

func (c *GroupReaderClient) grpcRemoveFromGroup(ctx context.Context, name, sourceID string) error {
	_, err := c.grpcClient.RemoveFromGroup(ctx, &logcache_v1.RemoveFromGroupRequest{
		Name:     name,
		SourceId: sourceID,
	})
	return err
}

// GroupMeta gives the information about given group.
type GroupMeta struct {
	SourceIDs    []string
	RequesterIDs []uint64
}

// Group returns the meta information about a group.
func (c *GroupReaderClient) Group(ctx context.Context, name string) (GroupMeta, error) {
	if c.grpcClient != nil {
		return c.grpcGroup(ctx, name)
	}

	u, err := url.Parse(c.addr)
	if err != nil {
		return GroupMeta{}, err
	}
	u.Path = fmt.Sprintf("v1/group/%s/meta", name)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return GroupMeta{}, err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GroupMeta{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GroupMeta{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GroupMeta{}, err
	}

	gresp := logcache_v1.GroupResponse{}
	if err := json.Unmarshal(data, &gresp); err != nil {
		return GroupMeta{}, err
	}

	return GroupMeta{
		SourceIDs:    gresp.SourceIds,
		RequesterIDs: gresp.RequesterIds,
	}, nil
}

func (c *GroupReaderClient) grpcGroup(ctx context.Context, name string) (GroupMeta, error) {
	resp, err := c.grpcClient.Group(ctx, &logcache_v1.GroupRequest{
		Name: name,
	})

	if err != nil {
		return GroupMeta{}, err
	}

	return GroupMeta{
		SourceIDs:    resp.SourceIds,
		RequesterIDs: resp.RequesterIds,
	}, nil
}
