package cadence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/balerter/balerter/internal/config/cadence"

	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
)

type runner interface {
	RunScript(name string, req *http.Request) error
}

type Cadence struct {
	cfg    *cadence.Cadence
	rnr    runner
	logger *zap.Logger
}

func New(cfg *cadence.Cadence, rnr runner, logger *zap.Logger) *Cadence {
	cd := &Cadence{
		cfg:    cfg,
		rnr:    rnr,
		logger: logger,
	}

	return cd
}

func (cd *Cadence) Run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	var workers []worker.Worker

	for _, w := range cd.cfg.Worker {
		cd.logger.Debug("build cadence client", zap.Any("w", w))

		service, err := buildCadenceClient(w)
		if err != nil {
			cd.logger.Error("error cadence init", zap.Error(err))
			cancel()
			return
		}

		workerOptions := worker.Options{
			Logger: cd.logger,
		}

		wrk := worker.New(service, w.Domain, w.TaskListName, workerOptions)

		wrk.RegisterWorkflowWithOptions(cd.workflow, workflow.RegisterOptions{
			Name: w.WorkflowAlias,
		})

		wrk.RegisterActivity(cd.activityRunScript)

		errRun := wrk.Start()
		if errRun != nil {
			cd.logger.Error("failed to start workers", zap.Error(errRun))
			cancel()
			return
		}

		workers = append(workers, wrk)
	}

	<-ctx.Done()

	cd.logger.Debug("stopping cadence workers")

	for _, wrk := range workers {
		wrk.Stop()
	}
}

func (cd *Cadence) workflow(ctx workflow.Context, payload []byte) error {
	cd.logger.Debug("start workflow", zap.ByteString("payload", payload))

	var scripts []string
	errDecode := json.Unmarshal(payload, &scripts)
	if errDecode != nil {
		cd.logger.Error("error decode cadence payload", zap.Error(errDecode))
		return errDecode
	}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	for _, s := range scripts {
		var resErr error
		errActivity := workflow.ExecuteActivity(ctx, cd.activityRunScript, s).Get(ctx, &resErr)
		if errActivity != nil {
			cd.logger.Error("error run activity", zap.String("script", s), zap.Error(errActivity))
			continue
		}
		if resErr != nil {
			cd.logger.Error("activity returns error", zap.String("script", s), zap.Error(resErr))
			continue
		}
	}

	return nil
}

func (cd *Cadence) activityRunScript(name string) error {
	cd.logger.Debug("run activity", zap.String("script", name))

	return cd.rnr.RunScript(name, nil)
}

func buildCadenceClient(w cadence.Worker) (workflowserviceclient.Interface, error) {
	ch, errNewChannel := tchannel.NewChannelTransport(tchannel.ServiceName(w.ClientName))
	if errNewChannel != nil {
		return nil, fmt.Errorf("error create channel transport, %w", errNewChannel)
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: w.ClientName,
		Outbounds: yarpc.Outbounds{
			w.CadenceService: {Unary: ch.NewSingleOutbound(w.Host)},
		},
	})
	errStartDispatcher := dispatcher.Start()
	if errStartDispatcher != nil {
		return nil, fmt.Errorf("error start dispatcher, %w", errStartDispatcher)
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(w.CadenceService)), nil
}
