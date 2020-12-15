// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
	process "github.com/shirou/gopsutil/process"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/nux"
	"log"
	"strings"
)

func ProcMetrics() (L []*model.MetricValue) {

	reportProcs := g.ReportProcs()
	sz := len(reportProcs)
	if sz == 0 {
		return
	}

	ps, err := nux.AllProcs()
	if err != nil {
		log.Println(err)
		return
	}

	pslen := len(ps)

	for tags, m := range reportProcs {
		cnt := 0
		for i := 0; i < pslen; i++ {
			if is_a(ps[i], m) {
				cnt++
				procCache := m
				mem, err := procCache.MemoryInfo()
				if err != nil {
					log.Println("获取内存错误",err)
					continue
				}
				memUsedTotal += mem.RSS
				memUtil, err := procCache.MemoryPercent()
			if err != nil {
				log.Println(err)
				continue
			}
			memUtilTotal += float64(memUtil)
			cpuUtil, err := procCache.Percent(0)
			if err != nil {
				log.Println(err)
				continue
			}
			cpuUtilTotal += cpuUtil
		}

		L = append(L, GaugeValue(g.PROC_NUM, cnt, tags))
		L = append(L, GaugeValue("proc.mem.used", memUsedTotal, tags))
		L = append(L, GaugeValue("proc.mem.util", memUtilTotal, tags))
		L = append(L, GaugeValue("proc.cpu.util", cpuUtilTotal, tags))
	}

	return
}

func is_a(p *nux.Proc, m map[int]string) bool {
	// only one kv pair
	for key, val := range m {
		if key == 1 {
			// name
			if val != p.Name {
				return false
			}
		} else if key == 2 {
			// cmdline
			if !strings.Contains(p.Cmdline, val) {
				return false
			}
		}
	}
	return true
}
