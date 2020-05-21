package main

import (
	"strconv"
	"regexp"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

/**
 Stores information needed to initialize a node on the
 GlobalStore network 
*/
type GlobalStoreSession struct {
	Host      string
	Keyspace  string
	Session   gocqlx.Session
}

/**
 Simple extensions to the Interval class to facilitate marshaling and
 unmarshaling during DB transactions
*/
func EncodeToString(i Interval) string {
	return fmt.Sprintf("[%d,%d]", i.Lo, i.Hi)
}

func DecodeToInterval(s string) Interval {
	// Find all numbers in the string
	matches := regexp.MustCompile("[0-9]+").FindAllString(s, -1)
	
	// First one is lo, second is hi
	lo, _ := strconv.Atoi(matches[0])
	hi, _ := strconv.Atoi(matches[1])

	return Interval{lo, hi}
}

/**
 Uses a GlobalStoreSession to establish a live connection to the
 decentralized network
*/
func CreateSession(host string, keyspace string) *GlobalStoreSession {

	// Create cluster
	cluster := gocql.NewCluster(host)

	// Initialize a session
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		panic(err.Error())
	}
	
	// Ensure keyspace exists. If not, create one.
	query := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};", keyspace)

	err = session.ExecStmt(query)
	if err != nil {
		panic(err.Error())
	}

	// Ensure table exists. If not, create one.
	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.global_traces(id uuid PRIMARY KEY, point_code text, intervals list<text>)", keyspace)

	err = session.ExecStmt(query)
	if err != nil {
		panic(err.Error())
	}
	
	// Emit struct for future network interactions
	return &GlobalStoreSession {
		Host: host,
		Keyspace: keyspace,
		Session: session,
	}
}


/**
 Gracefully closes an open session. Defer this function 
 after every CreateSession() call
*/
func (gs *GlobalStoreSession) Close() {
	gs.Session.Close()
}

/**
 Adds a LocalTrace struct to the global store
*/
func (gs *GlobalStoreSession) AddToStore(trace *LocalTrace) {

	// Query builder
	targetTable := fmt.Sprintf("%s.global_traces", gs.Keyspace)
	insert := gs.Session.Query(qb.Insert(targetTable).Columns("id","point_code","intervals").ToCql())
	
	// Adds the map entries in one-by-one
	for pointCode, intervals := range trace.data {
		id, _ := gocql.RandomUUID()

		// Pack the intervals into string encodings
		// TODO: Fix this.
		var encodedIntervals []string
		for _, intrvl := range intervals {
			encodedIntervals = append(encodedIntervals, EncodeToString(intrvl))
		}
		
		insert.BindMap(qb.M {
			"id": id,
			"point_code": pointCode,
			"intervals": encodedIntervals,
		})

		err := insert.Exec()
		if err != nil {
			panic(err.Error())
		}
	}
	
}

/**
 Returns all the GlobalTraces in the global store
*/
func (gs *GlobalStoreSession) ReadAll() *LocalTrace {

	// Build the desired query
	targetTable := fmt.Sprintf("%s.global_traces", gs.Keyspace)
	query := gs.Session.Query(qb.Select(targetTable).Columns("point_code","intervals").ToCql())
	
	// Allocate query results buffer
	type Item struct {
		PointCode  string
		Intervals  []string
	}
	var items []*Item
	
	// Run query
	err := query.Select(&items)
	if err != nil {
		panic(err.Error())
	}

	// Construct output
	data := make(map[PointCode][]Interval)
	for _, item := range items {

		// Decode intervals
		var decodedIntervals []Interval
		for _, intrvl := range item.Intervals {
			decodedIntervals = append(decodedIntervals, DecodeToInterval(intrvl))
		}
		
		value := data[PointCode(item.PointCode)]
		data[PointCode(item.PointCode)] = append(value, decodedIntervals...)
	}

	return &LocalTrace {data: data}	
}

/**
 Returns all the GlobalTraces in the global store
*/
func (gs *GlobalStoreSession) Read(cell PointCode) *LocalTrace {

	// Build the desired query
	targetTable := fmt.Sprintf("%s.global_traces", gs.Keyspace)
	query := gs.Session.Query(qb.Select(targetTable).
		Columns("point_code", "intervals").
		Where(qb.EqLit("point_code", fmt.Sprintf("'%s'", cell))).
		AllowFiltering().
		ToCql())
	
	// Allocate query results buffer
	type Item struct {
		PointCode  string
		Intervals  []string
	}
	var items []*Item
	
	// Run query
	err := query.Select(&items)
	if err != nil {
		panic(err.Error())
	}

	// Construct output
	data := make(map[PointCode][]Interval)
	for _, item := range items {

		// Decode intervals
		var decodedIntervals []Interval
		for _, intrvl := range item.Intervals {
			decodedIntervals = append(decodedIntervals, DecodeToInterval(intrvl))
		}
		
		value := data[PointCode(item.PointCode)]
		data[PointCode(item.PointCode)] = append(value, decodedIntervals...)
	}

	return &LocalTrace {data: data}	
}


/**
 Test driver code for this module
*/
func main() {

	// (1) Create DB session
	host := "127.0.0.1"
	keyspace := "demo"
	session := CreateSession(host, keyspace)

	defer session.Close()
	
	// (2) Add Some entries
	trace := NewLocalTrace()

	input := map[PointCode][]Interval {
		"Loc1": []Interval {
			Interval{0,1},
			Interval{3,5},
		},
		"Loc2": []Interval {
			Interval{1,3},
		},
	}
	
	trace.data = input

	// (3) Load data into DB
	session.AddToStore(trace)

	// (4) Fetch all from DB
	output := session.ReadAll()
	fmt.Println(output.data)

	// (5) Fetch some from DB
	output = session.Read("Loc1")
	fmt.Println(output.data)
	
}
