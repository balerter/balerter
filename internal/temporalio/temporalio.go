package temporalio

import (
	"context"
	"github.com/balerter/balerter/internal/config/temporalio"
	"net/http"
	"sync"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

type runner interface {
	RunScript(name string, req *http.Request) error
}

type Temporalio struct {
	cfg    *temporalio.Temporalio
	rnr    runner
	logger *zap.Logger
}

func New(cfg *temporalio.Temporalio, rnr runner, logger *zap.Logger) *Temporalio {
	t := &Temporalio{
		cfg:    cfg,
		rnr:    rnr,
		logger: logger,
	}

	return t
}

func (t *Temporalio) Run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	var clients []client.Client
	var workers []worker.Worker

	for _, w := range t.cfg.Worker {
		clientOpts := client.Options{
			HostPort:  w.Host,
			Namespace: w.Namespace,
		}

		c, errCreateClient := client.NewClient(clientOpts)
		if errCreateClient != nil {
			t.logger.Error("error create temporalio client", zap.Error(errCreateClient))
			cancel()
			return
		}
		clients = append(clients, c)

		workerOpts := worker.Options{}

		wrk := worker.New(c, w.TaskQueueName, workerOpts)

		wfo := workflow.RegisterOptions{
			Name: w.WorkflowAlias,
		}

		wrk.RegisterWorkflowWithOptions(t.workflow, wfo)
		wrk.RegisterActivity(t.activity)
		errRun := wrk.Run(worker.InterruptCh())
		if errRun != nil {
			t.logger.Error("error run worker", zap.Error(errRun))
			cancel()
			return
		}
		workers = append(workers, wrk)
	}

	<-ctx.Done()

	for _, c := range clients {
		c.Close()
	}

}

func (t *Temporalio) workflow(ctx workflow.Context, payload []string) error {
	t.logger.Debug("temporalio workflow run", zap.Any("payload", payload))

	if len(payload) == 0 {
		return nil
	}

	for _, name := range payload {
		var resErr error

		ao := workflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
		}
		ctx1 := workflow.WithActivityOptions(ctx, ao)

		errRunActivity := workflow.ExecuteActivity(ctx1, t.activity, name).Get(ctx, &resErr)
		if errRunActivity != nil {
			t.logger.Error("error run activity", zap.Error(errRunActivity))
			return errRunActivity
		}

		if resErr != nil {
			t.logger.Error("activity returns error", zap.Error(resErr))
			return resErr
		}
	}

	return nil
}

func (t *Temporalio) activity(_ context.Context, name string) error {
	t.logger.Debug("temporalio activity run", zap.String("name", name))
	return t.rnr.RunScript(name, nil)
}
