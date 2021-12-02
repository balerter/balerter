package temporalio

type Worker struct {
	Name string `json:"name" yaml:"name" hcl:"name,label"`

	Host          string `json:"host" yaml:"host" hcl:"host"`
	Namespace     string `json:"namespace" yaml:"namespace" hcl:"namespace"`
	WorkflowAlias string `json:"workflowAlias" yaml:"workflowAlias" hcl:"workflowAlias"`
	TaskQueueName string `json:"taskQueueName" yaml:"taskQueueName" hcl:"taskQueueName"`
}

func (w *Worker) Validate() error {
	return nil
}

type Temporalio struct {
	Worker []Worker `json:"worker" yaml:"worker" hcl:"worker,block"`
}

func (c *Temporalio) Validate() error {
	return nil
}
