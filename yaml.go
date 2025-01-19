package main

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

type ClusterGroupEntry struct {
	Name string
	Values []string
}

type GroupEntry struct {
	Name string
	Clusters []string
	//RenameTo string `yaml:"rename_to,omitempty"`
}

type ClusterGroup []ClusterGroupEntry
type Group []GroupEntry 

type Config struct {
	ClusterGroups ClusterGroup `yaml:"cluster_groups"`
	Groups Group `yaml:"groups"`
}

func (cg *ClusterGroup) UnmarshalYAML(value *yaml.Node) error {
	*cg = make([]ClusterGroupEntry, len(value.Content)/2)
	for i := 0; i < len(value.Content); i += 2 {
		var res = &(*cg)[i/2]
		if err := value.Content[i].Decode(&res.Name); err != nil {
			return err
		}
		if err := value.Content[i+1].Decode(&res.Values); err != nil {
			return err
		}
	}
	return nil
}

func (g *Group) UnmarshalYAML(value *yaml.Node) error {
    *g = make([]GroupEntry, len(value.Content)/2)
    for i := 0; i < len(value.Content); i += 2 {
        var res = &(*g)[i/2]
        if err := value.Content[i].Decode(&res.Name); err != nil {
            return err
        }
	content := value.Content[i+1]
	for x := 0; x < len(content.Content); x += 2 {
		res.Clusters = append(res.Clusters, content.Content[x+1].Content[0].Value)
	}
    }
    return nil
}

var input []byte = []byte(`
cluster_groups:
  build-farm:
  - app.ci
  - build01
  - build02
  - build03
  - build04
  - build05
  - build06
  - build09
  - build10
  - vsphere02
  dp-managed:
  - build01
  - build02
groups:
  art-admins:
    clusters:
    - app.ci
  build-api-team:
    clusters:
    - app.ci
  containers:
    cluster_groups:
    - dp-managed
  hypershift-team:
    clusters:
    - hosted-mgmt
    - app.ci
    rename_to: hypershift-pool-admins
  microshift-admins:
    clusters:
    - app.ci
`)

func main() {
	var f Config
	var err error
	if err = yaml.Unmarshal(input, &f); err != nil {
	    panic(err)
	}
	fmt.Printf("%v", f)
}
