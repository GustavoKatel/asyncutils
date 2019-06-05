package executor

import "github.com/GustavoKatel/asyncutils/executor/interfaces"

type jobImpl struct {
	jobFn interfaces.JobFn
}
