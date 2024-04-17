// Copyright 2022-2024 EMQ Technologies Co., Ltd.
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

package context

import (
	"github.com/lf-edge/ekuiper/contract/v2/api"
	"github.com/lf-edge/ekuiper/v2/internal/conf"
	"github.com/lf-edge/ekuiper/v2/internal/topo/context"
	"github.com/lf-edge/ekuiper/v2/internal/topo/state"
	"github.com/lf-edge/ekuiper/v2/internal/topo/transform"
)

func NewMockContext(ruleId string, opId string) api.StreamContext {
	contextLogger := conf.Log.WithField("rule", ruleId)
	ctx := context.WithValue(context.Background(), context.LoggerKey, contextLogger)
	tf, _ := transform.GenTransform("", "json", "", "", "", []string{})
	ctx = context.WithValue(ctx, context.TransKey, tf)
	tempStore, _ := state.CreateStore(ruleId, api.AtMostOnce)
	return ctx.WithMeta(ruleId, opId, tempStore)
}