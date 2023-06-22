//go:build rocksdb
// +build rocksdb

package cmd

import (
	"sort"

	versiondbclient "github.com/crypto-org-chain/cronos/versiondb/client"
	"github.com/linxGnu/grocksdb"
	"github.com/spf13/cobra"

	"github.com/numiadata/tools/cosi/utils/opendb"
)

func ChangeSetCmd() *cobra.Command {
	// keys, _, _ := app.StoreKeys()

	// add a flag for certain chains to allow users to index changesets
	storeNames := make([]string, 0, len("123")) //TODO
	// for name := range keys {
	// 	storeNames = append(storeNames, name)
	// }
	sort.Strings(storeNames)

	return versiondbclient.ChangeSetGroupCmd(versiondbclient.Options{
		DefaultStores:  storeNames,
		OpenReadOnlyDB: opendb.OpenReadOnlyDB,
		AppRocksDBOptions: func(sstFileWriter bool) *grocksdb.Options {
			return opendb.NewRocksdbOptions(nil, sstFileWriter)
		},
	})
}
