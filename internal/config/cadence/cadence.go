package cadence

type Worker struct {
	Name string `json:"name" yaml:"name" hcl:"name,label"`

	Host           string `json:"host" yaml:"host" hcl:"host"`
	Domain         string `json:"domain" yaml:"domain" hcl:"domain"`
	TaskListName   string `json:"taskListName" yaml:"taskListName" hcl:"taskListName"`
	ClientName     string `json:"clientName" yaml:"clientName" hcl:"clientName"`
	CadenceService string `json:"сadenceService" yaml:"сadenceService" hcl:"cadenceService"`
	WorkflowAlias  string `json:"workflowAlias" yaml:"workflowAlias" hcl:"workflowAlias"`
}

func (w *Worker) Validate() error {
	return nil
}

type Cadence struct {
	Worker []Worker `json:"worker" yaml:"worker" hcl:"worker,block"`
}

func (c *Cadence) Validate() error {
	return nil
}
