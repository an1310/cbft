//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

package main

import (
	"os"

	"github.com/blevesearch/bleve"

	log "github.com/couchbaselabs/clog"
)

func (pindex *PIndex) Run() {
	if pindex.IndexType == "bleve" {
		err := RunBleveStream(pindex.Stream, pindex.Impl.(bleve.Index))
		if err != nil {
			log.Printf("error: RunBleveStream, err: %v", err)
			return
		}
	}

	// NOTE: We expect the PIndexImpl to handle any inflight, concurrent
	// queries, access and Close() correctly with its own locking.
	pindex.Impl.Close()

	os.RemoveAll(pindex.Path)
}

func RunBleveStream(stream Stream, bindex bleve.Index) error {
	for m := range stream {
		// TODO: probably need things like stream reset/rollback
		// and snapshot kinds of ops here, too.

		// TODO: maybe need a more batchy API?  Perhaps, yet another
		// goroutine that clumps up up updates into bigger batches?

		switch m := m.(type) {
		case *StreamUpdate:
			bindex.Index(string(m.Id()), m.Body())
		case *StreamDelete:
			bindex.Delete(string(m.Id()))
		}
	}

	return nil
}
