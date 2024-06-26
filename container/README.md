**Read this in other languages: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [container](#container)
  * [bitmap_test.go](#bitmap_testgo)
  * [buff_pool_test.go](#buff_pool_testgo)
  * [const_hash_test.go](#const_hash_testgo)
  * [lockable_test.go](#lockable_testgo)
  * [mdb_test.go](#mdb_testgo)
  * [queue_test.go](#queue_testgo)
  * [red_black_tree_test.go](#red_black_tree_testgo)

<!-- tocstop -->

# container
## bitmap_test.go
### TestBitmapExists
```go

bitmap := initTestData()
t.Log(bitmap)

if !bitmap.Exists(122) {
	t.FailNow()
}

if bitmap.Exists(123) {
	t.FailNow()
}
```
### TestBitmapSet
```go

bitmap := initTestData()

if bitmap.Exists(1256) {
	t.FailNow()
}

bitmap.Set(1256)

if !bitmap.Exists(1256) {
	t.FailNow()
}
```
### TestBitmapUnionOr
```go

bitmap := initTestData()
bitmap2 := initTestData()
bitmap2.Set(256)

bitmap3 := bitmap.Union(&bitmap2)
if !bitmap3.Exists(256) {
	t.FailNow()
}

bitmap3.Set(562)

if !bitmap3.Exists(562) {
	t.FailNow()
}

if bitmap.Exists(562) {
	t.FailNow()
}
```
### TestBitmapBitInverse
```go

bitmap := initTestData()

if !bitmap.Exists(66) {
	t.FailNow()
}

bitmap.Inverse()

if bitmap.Exists(66) {
	t.FailNow()
}
```
## buff_pool_test.go
### TestBuffPool
```go

buf1 := PoolBuffer128K.Get()
ptr1 := uintptr(unsafe.Pointer(&buf1[0]))
len1 := len(buf1)

PoolBuffer128K.Put(buf1)

buf2 := PoolBuffer128K.Get()
ptr2 := uintptr(unsafe.Pointer(&buf2[0]))
len2 := len(buf2)
PoolBuffer128K.Put(buf2)

if len1 != 128*1024 {
	t.Error("pool get BUFF_128K len failed")
}

if len1 != len2 {
	t.Error("pool get BUFF_128K len failed")
}

if ptr1 != ptr2 {
	t.Error("pool get BUFF_128K failed")
}

//4M
buf3 := PoolBuffer4M.Get()
ptr3 := uintptr(unsafe.Pointer(&buf3[0]))
len3 := len(buf3)
PoolBuffer4M.Put(buf3)

buf4 := PoolBuffer4M.Get()
ptr4 := uintptr(unsafe.Pointer(&buf4[0]))
len4 := len(buf4)
PoolBuffer4M.Put(buf4)

if len3 != 4*1024*1024 {
	t.Error("pool get BUFF_4M len failed")
}

if len3 != len4 {
	t.Error("pool get BUFF_4M len failed")
}

if ptr3 != ptr4 {
	t.Error("pool get BUFF_4M failed")
}
```
## const_hash_test.go
### TestConstHash
```go


var ringchash CHashRing

var configs []CHashNode
for i := 0; i < 10; i++ {
	configs = append(configs, testConstHashNode("node"+strconv.Itoa(i)))
}

ringchash.Adds(configs)

t.Log("init:", ringchash.Debug())

if ringchash.GetByC32(100, false).Id() != "node0" {
	t.Fail()
}

if ringchash.GetByC32(134217727, false).Id() != "node0" {
	t.Fail()
}

if ringchash.GetByC32(134217728, false).Id() != "node8" {
	t.Fail()
}

var configs2 []CHashNode
for i := 0; i < 2; i++ {
	configs2 = append(configs2, testConstHashNode("node"+strconv.Itoa(10+i)))
}
ringchash.Adds(configs2)

t.Log("add 2 nodes", ringchash.Debug())

if ringchash.GetByC32(134217727, false).Id() != "node10" {
	t.Fail()
}

if ringchash.GetByC32(134217728, false).Id() != "node10" {
	t.Fail()
}

ringchash.Del("node0")
t.Log("del 1 node", ringchash.Debug())

if ringchash.GetByC32(100, false).Id() != "node10" {
	t.Fail()
}

t.Log(ringchash.GetByKey("goutils", false))
```
## lockable_test.go
### TestLockable
```go

var i Lockable[int]
i.Set(100)
if i.Get() != 100 {
	t.Error(i.Get())
}

var wg sync.WaitGroup
wg.Add(3)
go func() {
	i.Update(func(i int) int { return i + 1 })
	wg.Done()
}()
go func() {
	i.Update(func(i int) int { return i + 2 })
	wg.Done()
}()
go func() {
	i.Update(func(i int) int { return i - 3 })
	wg.Done()
}()
wg.Wait()

if i.Get() != 100 {
	t.Error(i.Get())
}
```
## mdb_test.go
### TestDataTable
```go

if len(testDt.Rows()) != 10 {
	t.Error(len(testDt.Rows()))
}

if !reflect.DeepEqual(testDt.Cols(), []string{"id", "code", "name"}) {
	t.Error(testDt.Cols())
}

if testDt.PkCol() != "id" {
	t.Error(testDt.PkCol())
}

if !reflect.DeepEqual(testDt.Indexes(), []string{"code"}) {
	t.Error(testDt.Indexes())
}

if testDt.PkString(testDt.Row("9")) != "9" {
	t.Error(testDt.PkString(testDt.Row("9")))
}

if testDt.PkInt(testDt.Row("8")) != 8 {
	t.Error(testDt.PkInt(testDt.Row("8")))
}

if testDt.Row("9").Int64("id") != 9 {
	t.Error(testDt.Row("9").Int64("id"))
}

if testDt.Row("9").UInt64("id") != 9 {
	t.Error(testDt.Row("9").Int64("id"))
}

if testDt.Row("9").String("code") != "C9" {
	t.Error(testDt.Row("9").String("code"))
}

if testDt.Row("9").String("nop") != "" {
	t.Error(testDt.Row("9").String("nop"))
}

if !reflect.DeepEqual(testDt.Row("2").Data(), []string{"2", "C2", "N2"}) {
	t.Error(testDt.Row("2").Data())
}

if !reflect.DeepEqual(testDt.RowsBy("code", "C2")[0].Data(), []string{"2", "C2", "N2"}) {
	t.Error(testDt.RowsBy("code", "C2")[0].Data())
}

if !reflect.DeepEqual(testDt.RowsByPredicate(func(dr *DataRow) bool { return dr.String("name") == "N4" })[0].Data(), []string{"4", "C4", "N4"}) {
	t.Error("RowsByPredicate")
}

testDt.Push([]string{"2", "C2", "N3"})

if !reflect.DeepEqual(testDt.RowsByIndexPredicate("code", "C2", func(dr *DataRow) bool { return dr.String("name") == "N3" })[0].Data(), []string{"2", "C2", "N3"}) {
	t.Error("RowsByIndexPredicate")
}
```
### TestDataSet
```go

if testDs.Table("testDt") != testDt {
	t.FailNow()
}
```
## queue_test.go
### TestEnqueueBack
```go

q := InitQueue()
if !q.IsEmpty() {
	t.FailNow()
}

for i := 0; i < 25; i++ {
	deItem := q.EnqueueBack(fmt.Sprint(i))

	if i >= q.Len() {
		expectedItem := fmt.Sprint(i - q.Len())
		if deItem != expectedItem {
			t.Error(deItem)
		}
	} else {
		if deItem != "" {
			t.Error(deItem)
		}
	}
}

if q.IsEmpty() {
	t.FailNow()
}

// q.Range(func(i string) bool {
// 	t.Log("left:", i)
// 	return true
// })

```
### TestDequeueFront
```go

q := InitQueue()
if !q.IsEmpty() {
	t.FailNow()
}

for i := 0; i < 25; i++ {
	q.EnqueueBack(fmt.Sprint(i))
}

for i := 0; i < 5; i++ {
	deItem := q.DequeueFront()
	expectedItem := fmt.Sprint(i + 25 - 10)
	if deItem != expectedItem {
		t.Error(deItem)
	}
}

// q.Range(func(i string) bool {
// 	t.Log("left:", i)
// 	return true
// })

if q.IsEmpty() {
	t.FailNow()
}

for i := 0; i < 5; i++ {
	deItem := q.DequeueFront()
	expectedItem := fmt.Sprint(i + 25 - 10 + 5)
	if deItem != expectedItem {
		t.Error(deItem)
	}
}

// q.Range(func(i string) bool {
// 	t.Log("left:", i)
// 	return true
// })

if !q.IsEmpty() {
	t.FailNow()
}
```
### TestEnqueueFront
```go

q := InitQueue()
if !q.IsEmpty() {
	t.FailNow()
}

for i := 0; i < 25; i++ {
	deItem := q.EnqueueFront(fmt.Sprint(i))
	if i >= q.Len() {
		expectedItem := fmt.Sprint(i - q.Len())
		if deItem != expectedItem {
			t.Error(deItem)
		}
	} else {
		if deItem != "" {
			t.Error(deItem)
		}
	}
}

if q.IsEmpty() {
	t.FailNow()
}

// q.Range(func(i string) bool {
// 	t.Log("left:", i)
// 	return true
// })
```
### TestDequeueBack
```go

q := InitQueue()
if !q.IsEmpty() {
	t.FailNow()
}

for i := 0; i < 25; i++ {
	q.EnqueueFront(fmt.Sprint(i))
}

for i := 0; i < 5; i++ {
	deItem := q.DequeueBack()
	expectedItem := fmt.Sprint(i + 25 - 10)
	if deItem != expectedItem {
		t.Error(deItem)
	}
}
if q.IsEmpty() {
	t.FailNow()
}

// q.Range(func(i string) bool {
// 	t.Log("left:", i)
// 	return true
// })

for i := 0; i < 5; i++ {
	deItem := q.DequeueBack()
	expectedItem := fmt.Sprint(i + 25 - 10 + 5)
	if deItem != expectedItem {
		t.Error(deItem)
	}
}

if !q.IsEmpty() {
	t.FailNow()
}
```
### TestQueueClear
```go

q := InitQueue()
if q.Cap() != 10 {
	t.FailNow()
}

if !q.IsEmpty() {
	t.FailNow()
}

for i := 0; i < 25; i++ {
	q.EnqueueBack(fmt.Sprint(i))
}

// q.Range(func(i string) bool {
// 	t.Log("left:", i)
// 	return true
// })

q.Clear()

if q.Cap() != 10 {
	t.FailNow()
}

if !q.IsEmpty() {
	t.FailNow()
}

// q.Range(func(i string) bool {
// 	t.Log("left:", i)
// 	return true
// })
```
### TestQueueFindBy
```go

q := NewQueue[int](25)

for i := 0; i < 25; i++ {
	q.EnqueueBack(i)
}

i20 := q.FindOneBy(func(i int) bool {
	return i == 20
})

if i20 != 20 {
	t.FailNow()
}

items := q.FindBy(func(i int) bool {
	return i%3 == 0
})

for _, item := range items {
	if item%3 != 0 {
		t.Error(item)
	}
}
```
### TestQueueRange
```go

q := NewQueue[int](10)
for i := 0; i < 25; i++ {
	q.EnqueueBack(i)
}

j := 15
q.Range(func(i int) bool {
	if i != j {
		t.Error(i)
	}
	j++
	return true
})
```
## red_black_tree_test.go
### TestReaBlackTree
```go


type personT struct {
	name   string
	age    int
	gender bool
	score  int64
}

tree := RedBlackTree{}

for i := 0; i < 1000000; i++ {

	name := fmt.Sprintf("rongo%d", i)

	tree.Put(name, &personT{
		name:   name,
		age:    i,
		gender: true,
		score:  int64(i) + 100,
	})
}

nodeVal := tree.Get("rongo999999")

personVal, ok := nodeVal.(*personT)
if !ok {
	t.FailNow()
}

if personVal.name != "rongo999999" {
	t.FailNow()
}

if personVal.age != 999999 {
	t.FailNow()
}

if personVal.score != 1000099 {
	t.FailNow()
}

```
