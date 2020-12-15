package funcs

import (
	"log"
	"math"
	"nightingale-master/src/dataobj"
	"strconv"
"strings"

"github.com/toolkits/pkg/file"
"github.com/toolkits/pkg/nux"

 "github.com/open-falcon/falcon-plus/common/model"
)

func FsKernelMetrics() []*model.MetricValue {
	maxFiles, err := nux.KernelMaxFiles()
	if err != nil {
		log.Fatalln("failed to call collect KernelMaxFiles:%v\n", err)
		return nil
	}

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Fatalln("failed to call KernelAllocateFiles:%v\n", err)
		return nil
	}

	v := math.Ceil(float64(allocateFiles) * 100 / float64(maxFiles))
	return []*model.MetricValue{
		GaugeValue("sys.fs.files.max", maxFiles),
		GaugeValue("sys.fs.files.free", maxFiles-allocateFiles),
		GaugeValue("sys.fs.files.used", allocateFiles),
		GaugeValue("sys.fs.files.used.percent", v),
	}
}

func ProcsNumMetrics() []*model.MetricValue {
	var dirs []string
	num := 0
	dirs, err := file.DirsUnder("/proc")
	if err != nil {
		log.Fatalln("read /proc err:", err)
		return nil
	}

	size := len(dirs)
	if size == 0 {
		log.Fatalln("dirs is null")
		return nil
	}

	for i := 0; i < size; i++ {
		_, e := strconv.Atoi(dirs[i])
		if e != nil {
			continue
		}
		num += 1
	}

	return []*dataobj.MetricValue{
		GaugeValue("sys.ps.process.total", num),
	}
}
func EntityNumMetrics() []*model.MetricValue {
	data, err := file.ToTrimString("/proc/loadavg")
	if err != nil {
		return nil
	}

	L := strings.Fields(data)
	if len(L) < 5 {
		log.Fatalln("get entity num err: %v", data)
		return nil
	}

	arr := strings.Split(L[3], "/")
	if len(arr) != 2 {
		log.Fatalln("get entity num err: %v", data)
		return nil
	}

	num, err := strconv.ParseFloat(arr[1], 64)
	if err != nil {
		log.Fatalln("get entity num err: %v", err)
		return nil
	}

	return []*model.MetricValue{
		GaugeValue("sys.ps.entity.total", num),
	}
}
