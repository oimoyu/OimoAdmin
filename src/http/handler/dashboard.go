package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"time"
)

func Logout(c *gin.Context) {
	c.SetCookie("admin_token", "", -3600, "/", "", false, true)
	g.AdminLoginData.IP = ""
	g.AdminLoginData.Token = ""
	restful.Ok(c, "Logout Success")
}

func getSystemInfo() map[string]interface{} {
	cpuUsage, _ := cpu.Percent(0, false)
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")

	cpuInfo, _ := cpu.Info()

	cpuPercentage := fmt.Sprintf("%.1f", cpuUsage[0])

	memoryTotalMB := float64(memInfo.Total) / (1024 * 1024)
	memoryUsageMB := float64(memInfo.Used) / (1024 * 1024)

	memoryTotalMBString := fmt.Sprintf("%.0f", memoryTotalMB)
	memoryUsageMBString := fmt.Sprintf("%.0f", memoryUsageMB)
	memoryPercentageString := fmt.Sprintf("%.1f", memoryUsageMB/memoryTotalMB*100)

	diskTotalMB := float64(diskInfo.Total) / (1024 * 1024)
	diskUsageMB := float64(diskInfo.Used) / (1024 * 1024)

	diskTotalMBString := fmt.Sprintf("%.0f", diskTotalMB)
	diskUsageMBString := fmt.Sprintf("%.0f", diskUsageMB)
	diskPercentageString := fmt.Sprintf("%.1f", diskUsageMB/diskTotalMB*100)

	uploadSpeed, downloadSpeed, _ := getNetworkSpeed(time.Millisecond * 3000)
	uploadSpeedString := fmt.Sprintf("%.1f", uploadSpeed)
	downloadSpeedString := fmt.Sprintf("%.1f", downloadSpeed)

	return map[string]interface{}{
		"cpu_percentage_usage":    cpuPercentage,
		"cpu_cores":               cpuInfo[0].Cores,
		"memory_total":            memoryTotalMBString,
		"memory_usage":            memoryUsageMBString,
		"memory_percentage_usage": memoryPercentageString,
		"disk_total":              diskTotalMBString,
		"disk_usage":              diskUsageMBString,
		"disk_percentage_usage":   diskPercentageString,
		"upload_speed":            uploadSpeedString,
		"download_speed":          downloadSpeedString,
	}
}
func getNetworkSpeed(interval time.Duration) (float64, float64, error) {
	// TODO: consider using websocket to get accurate period speed

	// IOCounters return the accumulate data count from system start
	ioCounters, err := net.IOCounters(false)
	if err != nil {
		return 0, 0, err
	}

	if len(ioCounters) == 0 {
		return 0, 0, fmt.Errorf("no network interfaces found")
	}

	// init data
	initialBytesSent := ioCounters[0].BytesSent
	initialBytesRecv := ioCounters[0].BytesRecv

	// wait interval
	time.Sleep(interval)

	// new data
	ioCounters, err = net.IOCounters(false)
	if err != nil {
		return 0, 0, err
	}

	bytesSent := ioCounters[0].BytesSent - initialBytesSent
	bytesRecv := ioCounters[0].BytesRecv - initialBytesRecv

	uploadSpeed := float64(bytesSent) / interval.Seconds() / (1024 * 1024)
	downloadSpeed := float64(bytesRecv) / interval.Seconds() / (1024 * 1024)

	return uploadSpeed, downloadSpeed, nil
}

func SysInfo(c *gin.Context) {
	sysInfo := getSystemInfo()
	restful.Ok(c, sysInfo)
}

func Dashboard(c *gin.Context) {
	returnData := map[string]interface{}{
		"db_name":     g.OimoAdmin.DBName,
		"driver_name": g.OimoAdmin.DB.Dialector.Name(),
	}
	restful.Ok(c, returnData)
}
