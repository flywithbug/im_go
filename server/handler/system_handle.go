package handler



import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"net/http"
	"time"
	"fmt"
	"im_go/model"
	"strconv"
	"github.com/gin-gonic/gin"
)


// 系统状态信息
func handleSystem(c *gin.Context) {
	mem, _ := mem.VirtualMemory()
	cpuNum, _ := cpu.Counts(true);
	cpuInfo, _ := cpu.Percent(10 * time.Microsecond, true);

	data := make(map[string]interface{})
	//data["im.conn"] = len(ClientMaps)
	data["mem.total"] = fmt.Sprintf("%vMB", mem.Total/1024/1024)
	data["mem.free"] = fmt.Sprintf("%vMB", mem.Free/1024/1024)
	data["mem.used_percent"] = fmt.Sprintf("%s%%", strconv.FormatFloat(mem.UsedPercent, 'f', 2, 64))
	data["cpu.num"] = cpuNum
	data["cpu.info"] = cpuInfo

	c.IndentedJSON(http.StatusOK,model.NewIMResponseData(data,""))
}
