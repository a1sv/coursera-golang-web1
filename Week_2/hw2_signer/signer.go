package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// SingleHash считает значение crc32(data)+"~"+crc32(md5(data))
// ( конкатенация двух строк через ~), где data - то что пришло на вход (по сути - числа из первой функции)
func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for data := range in {
		value, ok := data.(string)
		if !ok {
			value = strconv.Itoa(data.(int))
		}

		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			dataHash1 := make(chan string)
			dataHash2 := make(chan string)

			go func() {
				dataHash1 <- DataSignerCrc32(data)
			}()

			go func() {
				mu.Lock()
				md5hash := DataSignerMd5(data)
				mu.Unlock()
				dataHash2 <- DataSignerCrc32(md5hash)
			}()

			out <- <-dataHash1 + "~" + <-dataHash2
		}(value)
	}

	wg.Wait()
}

// MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки), где th=0..5 ( т.е. 6 хешей на каждое входящее значение ),
// потом берёт конкатенацию результатов в порядке расчета (0..5), где data - то что пришло на вход (и ушло на выход из SingleHash)
func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for data := range in {
		value, ok := data.(string)
		if !ok {
			value = strconv.Itoa(data.(int))
		}

		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			workers := &sync.WaitGroup{}
			dataHashes := make([]string, 6)
			for th := 0; th < 6; th++ {
				workers.Add(1)
				go func(th int) {
					defer workers.Done()
					dataHashes[th] = DataSignerCrc32(strconv.Itoa(th) + data)
				}(th)
			}
			workers.Wait()
			out <- strings.Join(dataHashes, "")
		}(value)
	}

	wg.Wait()
}

// CombineResults получает все результаты,
// сортирует (https://golang.org/pkg/sort/), объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
//
func CombineResults(in, out chan interface{}) {
	results := make([]string, 0, 5)
	for data := range in {
		results = append(results, data.(string))
	}
	sort.Strings(results)
	result := strings.Join(results, "_")
	out <- result
}

// ExecutePipeline выполняет массив всех пришедших функций
func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	out := make(chan interface{})
	for _, j := range jobs {
		wg.Add(1)
		go jobWorker(wg, j, in, out)
		in = out
		out = make(chan interface{})
	}
	wg.Wait()
	close(out)
}

func jobWorker(wg *sync.WaitGroup, j job, in, out chan interface{}) {
	defer wg.Done()
	j(in, out)
	close(out)
}
