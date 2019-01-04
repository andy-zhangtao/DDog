package k8service

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/k8s/k8smodel"
	"github.com/andy-zhangtao/DDog/model/k8sconfig"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
	"time"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/14.
//获取K8s DeployMent数据

const (
	ModuleName = "K8s-maintainer-agent"
)

func AddNewK8sClusterConfigeDate(kc k8sconfig.K8sCluster) (err error) {
	if kc.Region == "" || kc.Endpoint == "" || kc.Token == "" {
		return errors.New(fmt.Sprintf("Regin/Endpoint/Token Empty! [%v]", kc))
	}

	if err = mongo.SaveK8sClusterData(kc); err != nil {
		err = errors.New(fmt.Sprintf("Add New K8s Configure Error [%s] K8s [%v]", err.Error(), kc))
	}

	return
}

func GetALlK8sCluster() (kc []k8sconfig.K8sCluster, err error) {
	if kc, err = mongo.GetAllK8sClusterData(); err != nil {
		err = errors.New(fmt.Sprintf("Query All K8s Data Error [%s]", err.Error()))
	}

	return
}

func GetK8sClusterWithNamespace(region, namespace string) (kc k8sconfig.K8sCluster, err error) {
	if region == "" {
		region = os.Getenv(_const.EnvRegion)
	}

	if region == "" {
		err = errors.New(fmt.Sprintf("Region Empty!"))
		return
	}

	kcs, err := mongo.GetAllK8sClusterData()
	if err != nil {
		err = errors.New(fmt.Sprintf("Query K8s Cluster Data Error [%s] Region [%s]", err.Error(), region))
		return
	}

	for _, k := range kcs {
		if strings.Compare(k.Namespace, namespace) == 0 {
			return k, nil
		}
	}

	return kc, errors.New("No K8s Cluster")
}

func GetK8sCluster(region string) (kc k8sconfig.K8sCluster, err error) {
	if region == "" {
		region = os.Getenv(_const.EnvRegion)
	}

	if region == "" {
		err = errors.New(fmt.Sprintf("Region Empty!"))
		return
	}

	if kc, err = mongo.GetK8sClusterData(region); err != nil {
		err = errors.New(fmt.Sprintf("Query K8s Cluster Data Error [%s] Region [%s]", err.Error(), region))
	}

	return
}

func UpdateK8sCluster(kc k8sconfig.K8sCluster) (err error) {
	if kc.Region == "" {
		kc.Region = os.Getenv(_const.EnvRegion)
	}

	if kc.Region == "" {
		err = errors.New(fmt.Sprintf("Region Empty!"))
		return
	}

	if err = mongo.UpdateK8sClusterData(kc); err != nil {
		errors.New(fmt.Sprintf("Update K8s Cluster Data Error [%s] K8s [%s]", err.Error(), kc))
	}

	return
}

func DeleteK8sCluster(region string) (err error) {
	if region == "" {
		region = os.Getenv(_const.EnvRegion)
	}

	if region == "" {
		err = errors.New(fmt.Sprintf("Region Empty!"))
		return
	}

	if kc, err := mongo.GetK8sClusterData(region); err != nil {
		return errors.New(fmt.Sprintf("Query K8s Cluster Data Error [%s] Region [%s]", err.Error(), region))
	} else {
		if err = mongo.DeleteK8sClusterDataByID(kc.ID); err != nil {
			return errors.New(fmt.Sprintf("Delete K8s Cluster Data Error [%s] ID [%s]", err.Error(), kc.ID))
		}
	}

	return
}

func DeleteK8sClusterByID(id string) (err error) {
	if err = mongo.DeleteK8sClusterDataByID(bson.ObjectIdHex(id)); err != nil {
		return errors.New(fmt.Sprintf("Delete K8s Cluster Data Error [%s] ID [%s]", err.Error(), id))
	}

	return
}

