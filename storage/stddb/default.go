package stddb

import (
	"github.com/MetaParticle/metaparticle/storage"
	"github.com/MetaParticle/metaparticle/storage/mgodb"
)

func Connect(db, host string) storage.DB {
	return mgodb.ConnectMgo(db, host)
}
