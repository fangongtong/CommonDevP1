package bs

import (
	"context"
	"fmt"
	"testing"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	api "github.com/influxdata/influxdb-client-go/v2/api"
	influxdb "github.com/influxdata/influxdb1-client/v2"
)

func test5(t *testing.T) {

	client := influxdb2.NewClient("http://192.168.1.133:8086", "")
	client.Options().SetPrecision(time.Millisecond)
	writeAPI := client.WriteAPI("", "RaspDb")

	timeStart := time.Now()
	uid := timeStart.Format("20060102150405")
	fmt.Println(uid)

	for i := 0; i < 100; i++ {
		T6(writeAPI, i, uid)
		time.Sleep(time.Millisecond * 2)
	}
	writeAPI.Flush()
	fmt.Println(time.Now().Sub(timeStart).Milliseconds())

}

func T5(writeAPI api.WriteAPI, idx int, uid string) {
	tNow := time.Now()
	tags := map[string]string{
		"Uid":      uid,
		"PointIdx": fmt.Sprint(idx),
	}
	fields1 := map[string]interface{}{
		"TskInfo": "{..}",
	}
	fields2 := map[string]interface{}{
		"Pos":         int32(1),
		"Wp2PeekPull": float32(1.1),
		"Wp2PeekPush": float32(1.2),
	}
	fields3 := map[string]interface{}{
		"Pos":         int32(2),
		"Wp2PeekPull": float32(2.1),
		"Wp2PeekPush": float32(2.2),
	}

	// create point
	p := influxdb2.NewPoint("TskPoint", tags, fields1, tNow)
	writeAPI.WritePoint(p)
	p = influxdb2.NewPoint("TskPointPos1", tags, fields2, tNow)
	writeAPI.WritePoint(p)
	p = influxdb2.NewPoint("TskPointPos2", tags, fields3, tNow)
	writeAPI.WritePoint(p)

}
func T6(writeAPI api.WriteAPI, idx int, uid string) {
	tNow := time.Now()
	tags := map[string]string{
		"Uid": uid,
	}
	fields1 := map[string]interface{}{
		"PointIdx": int32(idx),
		"TskInfo":  "{..}",
	}
	fields2 := map[string]interface{}{
		"Wp2PeekPull": float32(1.1),
		"Wp2PeekPush": float32(1.2),
	}
	fields3 := map[string]interface{}{
		"Wp2PeekPull": float32(2.1),
		"Wp2PeekPush": float32(2.2),
	}

	// create point
	p := influxdb2.NewPoint("TskPoint", tags, fields1, tNow)
	writeAPI.WritePoint(p)
	p = influxdb2.NewPoint("TskPointPos1", tags, fields2, tNow)
	writeAPI.WritePoint(p)
	p = influxdb2.NewPoint("TskPointPos2", tags, fields3, tNow)
	writeAPI.WritePoint(p)

}

// `
// TskPoint = from(bucket:"RaspDb")
// 	|> range(start: -24h)
// 	|> filter(fn: (r) => r._measurement == "TskPoint" and r.Uid == "20220513020830" )
// 	|> drop(columns: ["_start", "_stop"])
// 	|> pivot(rowKey:["Uid", "PointIdx"], columnKey: ["_measurement", "_field"], valueColumn: "_value")

// TskPointPos1 = from(bucket:"RaspDb")
// 	|> range(start: -24h)
// 	|> filter(fn: (r) => r._measurement == "TskPointPos1" and r._field =~ /Pos|Wp2PeekPull|Wp2PeekPush/  and r.Uid == "20220513020830")
// 	|> drop(columns: ["_start", "_stop"])
// 	|> pivot(rowKey:["Uid", "PointIdx"], columnKey: ["_measurement", "_field"], valueColumn: "_value")

// joined = join(tables: {TskPoint: TskPoint, TskPointPos1: TskPointPos1}, on: ["PointIdx", "Uid"])

// TskPointPos2 = from(bucket:"RaspDb")
// 	|> range(start: -24h)
// 	|> filter(fn: (r) => r._measurement == "TskPointPos2" and r._field =~ /Pos|Wp2PeekPull|Wp2PeekPush/  and r.Uid == "20220513020830")
// 	|> drop(columns: ["_start", "_stop"])
// 	|> pivot(rowKey:["Uid", "PointIdx"], columnKey: ["_measurement", "_field"], valueColumn: "_value")

