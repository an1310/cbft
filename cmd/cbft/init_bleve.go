// Copyright (c) 2016 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License. You
// may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"runtime"
	"strconv"

	"github.com/blevesearch/bleve"
	bleveMapping "github.com/blevesearch/bleve/mapping"
	bleveSearcher "github.com/blevesearch/bleve/search/searcher"

	"github.com/couchbase/cbft"
)

func InitBleveOptions(options map[string]string) error {
	bleveMapping.StoreDynamic = false
	bleveMapping.MappingJSONStrict = true
	bleveSearcher.DisjunctionMaxClauseCount = 1024

	bleveKVStoreMetricsAllow := options["bleveKVStoreMetricsAllow"]
	if bleveKVStoreMetricsAllow != "" {
		v, err := strconv.ParseBool(bleveKVStoreMetricsAllow)
		if err != nil {
			return err
		}

		cbft.BleveKVStoreMetricsAllow = v
	}

	bleveMaxOpsPerBatch := options["bleveMaxOpsPerBatch"]
	if bleveMaxOpsPerBatch != "" {
		v, err := strconv.Atoi(bleveMaxOpsPerBatch)
		if err != nil {
			return err
		}

		cbft.BleveMaxOpsPerBatch = v
	}

	bleveAnalysisQueueSize := runtime.NumCPU()

	bleveAnalysisQueueSizeStr := options["bleveAnalysisQueueSize"]
	if bleveAnalysisQueueSizeStr != "" {
		v, err := strconv.Atoi(bleveAnalysisQueueSizeStr)
		if err != nil {
			return err
		}

		if v > 0 {
			bleveAnalysisQueueSize = v
		} else {
			bleveAnalysisQueueSize = bleveAnalysisQueueSize - v
		}
	}

	if bleveAnalysisQueueSize < 1 {
		bleveAnalysisQueueSize = 1
	}

	bleve.Config.SetAnalysisQueueSize(bleveAnalysisQueueSize)

	return nil
}
