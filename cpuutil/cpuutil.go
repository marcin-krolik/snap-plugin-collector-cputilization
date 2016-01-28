// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cpuutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"

	str "github.com/intelsdi-x/snap-plugin-utilities/strings"
)

const (
	// NAME of the collector plugin
	NAME = "cpu_utilization"
	// VERSION of collector plugin
	VERSION = 1

	vendor = "intel"
	fs     = "procfs"
)

var utilInfoPath = "/proc/stat"

// New returns instance of cpu freq collector plugin
func New() *cpuUtilCollector {
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}
	return &cpuUtilCollector{host: host}
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (cpf *cpuUtilCollector) GetMetricTypes(_ plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	metricTypes := []plugin.PluginMetricType{}
	_, err := GetUtilStat()
	if err != nil {
		return nil, err
	}

	metricType := plugin.PluginMetricType{Namespace_: []string{vendor, fs, NAME}}
	metricTypes = append(metricTypes, metricType)
	return metricTypes, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (cpuu *cpuUtilCollector) CollectMetrics(metricTypes []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	metrics := []plugin.PluginMetricType{}
	skip := []string{"idle", "iowait", "steal", "quest", "quest_nice"}
	utilization, err := GetUtilStat()
	if err != nil {
		return nil, err
	}

	for _, metricType := range metricTypes {
		ns := metricType.Namespace()

		curIdle := utilization["idle"] + utilization["iowait"]
		var curNonIdle uint64
		for util, val := range utilization {
			if !str.Contains(skip, util) {
				fmt.Println(util)
				curNonIdle += val
			}
		}
//		fmt.Println("idle", curIdle)
//		fmt.Println("curNon", curNonIdle)

		// calculate difference
		totald := curIdle + curNonIdle - (cpuu.idle + cpuu.nonidle)
//		fmt.Println("totald", totald)
		idled := curIdle - cpuu.idle
//		fmt.Println("idled", idled)

		utilPerc := 100 * float64(totald - idled) / float64(totald)

		// update previous values
		cpuu.idle = curIdle
		cpuu.nonidle = curNonIdle

		metric := plugin.PluginMetricType{
			Namespace_: ns,
			Source_:    cpuu.host,
			Timestamp_: time.Now(),
			Data_:      utilPerc,
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// GetConfigPolicy returns config policy
// It returns error in case retrieval was not successful
func (cpf *cpuUtilCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}

type cpuUtilCollector struct {
	host    string
	idle    uint64
	nonidle uint64
}

func GetUtilStat() (map[string]uint64, error) {
	utilization := map[string]uint64{}
	entries := []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest", "guest_nice"}

	content, err := ioutil.ReadFile(utilInfoPath)
	if err != nil {
		return nil, err
	}

	cpu := strings.Split(string(content), "\n")[0]
	utilData := strings.Fields(cpu)
	if len(utilData) != 11 {
		return nil, fmt.Errorf("Error while parsing utilization data")
	}

	for i, entry := range entries {
		val, err := strconv.ParseUint(utilData[i+1], 10, 64)
		if err != nil {
			return nil, err
		}
		utilization[entry] = val
	}

	return utilization, nil
}
