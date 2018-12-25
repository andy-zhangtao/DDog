package agent

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/8.
import zmodel "github.com/openzipkin/zipkin-go/model"

type DeployMsg struct {
	//	Span 链路跟踪跨度数据
	Span      zmodel.SpanContext `json:"span"`
	SvcName   string             `json:"svc_name"`
	NameSpace string             `json:"name_space"`
	Upgrade   bool               `json:"upgrade"`
	Replicas  int                `json:"replicas"`
	DC        string             `json:"dc"`
}
