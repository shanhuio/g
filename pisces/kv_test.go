// Copyright (C) 2023  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package pisces

import (
	"testing"

	"bytes"
	"errors"
	"reflect"
	"sort"

	"shanhu.io/g/errcode"
)

type testData struct {
	Value string
}

var kvTestSuite = []struct {
	name    string
	f       func(t *testing.T, kv *KV)
	ordered bool
}{
	{"add", testKVAdd, false},
	{"clear", testKVClear, false},
	{"remove", testKVRemove, false},
	{"remove-notfound", testKVRemoveNotFound, false},
	{"emplace", testKVEmplace, false},
	{"replace", testKVReplace, false},
	{"append", testKVAppendBytes, false},
	{"mutate", testKVMutate, false},
	{"mutate-cancel", testKVMutateCancel, false},
	{"mutate-error", testKVMutateError, false},
	{"walk", testKVWalk, false},
	{"ordered-walk", testKVOrderedWalk, true},
	{"walk-class", testKVWalkClass, false},
	{"walk-partial", testKVWalkPartial, true},
	{"walk-partial-desc", testKVWalkPartialDesc, true},
	{"walk-partial-class", testKVWalkPartialClass, true},
}

func testAdd(t *testing.T, kv *KV, k, v string) {
	if err := kv.Add(k, &testData{Value: v}); err != nil {
		t.Fatal(err)
	}
}

func testAddClass(t *testing.T, kv *KV, k, cls, v string) {
	if err := kv.AddClass(k, cls, &testData{Value: v}); err != nil {
		t.Fatal(err)
	}
}

func testAppendBytes(t *testing.T, kv *KV, k string, bs []byte) {
	if err := kv.AppendBytes(k, bs); err != nil {
		t.Fatal(err)
	}
}

func testGetBytes(t *testing.T, kv *KV, k string, bs []byte) {
	got, err := kv.GetBytes(k)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(bs, got) {
		t.Errorf("get bytes %q got %v, want %v", k, got, bs)
	}
}

func testGet(t *testing.T, kv *KV, k, v string) {
	d := new(testData)
	if err := kv.Get(k, d); err != nil {
		t.Fatal(err)
	}
	if d.Value != v {
		t.Errorf(`key %q: got %q, want %q`, k, d.Value, v)
	}
}

func testHas(t *testing.T, kv *KV, k string, has bool) {
	got, err := kv.Has(k)
	if err != nil {
		t.Fatal(err)
	}
	if got != has {
		t.Errorf(`has key %q, got %t, want %t`, k, got, has)
	}
}

func testGetNotFound(t *testing.T, kv *KV, k string) {
	d := new(testData)
	if err := kv.Get(k, d); err == nil {
		t.Errorf("get %q, want an error, got nil", k)
	} else if !errcode.IsNotFound(err) {
		t.Errorf("get %q want not found error, got %s", k, err)
	}
}

type testValueIter struct {
	values []string
}

func (it *testValueIter) iter() *Iter {
	return &Iter{
		Make: func() interface{} { return new(testData) },
		Do: func(_ string, v interface{}) error {
			it.values = append(it.values, v.(*testData).Value)
			return nil
		},
	}
}

func testListValues(t *testing.T, kv *KV) []string {
	iter := new(testValueIter)
	if err := kv.Walk(iter.iter()); err != nil {
		t.Fatal(err)
	}
	return iter.values
}

func testListPartialValues(t *testing.T, kv *KV, p *KVPartial) []string {
	iter := new(testValueIter)
	if err := kv.WalkPartial(p, iter.iter()); err != nil {
		t.Fatal(err)
	}
	return iter.values
}

func testKVAdd(t *testing.T, kv *KV) {
	testGetNotFound(t, kv, "k")
	testAdd(t, kv, "k", "v")
	testGet(t, kv, "k", "v")
	testGetNotFound(t, kv, "miss")
}

func testKVClear(t *testing.T, kv *KV) {
	testAdd(t, kv, "k1", "v1")
	testAdd(t, kv, "k2", "v2")

	testGet(t, kv, "k1", "v1")
	testGet(t, kv, "k2", "v2")
	testHas(t, kv, "k1", true)
	testHas(t, kv, "k2", true)

	if err := kv.Clear(); err != nil {
		t.Fatal(err)
	}

	testGetNotFound(t, kv, "k1")
	testGetNotFound(t, kv, "k2")
	testHas(t, kv, "k1", false)
	testHas(t, kv, "k2", false)
}

func testKVRemove(t *testing.T, kv *KV) {
	testAdd(t, kv, "k1", "v1")
	testAdd(t, kv, "k2", "v2")

	if err := kv.Remove("k1"); err != nil {
		t.Fatal(err)
	}

	testGetNotFound(t, kv, "k1")
	testGet(t, kv, "k2", "v2")
}

func testKVRemoveNotFound(t *testing.T, kv *KV) {
	err := kv.Remove("k1")
	if err == nil {
		t.Error("got nil error, want not found")
	}
	if !errcode.IsNotFound(err) {
		t.Errorf("got %q, want not found", err)
	}
}

func testKVEmplace(t *testing.T, kv *KV) {
	testAdd(t, kv, "k1", "v1")
	if err := kv.Emplace("k1", &testData{Value: "v1-other"}); err != nil {
		t.Fatal(err)
	}
	testGet(t, kv, "k1", "v1")

	testGetNotFound(t, kv, "k2")
	if err := kv.Emplace("k2", &testData{Value: "v2"}); err != nil {
		t.Fatal(err)
	}
	testGet(t, kv, "k2", "v2")
}

