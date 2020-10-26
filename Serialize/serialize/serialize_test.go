package serialize

import (
	"reflect"
	"testing"
)

func TestJsonMarshalInt(t *testing.T) {
	res, err := JsonMarshal(-1234)
	want:=[]byte{'-','1', '2', '3', '4'};
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalUint(t *testing.T) {
	res, err := JsonMarshal(1234)
	want:=[]byte{'1', '2', '3', '4'};
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalString(t *testing.T) {
	res, err := JsonMarshal("test")
	want:=[]byte{'"', 't', 'e', 's', 't', '"'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalSlice(t *testing.T) {
	res, err := JsonMarshal([]int{})
	want:=[]byte{'[', ']'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalSlice2(t *testing.T) {
	res, err := JsonMarshal([]int{1, -2})
	want:=[]byte{'[', '1', ',', '-', '2', ']'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalArray(t *testing.T) {
	res, err := JsonMarshal([0]int{})
	want:=[]byte{'[', ']'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalArray2(t *testing.T) {
	res, err := JsonMarshal([3]uint{1, 2, 3})
	want:=[]byte{'[', '1', ',', '2', ',', '3', ']'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalEmptyStruct(t *testing.T) {
	type Empty struct{}
	res, err := JsonMarshal(Empty{})
	want:=[]byte{'{', '}'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalUnexportedFields(t *testing.T) {
	type Private struct{ a int }
	res, err := JsonMarshal(Private{})
	want:=[]byte{'{', '}'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalExportedFields(t *testing.T) {
	type Public struct{ A int }
	res, err := JsonMarshal(Public{})
	want:=[]byte{'{', '"', 'A', '"', ':', '0', '}'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalExportedFields2(t *testing.T) {
	type Public struct{ A, B int }
	res, err := JsonMarshal(Public{})
	want:=[]byte{'{', '"', 'A', '"', ':', '0', ',', '"', 'B', '"', ':', '0', '}'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalMixedFields(t *testing.T) {
	type Public struct{ A, b, C int }
	res, err := JsonMarshal(Public{})
	want:=[]byte{'{', '"', 'A', '"', ':', '0', ',', '"', 'C', '"', ':', '0', '}'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}


func TestJsonMarshalWithTag(t *testing.T) {
	type PublicWithTag struct {A int `mytag:"thisisA"`}
	res, err := JsonMarshal(PublicWithTag{A: 1})
	want:=[]byte{'{','"', 't', 'h', 'i','s','i','s','A','"', ':', '1', '}'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalInnerString(t *testing.T) {
	type O struct{ B string	}
	res, err := JsonMarshal(O{B: "test"})
	want:=[]byte{'{', '"','B','"', ':', '"', 't', 'e', 's', 't', '"', '}'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalInnerStruct(t *testing.T) {
	type I struct{ A int }
	type O struct{ B I }
	res, err := JsonMarshal(O{})
	want:=[]byte{'{', '"', 'B','"', ':', '{', '"','A', '"',':', '0', '}', '}'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}

func TestJsonMarshalInnerArray(t *testing.T) {
	type O struct{ B []int }
	res, err := JsonMarshal(O{B: []int{1, 2}})
	want:=[]byte{'{', '"','B','"', ':', '[', '1', ',', '2', ']', '}'}
	if err != nil{
		t.Errorf("JsonMarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("JsonMarshal() = %v, want %v", res, want)
	}
}


