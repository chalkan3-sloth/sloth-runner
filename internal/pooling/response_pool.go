package pooling

import (
	"sync"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
)

// ResponsePool provides object pooling for gRPC responses to reduce allocations
type ResponsePool struct {
	resourceUsagePool sync.Pool
	processListPool   sync.Pool
	metricsDataPool   sync.Pool
}

// Global response pool instance
var GlobalResponsePool = NewResponsePool()

// NewResponsePool creates a new response pool
func NewResponsePool() *ResponsePool {
	return &ResponsePool{
		resourceUsagePool: sync.Pool{
			New: func() interface{} {
				return &pb.ResourceUsageResponse{}
			},
		},
		processListPool: sync.Pool{
			New: func() interface{} {
				return &pb.ProcessListResponse{
					Processes: make([]*pb.ProcessInfo, 0, 50),
				}
			},
		},
		metricsDataPool: sync.Pool{
			New: func() interface{} {
				return &pb.MetricsData{}
			},
		},
	}
}

// GetResourceUsageResponse gets a response from pool
func (p *ResponsePool) GetResourceUsageResponse() *pb.ResourceUsageResponse {
	resp := p.resourceUsagePool.Get().(*pb.ResourceUsageResponse)
	// Reset fields
	resp.Reset()
	return resp
}

// PutResourceUsageResponse returns response to pool
func (p *ResponsePool) PutResourceUsageResponse(resp *pb.ResourceUsageResponse) {
	p.resourceUsagePool.Put(resp)
}

// GetProcessListResponse gets a response from pool
func (p *ResponsePool) GetProcessListResponse() *pb.ProcessListResponse {
	resp := p.processListPool.Get().(*pb.ProcessListResponse)
	resp.Processes = resp.Processes[:0] // Clear slice but keep capacity
	return resp
}

// PutProcessListResponse returns response to pool
func (p *ResponsePool) PutProcessListResponse(resp *pb.ProcessListResponse) {
	p.processListPool.Put(resp)
}

// GetMetricsData gets a metrics data from pool
func (p *ResponsePool) GetMetricsData() *pb.MetricsData {
	data := p.metricsDataPool.Get().(*pb.MetricsData)
	data.Reset()
	return data
}

// PutMetricsData returns metrics data to pool
func (p *ResponsePool) PutMetricsData(data *pb.MetricsData) {
	p.metricsDataPool.Put(data)
}
