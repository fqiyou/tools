/**
 *  Author: SongLee24
 *  Email: lisong.shine@qq.com
 *  Date: 2018-08-15
 *
 *
 *  prometheus.Desc是指标的描述符，用于实现对指标的管理
 *
 */

package collector

import (
	"github.com/fqiyou/tools/foo/system"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

// 指标结构体
type Metrics struct {
	metrics map[string]*prometheus.Desc
	mutex   sync.Mutex
}

/**
 * 函数：newGlobalMetric
 * 功能：创建指标描述符
 */
func newGlobalMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}


/**
 * 工厂方法：NewMetrics
 * 功能：初始化指标信息，即Metrics结构体
 */
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"disk_raid_monitor_raid_devices": newGlobalMetric(namespace, "disk_raid_monitor_raid_devices",
				"The description of disk_raid_monitor_raid_devices",[]string{"hostname","raidname","raidlevel","version"}),
			"disk_raid_monitor_total_devices": newGlobalMetric(namespace, "disk_raid_monitor_total_devices",
				"The description of disk_raid_monitor_total_devices",[]string{"hostname","raidname","raidlevel","version"}),
			"disk_raid_monitor_active_devices": newGlobalMetric(namespace, "disk_raid_monitor_active_devices",
				"The description of disk_raid_monitor_active_devices",[]string{"hostname","raidname","raidlevel","version"}),
			"disk_raid_monitor_working_devices": newGlobalMetric(namespace, "disk_raid_monitor_working_devices",
				"The description of disk_raid_monitor_working_devices",[]string{"hostname","raidname","raidlevel","version"}),
			"disk_raid_monitor_failed_devices": newGlobalMetric(namespace, "disk_raid_monitor_failed_devices",
				"The description of disk_raid_monitor_failed_devices",[]string{"hostname","raidname","raidlevel","version"}),
			"disk_raid_monitor_spare_devices": newGlobalMetric(namespace, "disk_raid_monitor_spare_devices",
				"The description of disk_raid_monitor_spare_devices",[]string{"hostname","raidname","raidlevel","version"}),
			"disk_raid_test": newGlobalMetric(namespace, "disk_raid_test",
				"The description of disk_raid_test",[]string{"hostname"}),

		},

	}
}

/**
 * 接口：Describe
 * 功能：传递结构体中的指标描述符到channel
 */
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}


/**
 * 接口：Collect
 * 功能：抓取最新的数据，传递给channel
 */
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()  // 加锁
	defer c.mutex.Unlock()

	//mockCounterMetricData, mockGaugeMetricData := c.GenerateMockData()

	info := system.RaidDisk{}
	info.Collect()
	for _,v := range info.RaidDiskSystemInfo{

		ch <-prometheus.MustNewConstMetric(c.metrics["disk_raid_monitor_raid_devices"], prometheus.GaugeValue, float64(v.RaidDevices),
			info.HostName,
			v.RaidName,v.RaidLevel,v.Version,
		)

		ch <-prometheus.MustNewConstMetric(c.metrics["disk_raid_monitor_total_devices"], prometheus.GaugeValue, float64(v.TotalDevices),
			info.HostName,
			v.RaidName,v.RaidLevel,v.Version,
		)

		ch <-prometheus.MustNewConstMetric(c.metrics["disk_raid_monitor_active_devices"], prometheus.GaugeValue, float64(v.ActiveDevices),
			info.HostName,
			v.RaidName,v.RaidLevel,v.Version,
		)

		ch <-prometheus.MustNewConstMetric(c.metrics["disk_raid_monitor_working_devices"], prometheus.GaugeValue, float64(v.WorkingDevices),
			info.HostName,
			v.RaidName,v.RaidLevel,v.Version,
		)

		ch <-prometheus.MustNewConstMetric(c.metrics["disk_raid_monitor_failed_devices"], prometheus.GaugeValue, float64(v.FailedDevices),
			info.HostName,
			v.RaidName,v.RaidLevel,v.Version,
		)

		ch <-prometheus.MustNewConstMetric(c.metrics["disk_raid_monitor_spare_devices"], prometheus.GaugeValue, float64(v.SpareDevices),
			info.HostName,
			v.RaidName,v.RaidLevel,v.Version,
		)


	}
	ch <-prometheus.MustNewConstMetric(c.metrics["disk_raid_test"], prometheus.GaugeValue, float64(111),
		info.HostName,
	)



}

