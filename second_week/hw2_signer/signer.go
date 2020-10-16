package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func ExecutePipeline(hashSignJobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for _, someHashJob := range hashSignJobs {
		wg.Add(1)
		out := make(chan interface{})

		go func(someJobWork job, in chan interface{}, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)
			someJobWork(in, out)
		}(someHashJob, in, out, wg)

		in = out
	}

	wg.Wait()
}

func SingleHash(in chan interface{}, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for input := range in {
		wg.Add(1)
		input := fmt.Sprintf("%v", input)
		inputMd5 := DataSignerMd5(input)

		go func(input string, inputMd5 string) {
			defer wg.Done()

			chanLeftCrc32 := func(input string) chan string {
				resultChan := make(chan string, 1)

				go func(out chan<- string) {
					out <- DataSignerCrc32(input)
				}(resultChan)
				return resultChan

			}(input)

			rightCrc32Md5 := DataSignerCrc32(inputMd5)
			leftCrc32 := <-chanLeftCrc32

			out <- leftCrc32 + "~" + rightCrc32Md5
		}(input, inputMd5)

	}

	defer wg.Wait()
}

func MultiHash(in chan interface{}, out chan interface{}) {

	var wgOutter sync.WaitGroup
	var th = 6

	for input := range in {
		wgOutter.Add(1)

		go func(wg *sync.WaitGroup, inputData interface{}, outChan chan interface{}) {
			var wgInner sync.WaitGroup
			resultOfWork := make([]string, th)

			defer wg.Done()

			for i := 0; i < th; i++ {
				wgInner.Add(1)
				data := fmt.Sprintf("%v%v", i, inputData)

				go func(wg *sync.WaitGroup, str string, arr []string, idx int) {
					defer wg.Done()

					convertedString := DataSignerCrc32(str)
					arr[idx] = convertedString
				}(&wgInner, data, resultOfWork, i)

			}

			wgInner.Wait()
			outputData := strings.Join(resultOfWork, "")

			outChan <- outputData
		}(&wgOutter, input, out)

	}

	wgOutter.Wait()
}

func CombineResults(in, out chan interface{}) {

	var combineResults []string

	for input := range in {
		combineResults = append(combineResults, (input).(string))
	}

	sort.Strings(combineResults)
	output := strings.Join(combineResults, "_")

	out <- output
}
