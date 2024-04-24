// Copyright 2024 EMQ Technologies Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package node

import (
	"fmt"

	"github.com/lf-edge/ekuiper/contract/v2/api"
	"github.com/lf-edge/ekuiper/v2/internal/compressor"
	"github.com/lf-edge/ekuiper/v2/internal/pkg/def"
	"github.com/lf-edge/ekuiper/v2/pkg/infra"
	"github.com/lf-edge/ekuiper/v2/pkg/message"
)

type CompressOp struct {
	*defaultSinkNode
	tool message.Compressor
}

func NewCompressOp(name string, rOpt *def.RuleOption, compressMethod string) (*CompressOp, error) {
	dc, err := compressor.GetCompressor(compressMethod)
	if err != nil {
		return nil, fmt.Errorf("get compressor %s fail with error: %v", compressMethod, err)
	}
	return &CompressOp{
		defaultSinkNode: newDefaultSinkNode(name, rOpt),
		tool:            dc,
	}, nil
}

func (o *CompressOp) Exec(ctx api.StreamContext, errCh chan<- error) {
	o.prepareExec(ctx, errCh, "op")
	go func() {
		err := infra.SafeRun(func() error {
			runWithOrder(ctx, o.defaultSinkNode, o.concurrency, o.Worker)
			return nil
		})
		if err != nil {
			infra.DrainError(ctx, err, errCh)
		}
	}()
}

func (o *CompressOp) Worker(_ api.Logger, item any) []any {
	o.statManager.ProcessTimeStart()
	defer o.statManager.ProcessTimeEnd()
	switch d := item.(type) {
	case []byte:
		if r, err := o.tool.Compress(d); err != nil {
			return []any{err}
		} else {
			return []any{r}
		}
	default:
		return []any{fmt.Errorf("unsupported data received: %v", d)}
	}
}