package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/db"
	"github.com/evergreen-ci/evergreen/model/event"
	"github.com/evergreen-ci/evergreen/model/task"
	"github.com/evergreen-ci/evergreen/plugin"
	"github.com/evergreen-ci/evergreen/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tychoish/grip"
	"github.com/tychoish/grip/message"
)

func TestResourceInfoEndPoints(t *testing.T) {
	testConfig := evergreen.TestConfig()
	testApiServer, err := CreateTestServer(testConfig, nil, plugin.APIPlugins, true)
	defer testApiServer.Close()
	testutil.HandleTestingErr(err, t, "failed to create new API server")

	err = db.ClearCollections(event.TaskCollection, event.SystemCollection)
	testutil.HandleTestingErr(err, t, "problem clearing event collection")
	err = db.ClearCollections(task.Collection)
	testutil.HandleTestingErr(err, t, "problem clearing task collection")

	const (
		url    = "http://localhost:8182/api/2/task/"
		taskId = "the_task_id"
	)

	_, err = insertTaskForTesting(taskId, "version", "project", task.TestResult{})
	testutil.HandleTestingErr(err, t, "problem creating task")

	Convey("For the system info endpoint", t, func() {
		data := message.CollectSystemInfo().(*message.SystemInfo)
		data.Base = message.Base{}
		data.Errors = []string{}
		Convey("the system info endpoint should return 200", func() {
			payload, err := json.Marshal(data)
			So(err, ShouldBeNil)

			request, err := http.NewRequest("POST", url+taskId+"/system_info", bytes.NewBuffer(payload))
			So(err, ShouldBeNil)
			resp, err := http.DefaultClient.Do(request)
			testutil.HandleTestingErr(err, t, "problem making request")
			So(resp.StatusCode, ShouldEqual, 200)
		})

		Convey("the system data should persist in the database", func() {
			events, err := event.FindTask(event.TaskSystemInfoEvents(taskId, 0))
			testutil.HandleTestingErr(err, t, "problem finding task event")
			So(len(events), ShouldEqual, 1)
			e := events[0]
			So(e.ResourceId, ShouldEqual, taskId)
			taskData, ok := e.Data.Data.(*event.TaskSystemResourceData)
			So(ok, ShouldBeTrue)
			grip.Info(taskData.SystemInfo)
		})
	})

	Convey("For the process info endpoint", t, func() {
		data := message.CollectProcessInfoSelfWithChildren()
		Convey("the process info endpoint should return 200", func() {
			payload, err := json.Marshal(data)
			So(err, ShouldBeNil)

			request, err := http.NewRequest("POST", url+taskId+"/process_info", bytes.NewBuffer(payload))
			resp, err := http.DefaultClient.Do(request)
			testutil.HandleTestingErr(err, t, "problem making request")
			So(resp.StatusCode, ShouldEqual, 200)
		})

		Convey("the process data should persist in the database", func() {
			events, err := event.FindTask(event.TaskProcessInfoEvents(taskId, 0))
			testutil.HandleTestingErr(err, t, "problem finding task event")
			So(len(events), ShouldEqual, 1)
			e := events[0]
			So(e.ResourceId, ShouldEqual, taskId)
			taskData, ok := e.Data.Data.(*event.TaskProcessResourceData)
			So(ok, ShouldBeTrue)
			grip.Info(taskData.Processes)
		})
	})
}
