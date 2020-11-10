package rxgo

import (
	"context"
	"reflect"
	"sync"
	"time"
)
// filter node implementation of streamOperator
type filOperater struct {
	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
}

func (ftop filOperater) op(ctx context.Context, o *Observable) {
	// must hold defintion of flow resourcs here, such as chan etc., that is allocated when connected
	// this resurces may be changed when operation routine is running.
	in := o.pred.outflow
	out := o.outflow
	//fmt.Println(o.name, "operator in/out chan ", in, out)
	var wg sync.WaitGroup
	if o.computation{	
		wg.Add(1)
		go func(){
			defer wg.Done()
			for o.computation{
				select {
				case <-ctx.Done():
					return
				case <-time.After(o.timespan):
					if o.flip!=nil{
						buffer,_:=o.flip.([]interface{})//通过断言实现类型转换
						for _,v := range buffer{
							o.sendToFlow(ctx, v, out)
						}
						o.flip=nil
					}
				}
			}	
		}()	
	}
	go func() {
		end := false
		for x := range in {
			if end {
				continue
			}
			// can not pass a interface as parameter (pointer) to gorountion for it may change its value outside!
			xv := reflect.ValueOf(x)
			// send an error to stream if the flip not accept error
			if e, ok := x.(error); ok && !o.flip_accept_error {
				o.sendToFlow(ctx, e, out)
				continue
			}
			// scheduler
			switch threading := o.threading; threading {
			case ThreadingDefault:
				if ftop.opFunc(ctx, o, xv, out) {
					end = true
				}
			case ThreadingIO:
				fallthrough
			case ThreadingComputing:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if ftop.opFunc(ctx, o, xv, out) {
						end = true
					}
				}()
			default:
			}
		}
		o.computation=false
		
		wg.Wait() //waiting all go-routines completed
		if o.flip!=nil{
			buffer,_:=o.flip.([]interface{})//通过断言实现类型转换
			for _,v := range buffer{
				o.sendToFlow(ctx, v, out)
			}
		}		
		o.closeFlow(out)
	}()
}


func (parent *Observable) Debounce(timespan time.Duration) (o *Observable) {
	o = parent.newTransformObservable("debounce")
	
	var latest interface{}
	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		latest=x.Interface()
		go func() {
			for{
				select {
				case <-ctx.Done():
					return
				case <-time.After(timespan):
					if latest==x.Interface() {
						o.sendToFlow(ctx, x.Interface(), out)
						return 
					}
				}
			}
		}()
			
		return
	}}
	return o
}

func (parent *Observable) Distinct() (o *Observable) {
	o = parent.newTransformObservable("distinct")

	var slice []interface{}
	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		flag := true
        for i := range slice{
            if x.Interface() == slice[i] {
                flag = false  // 存在重复元素，标识为false
                break
            }
        }
        if flag {  // 标识为false，不添加进结果
			o.sendToFlow(ctx, x.Interface(), out)
            slice = append(slice, x.Interface())

        }
		return
	}}
	return o
}

func (parent *Observable) ElementAt(index int) (o *Observable) {
	o = parent.newTransformObservable("elementat")

	takeCount := 0
	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		if takeCount==index {
			o.sendToFlow(ctx, x.Interface(), out)
			return true
		}
		takeCount++
		return
	}}
	return o
}

func (parent *Observable) First() (o *Observable) {
	o = parent.newTransformObservable("first")

	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		o.sendToFlow(ctx, x.Interface(), out)
		return true
	}}
	return o
}


func (parent *Observable) IgnoreElements() (o *Observable) {
	o = parent.newTransformObservable("ignoreElements")

	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		return
	}}
	return o
}


func (parent *Observable) Last() (o *Observable) {
	o = parent.newTransformObservable("last")

	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		var slice []interface{}
		o.flip=append(slice,x.Interface())
		return 
	}}
	return o
}


func (parent *Observable) Sample(timespan time.Duration) (o *Observable) {
	o = parent.newTransformObservable("sample")
	o.computation = true
	o.timespan = timespan
	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		var slice []interface{}
		o.flip=append(slice,x.Interface())
		return
	}}
	return o
}

func (parent *Observable) Skip(n int) (o *Observable) {
	o = parent.newTransformObservable("skip")
	skipCount:=0
	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		skipCount++;
		if skipCount>n{
			o.sendToFlow(ctx, x.Interface(), out)
		}
		return
	}}
	return o
}

func (parent *Observable) SkipLast(n int) (o *Observable) {
	o = parent.newTransformObservable("skipLast")
	var slice []interface{}
	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		slice=append(slice,x.Interface())
		if len(slice)>n{
			o.flip=slice[0:len(slice)-n]
		} 		
		return 
	}}
	return o
}

func (parent *Observable) Take(n int) (o *Observable) {
	o = parent.newTransformObservable("take")
	takeCount :=0
	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		takeCount++
		if takeCount>n {return true}
		o.sendToFlow(ctx, x.Interface(), out)
		return
	}}
	return o
}

func (parent *Observable) TakeLast(n int) (o *Observable) {
	o = parent.newTransformObservable("takeLast")

	var slice []interface{}
	o.operator  = filOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
		slice=append(slice,x.Interface())
		if len(slice)<n{
			o.flip=slice
		}else{
			o.flip=slice[len(slice)-n:len(slice)]
		} 		
		return 
	}}
	return o
}



