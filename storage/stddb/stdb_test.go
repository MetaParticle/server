package stddb

import (
	"github.com/MetaParticle/metaparticle/testing"
	"github.com/MetaParticle/metaparticle/entity"
)

func Test_PlayerStorage(t *testing.T) {
	p := entity.GetTestPlayer()
	db := stddb.Connect("mp", "")
	set := db.Get("players")
	set.Insert(p)
	o := new(Player)
	o.Id = 1
	o = set.Fill(o)
	q := set.GetOne(new(Player), "Id", "1")
	
	if p.Hashpass != o.Hashpass {
		t.Fail()
	}
	
	if p.Hashpass != q.Hashpass {
		t.Fail()
	}
	
	if q.Hashpass != o.Hashpass {
		t.Fail()
	}
}