// join(tables: {TskPoint: joined,  TskPointPos2: TskPointPos2}, on: ["PointIdx", "Uid"])
// 	`)
func Test3(t *testing.T) {
	client := influxdb2.NewClient("http://192.168.1.133:8086", "")
	client.Options().SetPrecision(time.Millisecond)
	queryAPI := client.QueryAPI("")
	//result, err := queryAPI.Query(context.Background(), `from(bucket:"RaspDb")|> range(start: -24h)`)

	// range说明: https://docs.influxdata.com/flux/v0.x/data-types/basic/duration/
	result, err := queryAPI.Query(context.Background(), `
		
	TskPoint = from(bucket:"RaspDb")
		|> range(start: 0)
		|> filter(fn: (r) => r._measurement == "TskPoint" and r.Uid == "20220513153258")
		|> drop(columns: ["_start", "_stop", "Uid"])
		|> pivot(rowKey:["_time"], columnKey: ["_measurement", "_field"], valueColumn: "_value")
		|> limit(n:50, offset:50)
		
	TskPointPos1 = from(bucket:"RaspDb")
		|> range(start: 0)
		|> filter(fn: (r) => r._measurement == "TskPointPos1"  and r.Uid == "20220513153258"  )
		|> drop(columns: ["_start", "_stop", "Uid"])
		|> pivot(rowKey:["_time"], columnKey: ["_measurement", "_field"], valueColumn: "_value")


	joined = join(tables: {TskPoint: TskPoint, TskPointPos1: TskPointPos1}, on: ["_time"])
		
	TskPointPos2 = from(bucket:"RaspDb")
		|> range(start: 0)
		|> filter(fn: (r) => r._measurement == "TskPointPos2"  and r.Uid == "20220513153258")
		|> drop(columns: ["_start", "_stop", "Uid"])
		|> pivot(rowKey:["_time"], columnKey: ["_measurement", "_field"], valueColumn: "_value")
		
	join(tables: {TskPoint: joined,  TskPointPos2: TskPointPos2}, on: ["_time"]) 
	`)
	/*
		 and r._field =~ /Wp2PeekPull|Wp2PeekPush/
				TskPointPos1 = from(bucket:"RaspDb")
					|> range(start: 0)
					|> filter(fn: (r) => r._measurement == "TskPointPos1" and r._field =~ /Wp2PeekPull|Wp2PeekPush/  )
					|> drop(columns: ["_start", "_stop"])
					|> pivot(rowKey:["_time"], columnKey: ["_measurement", "_field"], valueColumn: "_value")


				join(tables: {TskPoint: TskPoint, TskPointPos1: TskPointPos1}, on: ["_time"])
	*/
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				//fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// read result
			fmt.Printf("row: %s\n", result.Record().String())
			/*
				fmt.Println("\r\n------------")
				fmt.Print(",PointIdx: ", result.Record().ValueByKey("PointIdx"))
				fmt.Print(",TskPointPos1_Wp2PeekPull: ", result.Record().ValueByKey("TskPointPos1_Wp2PeekPull"))
				fmt.Print(",TskPointPos1_Wp2PeekPush: ", result.Record().ValueByKey("TskPointPos1_Wp2PeekPush"))
				fmt.Print(",TskPointPos2_Wp2PeekPull: ", result.Record().ValueByKey("TskPointPos2_Wp2PeekPull"))
				fmt.Print(",TskPointPos2_Wp2PeekPush: ", result.Record().ValueByKey("TskPointPos2_Wp2PeekPush"))
			*/
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	} else {
		fmt.Println(err.Error())
	}
	// Ensures background processes finishes
	client.Close()
}

func test1(t *testing.T) {

	// tm := []byte(time.Now().Format("20060102150405.000"))
	// tm = append(tm[:14], tm[15:]...)
	// fmt.Println(string(tm))
	// return
	var err error
	var client influxdb.Client
	var bp1, bp2 influxdb.BatchPoints

	client, err = influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: "http://192.168.1.133:8086",
	})
	if err != nil {
		fmt.Println("NewHTTPClient failed:", err.Error())
		return
	}

	bp1, err = influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  "RaspDb",
		Precision: "ms",
	})
	if err != nil {
		fmt.Println("NewBatchPoints failed:", err.Error())
		return
	}
	bp2, err = influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  "RaspDb",
		Precision: "ms",
	})
	if err != nil {
		fmt.Println("NewBatchPoints failed:", err.Error())
		return
	}

	uid := time.Now().Format("20060102150405")
	for i := 1; i < 100; i++ {
		T2(client, bp1, bp2, i, uid)
		time.Sleep(time.Millisecond * 10)
	}

}

func T2(client influxdb.Client, bp1, bp2 influxdb.BatchPoints, idx int, uid string) {
	//tm := []byte(time.Now().Format("20060102150405.000"))
	//tm = append(tm[:14], tm[15:]...)
	tNow := time.Now()
	tags := map[string]string{
		"Uid":      uid,
		"PointIdx": fmt.Sprint(idx),
	}
	fields1 := map[string]interface{}{
		"TskInfo": "",
	}
	fields2 := map[string]interface{}{
		"Pos":         1,
		"Wp2PeekPull": float32(0),
		"Wp2PeekPush": float32(0),
	}
	fields3 := map[string]interface{}{
		"Pos":         2,
		"Wp2PeekPull": float32(0),
		"Wp2PeekPush": float32(0),
	}

	//influxdb.

	{
		pt, err := influxdb.NewPoint("TskPoint", tags, fields1, tNow)
		if err != nil {
			fmt.Println("new point failed:", err.Error())
			return
		}
		bp1.AddPoint(pt)

		pt, err = influxdb.NewPoint("TskPointPos1", tags, fields2, tNow)
		if err != nil {
			fmt.Println("new point failed:", err.Error())
			return
		}
		bp1.AddPoint(pt)

		pt, err = influxdb.NewPoint("TskPointPos2", tags, fields3, tNow)
		if err != nil {
			fmt.Println("new point failed:", err.Error())
			return
		}
		bp1.AddPoint(pt)

		if err = client.Write(bp1); err != nil {
			fmt.Println("write failed:", err.Error())
			return
		}
	}

	// {
	// 	pt, err := influxdb.NewPoint("TskPointPos1", tags, fields2, tNow)

	// 	if err != nil {
	// 		fmt.Println("new point failed:", err.Error())
	// 		return
	// 	}

	// 	bp2.AddPoint(pt)
	// 	if err = client.Write(bp2); err != nil {
	// 		fmt.Println("write failed:", err.Error())
	// 		return
	// 	}
	// }

}
