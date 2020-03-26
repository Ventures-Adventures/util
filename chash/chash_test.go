package chash

import (
	"fmt"
	"time"
	"testing"
)

func Test_chash(t *testing.T){
	chash := NewCHash(10000)

	nodes := []string{"node1", "node2", "node3"}
	for _, node := range nodes{
		chash.Add(node)
	}
	
	var M = make(map[string]int)

	for i:=0; i<10000000; i++{
		node := chash.Get(fmt.Sprintf("%d", i))
		M[node]++
	}
	for k,v:= range M{
		fmt.Printf("%s:%d\n", k, v)
	}
}
func Test_chash_time(t *testing.T){
	chash := NewCHash(100000)

	nodes := []string{"node1", "node2", "node3"}
	for _, node := range nodes{
		chash.Add(node)
	}
	
	t0 := time.Now()

	for i:=0; i<10000000; i++{
		_ = chash.Get(fmt.Sprintf("%d", i))
	}
	t1 := time.Now()
	d1 := t1.Sub(t0)
	d2 := time.Since(t0)
	fmt.Printf("d1:%v, d2:%v, n1:%v, n2:%v\n", d1, d2, d1.Nanoseconds(), d2.Nanoseconds())
	fmt.Printf("avg:%v\n", d1.Nanoseconds()/10000000)
}