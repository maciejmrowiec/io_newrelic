package main

import (
	"bufio"
	"log"
	"os/exec"
)

type IItem interface {
	GetName() string
	Aggregate(item IItem)
}

type ISample interface {
	Append(item IItem)
	Empty() bool
	GetMap() map[string]IItem
}

type ISampleProcessor interface {
	GetCmdName() string
	GetCmdArgs() []string
	AggregateSample(sample ISample)
	NewSample() ISample
	IsHeader(row string) bool
	ParseRow(row string) (IItem, error)
}

type ISampleStats interface {
	Append(value float64)
	GetAverage() float64
}

type DynamicCollector struct {
	processor ISampleProcessor
}

func NewDynamicCollector(processor ISampleProcessor) *DynamicCollector {
	return &DynamicCollector{
		processor: processor,
	}
}

func (d *DynamicCollector) Run() {
	pipe := make(chan string, 1000)
	go d.executeCmd(d.processor.GetCmdName(), d.processor.GetCmdArgs(), pipe)
	go d.processOutput(pipe)
}

func (d *DynamicCollector) processOutput(ch <-chan string) {
	sample := d.processor.NewSample()

	for row := range ch {

		if d.processor.IsHeader(row) {
			if sample != nil && !sample.Empty() {
				d.processor.AggregateSample(sample)
			}

			sample = d.processor.NewSample()
			continue
		}

		p, err := d.processor.ParseRow(row)
		if err != nil {
			continue
		}

		sample.Append(p)
	}
}

func (d *DynamicCollector) executeCmd(name string, args []string, ch chan<- string) {

	execute_cmd := func() {
		cmd := exec.Command(name, args...)
		stdout, err := cmd.StdoutPipe()

		if err != nil {
			log.Fatal(err)
		}

		if err = cmd.Start(); err != nil {
			log.Fatal(err)
		}

		in := bufio.NewScanner(stdout)

		for in.Scan() {
			ch <- in.Text()
		}
	}

	for true {
		execute_cmd()
	}

}
