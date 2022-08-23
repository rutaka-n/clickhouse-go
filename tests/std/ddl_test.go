// Licensed to ClickHouse, Inc. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. ClickHouse, Inc. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package std

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuotedDDL(t *testing.T) {

	dsns := map[string]clickhouse.Protocol{"Native": clickhouse.Native, "Http": clickhouse.HTTP}
	ctx := context.Background()
	for name, protocol := range dsns {
		t.Run(fmt.Sprintf("%s Protocol", name), func(t *testing.T) {
			conn, err := GetStdDSNConnection(protocol, false, "false")
			require.NoError(t, err)
			require.NoError(t, conn.PingContext(context.Background()))
			require.NoError(t, err)
			require.NoError(t, conn.Ping())
			conn.Exec("DROP TABLE `test_string`")
			defer func() {
				conn.Exec("DROP TABLE `test_string`")
			}()
			_, err = conn.Exec("CREATE TABLE `test_string` (`1` String) Engine Memory")
			require.NoError(t, err)
			scope, err := conn.Begin()
			require.NoError(t, err)
			batch, err := scope.PrepareContext(ctx, "INSERT INTO `test_string`")
			require.NoError(t, err)
			_, err = batch.Exec("A")
			require.NoError(t, err)
			require.NoError(t, scope.Commit())
		})
	}
}