func testKVReplace(t *testing.T, kv *KV) {
	testAdd(t, kv, "k1", "v1")
	if err := kv.Replace("k1", &testData{Value: "v1-other"}); err != nil {
		t.Fatal(err)
	}
	testGet(t, kv, "k1", "v1-other")

	testGetNotFound(t, kv, "k2")
	if err := kv.Replace("k2", &testData{Value: "v2"}); err != nil {
		t.Fatal(err)
	}
	testGet(t, kv, "k2", "v2")
}

func testKVAppendBytes(t *testing.T, kv *KV) {
	testAppendBytes(t, kv, "k", []byte("hello"))
	testGetBytes(t, kv, "k", []byte("hello"))
	testAppendBytes(t, kv, "k", []byte("world"))
	testGetBytes(t, kv, "k", []byte("helloworld"))
}

func testKVMutate(t *testing.T, kv *KV) {
	testAdd(t, kv, "k", "v")

	if err := kv.Mutate("k", new(testData), func(v interface{}) error {
		d := v.(*testData)
		d.Value = "new-value"
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	testGet(t, kv, "k", "new-value")
}

func testKVMutateCancel(t *testing.T, kv *KV) {
	testAdd(t, kv, "k", "v")

	if err := kv.Mutate("k", new(testData), func(v interface{}) error {
		d := v.(*testData)
		d.Value = "new-value"
		return ErrCancel
	}); err != nil {
		t.Fatal(err)
	}

	testGet(t, kv, "k", "v")
}

func testKVMutateError(t *testing.T, kv *KV) {
	testAdd(t, kv, "k", "v")

	var errCustom = errors.New("custom")

	if err := kv.Mutate("k", new(testData), func(v interface{}) error {
		d := v.(*testData)
		d.Value = "new-value"
		return errCustom
	}); err != errCustom {
		t.Errorf("got %s, want error %s", err, errCustom)
	}

	testGet(t, kv, "k", "v")
}

func testKVWalk(t *testing.T, kv *KV) {
	testAdd(t, kv, "k1", "v1")
	testAdd(t, kv, "k2", "v2")

	values := testListValues(t, kv)
	sort.Strings(values)
	want := []string{"v1", "v2"}
	if !reflect.DeepEqual(values, want) {
		t.Errorf("got %v, want %v", values, want)
	}
}

func testKVOrderedWalk(t *testing.T, kv *KV) {
	testAdd(t, kv, "k1", "v1")
	testAdd(t, kv, "k3", "v3")
	testAdd(t, kv, "k2", "v2")

	values := testListValues(t, kv)
	want := []string{"v1", "v2", "v3"}
	if !reflect.DeepEqual(values, want) {
		t.Errorf("got %v, want %v", values, want)
	}
}

func testKVWalkClass(t *testing.T, kv *KV) {
	testAddClass(t, kv, "k1", "c1", "v1")
	testAddClass(t, kv, "k2", "c1", "v2")
	testAddClass(t, kv, "k3", "c3", "v3")

	iter := new(testValueIter)
	if err := kv.WalkClass("c1", iter.iter()); err != nil {
		t.Fatal(err)
	}
	sort.Strings(iter.values)
	want := []string{"v1", "v2"}
	if !reflect.DeepEqual(iter.values, want) {
		t.Errorf("got %v, want %v", iter.values, want)
	}
}

func testKVWalkPartial(t *testing.T, kv *KV) {
	testAdd(t, kv, "k1", "v1")
	testAdd(t, kv, "k3", "v3")
	testAdd(t, kv, "k2", "v2")

	values := testListPartialValues(t, kv, &KVPartial{
		Offset: 1,
		N:      1,
	})
	want := []string{"v2"}
	if !reflect.DeepEqual(values, want) {
		t.Errorf("got %v, want %v", values, want)
	}
}

func testKVWalkPartialDesc(t *testing.T, kv *KV) {
	testAdd(t, kv, "k1", "v1")
	testAdd(t, kv, "k3", "v3")
	testAdd(t, kv, "k2", "v2")

	values := testListPartialValues(t, kv, &KVPartial{
		Offset: 0,
		N:      2,
		Desc:   true,
	})
	want := []string{"v3", "v2"}
	if !reflect.DeepEqual(values, want) {
		t.Errorf("got %v, want %v", values, want)
	}
}

func testKVWalkPartialClass(t *testing.T, kv *KV) {
	testAddClass(t, kv, "k1", "odd", "v1")
	testAddClass(t, kv, "k2", "even", "v2")
	testAddClass(t, kv, "k3", "odd", "v3")
	testAddClass(t, kv, "k4", "even", "v4")
	testAddClass(t, kv, "k5", "odd", "v5")

	iter := new(testValueIter)
	if err := kv.WalkPartialClass("odd", &KVPartial{
		Offset: 0,
		N:      100,
	}, iter.iter()); err != nil {
		t.Fatal(err)
	}
	want := []string{"v1", "v3", "v5"}
	if !reflect.DeepEqual(iter.values, want) {
		t.Errorf("got %v, want %v", iter.values, want)
	}
}

func testKVCount(t *testing.T, kv *KV) {
	zero, err := kv.Count()
	if err != nil {
		t.Fatal(err)
	}
	if zero != 0 {
		t.Errorf("got %d, want 0", zero)
	}
	testAdd(t, kv, "k1", "v1")
	testAdd(t, kv, "k3", "v3")
	testAdd(t, kv, "k2", "v2")

	n, err := kv.Count()
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(3); n != want {
		t.Errorf("got %d, want %d", n, want)
	}
}