//BackupK8sCluster 备份指定的K8s控制数据
//region K8s所在机房
//name K8s控制命名空间名称
func BackupK8sCluster(region, name string) (fileName string, err error) {

	fileContent := make(map[string][]byte)
	serviceMap := make(map[string]k8smodel.K8sServiceBackup)
	services, err := GetK8sSpecService(region, name, "")
	if err != nil {
		err = errors.New(fmt.Sprintf("Query [%s] Service Error [%s]", name, err.Error()))
		return
	}

	for _, s := range services.Items {
		serviceMap[s.Metadata.Name] = k8smodel.K8sServiceBackup{
			Apiversion: services.ApiVersion,
			Kind:       "Service",
			Metadata: k8smodel.K8sServiceBackup_Metadata{
				Annotations:       s.Metadata.Annotations,
				CreationTimestamp: s.Metadata.CreationTimestamp,
				Labels:            s.Metadata.Labels,
				Namespace:         s.Metadata.Namespace,
				Name:              s.Metadata.Name,
				ResourceVersion:   s.Metadata.ResourceVersion,
				SelfLink:          s.Metadata.SelfLink,
				Uid:               s.Metadata.Uid,
			},
			Spec:   s.Spec,
			Status: s.Status,
		}
	}

	logrus.WithFields(logrus.Fields{"Cache Services": len(serviceMap)}).Info(ModuleName)

	deployments, err := GetK8sDeployMent(region, name)
	logrus.WithFields(logrus.Fields{"Deployments": deployments}).Info(ModuleName)

	//var fileName []string

	for _, d := range deployments.Items {
		depBack := k8smodel.K8sDeployBackup{
			Apiversion: "extensions/v1beta1",
			Kind:       "Deployment",
		}

		depBack.Metadata = k8smodel.K8sDeployBackup_Metadata{
			Annotations:       d.Metadata.Annotations,
			CreationTimestamp: d.Metadata.CreationTimestamp,
			Generation:        d.Metadata.Generation,
			Labels:            d.Metadata.Labels,
			Name:              d.Metadata.Name,
			Namespace:         d.Metadata.Namespace,
			ResourceVersion:   d.Metadata.ResourceVersion,
			SelfLink:          strings.Replace(d.Metadata.SelfLink, "apps", "extensions", 1),
			Uid:               d.Metadata.Uid,
		}

		depBack.Spec = k8smodel.K8sDeployBackup_spec{
			MinReadySeconds:      d.Spec.MinReadySeconds,
			Replicas:             d.Spec.Replicas,
			RevisionHistoryLimit: d.Spec.RevisionHistoryLimit,
			Selector:             d.Spec.Selector,
			Strategy:             d.Spec.Strategy,
			Template:             d.Spec.Template,
		}

		depBack.Status = d.Status
		if name, content, err := outputYamlFile(depBack, serviceMap[d.Metadata.Name]); err != nil {
			err = errors.New(fmt.Sprintf("Generate Yaml File Error [%s]", err.Error()))
			return "", err
		} else {
			fileContent[name] = content
		}
	}

	fileName = fmt.Sprintf("k8s-%s-%s-backup-%s.zip", region, name, time.Now().Format("2006-01-02T15:04"))
	fzip, _ := os.Create("/tmp/" + fileName)
	w := zip.NewWriter(fzip)

	for name, content := range fileContent {
		f, err := w.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(content)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Make sure to check the error on Close.
	err = w.Close()
	return
}

func outputYamlFile(d k8smodel.K8sDeployBackup, s k8smodel.K8sServiceBackup) (filename string, content []byte, err error) {
	filename = fmt.Sprintf("/tmp/deploy_%s_back.yaml", d.Metadata.Name)
	//sn := fmt.Sprintf("/tmp/deploy_%s_back.yaml", d.Metadata.Name)

	content, err = yaml.Marshal(&d)
	if err != nil {
		return
	}

	if s.Apiversion != "" {
		sd, err := yaml.Marshal(&s)
		if err != nil {
			return filename, content, err
		}

		content = append(content, []byte("\n---\n")...)
		content = append(content, sd...)
	}

	//err = ioutil.WriteFile(filename, content, os.ModePerm)
	return
}
