package services

import (
	"encoding/json"
	"fmt"
	db "gin-web-demo/db"
	model "gin-web-demo/models"
	uuid "github.com/satori/go.uuid"
)

func SubmitJob(jobName, cwd, cmds string) (string, error) {
	var err error
	jobid := uuid.NewV1().String()
	fmt.Println("jobid:" + jobid)
	job := model.LSFJobReq{
		JobID:   jobid,
		JobName: jobName,
		WorkDir: cwd,
		Command: cmds,
	}
	db.InsertJob(&job)
	datainfo, err := json.Marshal(job)
	return string(datainfo), err
}
