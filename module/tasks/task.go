package tasks

import (
	"github.com/zihang5127/easy-operation/conf"
	"github.com/zihang5127/easy-operation/model"
	wsChannel "github.com/zihang5127/easy-operation/module/channel"
	"github.com/zihang5127/easy-operation/module/client"
	"github.com/zihang5127/easy-operation/module/queue"

	"net/url"
	"strconv"
	"strings"
	_ "time"

	"github.com/astaxie/beego/logs"
)

var (
	queues = queue.NewQueue(conf.QueueSize())
)

type Task struct {
	ProjectId int
	ServerId  int
}

func Add(task Task) {
	name := strconv.Itoa(task.ProjectId) + "-" + strconv.Itoa(task.ServerId)
	queues.Enqueue(name, task)
}

func Handle(value interface{}) {
	chann := wsChannel.GetChannel()
	oc := wsChannel.GetOverChannel()
	if task, ok := value.(Task); ok {
		server := model.NewServer()
		server.Id = task.ServerId
		if err := server.Find(); err != nil {
			logs.Error("%s", err.Error())
			return
		}

		project := model.NewProject()
		project.Id = task.ProjectId

		if err := project.Find(); err != nil {
			logs.Error("%s", err.Error())
			return
		}
		if strings.TrimSpace(project.Shell) == "" {
			logs.Warn("", "Shell command does not exist.")
			return
		}

		channel := make(chan []byte, 10)

		client, err := CreateClient()
		//defer close(wsChannel.GetChannel())
		if err != nil {
			logs.Error("%s", err.Error())
			return
		}

		host := server.IpAddress + ":" + strconv.Itoa(server.Port)

		u, err := url.Parse(host)

		scheme := "http"

		if strings.HasPrefix(server.IpAddress, "http://") {
			scheme = "http"
		} else if strings.HasPrefix(server.IpAddress, "https://") {
			scheme = "https"
		} else {
			scheme = "ssh"
		}

		if err != nil {
			logs.Error("%s", err)
			u = &url.URL{Scheme: scheme, Host: host}
		}

		logs.Info("Connecting ... %s", u)

		go client.Command(*u, server.Username, server.PrivateKey, project.Shell, channel)

		if ok {
			for out := range channel {
				chann <- out
			}
			logs.Info("%s", "The command was executed successfully")
			oc <- []byte("quit")
		}

	} else {
		logs.Error("Can not be converted to Task: ", value)
	}
}

func CreateClient() (client.Interface, error) {
	return &client.SSHClient{}, nil
}

func init() {
	queues.Handle = Handle
}
