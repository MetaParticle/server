package mgodb

import (
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"

	// Definition of 'DB' and 'Set' interfaces
	"github.com/MetaParticle/metaparticle/storage"
)

type MgoDB struct {
	s *mgo.Session
	db *mgo.Database
}

type MgoSet struct {
	c *mgo.Collection
}


func ConnectMgo(db, host string) storage.DB {
	if host == "" {
		host = "localhost"
	}
	session, err := mgo.Dial(host)
    if err != nil {
            panic(err)
    }
    
    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)
    
    return MgoDB{session, session.DB(db)}
}

// DB

/*
 * Returns a SubSet of the database.
 */
func (m MgoDB) Get(set string) storage.Set {
	return MgoSet{m.db.C(set)}
}

/*
 * Closes the connection to the database.
 */
func (m MgoDB) Close() error {
	m.s.Close()
	return nil
}

/*
 * Returns the name of the database.
 * Ex. "db"
 */
func (m MgoDB) String() string {
	return m.db.Name
}

// SET

/*
 * Supplied with part of a whole, it attempts to complete
 * the struct.
 */
func (s MgoSet) Fill(data interface{}) error {
	m := bson.GetSearchMap(data)
	if m == nil {
		return storage.NewStorageError("Couldn't make SearchMap from data.")
	}
	return s.c.Find(m).One(data)	
}

func (s MgoSet) Insert(data ...interface{}) error {
	return s.c.Insert(data...)
}

/*
 * Gets an arbitrary amount of results.
 */
func (s MgoSet) Get(res interface{}, where ...string) (int, error) {
	query := s.getPrep(where...)
	
	count, err := query.Count()
	if err != nil {
		return count, err
	}
	
	return count, query.All(res)
}

/*
 * Gets the first result that satisfies the conditions.
 */
func (s MgoSet) GetOne(res interface{}, where ...string) (err error) {
	return s.getPrep(where...).One(res)
}

/*
 * Prepares a get-query from the supplied conditions.
 */
func (s MgoSet) getPrep(where ...string) *mgo.Query {
	num := len(where)/2+1
	M := make(bson.M)			// map[string]interface{}
	for i := 0; i < num; i+=2 {
		M[where[i]] = where[i+1]
	}
	query := s.c.Find(M)
	return query
}

/*
 * Returns the name of the Set, including that of the Database.
 * Ex. "database.set"
 */
func (s MgoSet) String() string {
	return s.c.FullName
}
