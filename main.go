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

package main

import (
	"os"

	"github.com/intelsdi-x/snap-plugin-collector-cputilization/cpuutil"
	"github.com/intelsdi-x/snap/control/plugin"
)

func main() {
	cpuUtil := cpuutil.New()
	if cpuUtil == nil {
		panic("Failed to initialize plugin!\n")
	}

//	stats, err := cpuutil.GetUtilStat()
//	if err != nil {
//		panic(err)
//	}
//	for k, v := range stats {
//		fmt.Printf("%v = %v\n", k, v)
//	}
//
//	mets, err := cpuUtil.GetMetricTypes(plugin.PluginConfigType{})
//	if err != nil {
//		panic(err)
//	}
//
//	mmm, err := cpuUtil.CollectMetrics(mets)
//	for _, m := range mmm {
//		fmt.Println(m.Namespace(), m.Data())
//	}
//
//	panic("blah")
	plugin.Start(
		plugin.NewPluginMeta(
			cpuutil.NAME,
			cpuutil.VERSION,
			plugin.CollectorPluginType,
			[]string{},
			[]string{plugin.SnapGOBContentType},
			plugin.ConcurrencyCount(1)),
		cpuUtil,
		os.Args[1],
	)
}
