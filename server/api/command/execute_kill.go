package command

import (
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ExecuteKillRequest struct {
	OwnerType string `uri:"owner_type" binding:"required"`
	OwnerName string `uri:"owner_name" binding:"required"`
	Pid       string `uri:"pid"`
	Operation string `uri:"operation"`
	Step      string `uri:"step"`
	Time      string `uri:"time"`
}

func ExecuteKill(c *gin.Context) {
	var req ExecuteKillRequest
	c.ShouldBindUri(&req)

	pid := req.Pid
	if pid == "" {
		pid = req.Operation + "/" + req.Step + "/" + req.Time
	}

	process := runningProcesses[pid]
	if process == nil {
		common.HandleError(c, http.StatusNotFound, "process doesnot exist or is already killed. "+pid, nil)
		return
	}

	defer delete(runningProcesses, pid)

	if err := process.Kill(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to kill process "+pid, err)
		return
	}

	go func() {
		time.Sleep(time.Duration(1) * time.Second)

		logFilePath := constants.GET_DATA_DIR() + "/" + req.OwnerType + "/" + req.OwnerName + "/history/" + pid + "/execute.log"
		logrus.Trace("logFilePath: ", logFilePath)
		logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logrus.Warn("cannot write logFile ", logFilePath, err)
			debug.PrintStack()
		} else {
			logFile.WriteString("\n\n")
			logFile.WriteString("\033[31m\033[01m\033[05m ###############" + "######################" + "###############\033[0m \n")
			logFile.WriteString("\033[31m\033[01m\033[05m ###############" + " You killed the task. " + "###############\033[0m \n")
			logFile.WriteString("\033[31m\033[01m\033[05m ###############" + " 您强制结束了该任务。 " + "###############\033[0m \n")
			logFile.WriteString("\033[31m\033[01m\033[05m ###############" + "######################" + "###############\033[0m \n")
			logFile.WriteString("\n\n")
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
	})
}
