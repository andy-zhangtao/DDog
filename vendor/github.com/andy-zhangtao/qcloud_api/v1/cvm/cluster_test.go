package cvm

import (
	"testing"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
)

func TestQueryCluster(t *testing.T) {
	q := Cluster{
		Cid:        "cls-rfje0azd",
		Cname:      "caas-test",
		Status:     "Running",
		OrderField: "createdAt",
		OrderType:  "desc",
		Offset:     0,
		Limit:      20,
	}
	field, req := q.queryCluster()
	if len(field) != 6 {
		t.Error("Field Size Error!")
	}

	if len(req) != 6 {
		t.Error("Request Error!")
	}

	q.Offset = 10
	field, req = q.queryCluster()
	if len(field) != 7 {
		t.Error("Field Size Error!")
	}

	if len(req) != 7 {
		t.Error("Request Error!")
	}
}

func TestGenerateSignature(t *testing.T) {
	q := Cluster{
		Cid:        "cls-rfje0azd",
		Cname:      "caas-test",
		Status:     "Running",
		OrderField: "createdAt",
		OrderType:  "desc",
		Offset:     0,
		Limit:      20,
	}
	field, reqMap := q.queryCluster()
	publicMap := public.PublicParam("DescribeInstances", "sh", "123456")
	public.GenerateSignatureString(field, reqMap, publicMap)
	//if req != "Action=DescribeInstances&Nonce=827870&Region=sh&SecretId=123456&Timestamp=1514276701&clusterIds.n=cls-rfje0azd&clusterName=caas-test&limit=20&orderField=createdAt&orderType=desc&status=Running" {
	//	t.Error("Req Generate Error!")
	//}
}

func TestCluster_QueryClusters(t *testing.T) {
	q := Cluster{
		Pub: public.Public{
			Region:   "sh",
			SecretId: "123",
		},
		Cid:        "cls-rfje0azd",
		Cname:      "caas-test",
		Status:     "Running",
		OrderField: "createdAt",
		OrderType:  "desc",
		Offset:     0,
		Limit:      20,
		SecretKey:  "321",
	}

	q.QueryClusters()
}

func TestCluster_QueryClusterNodes(t *testing.T) {
	q := Cluster{
		Pub: public.Public{
			Region:   "sh",
			SecretId: "123",
		},
		Cid:       "cls-rfje0azd",
		Namespace: "devenv",
		Offset:    0,
		Limit:     20,
		SecretKey: "321",
	}

	q.QueryClusterNodes()
}
