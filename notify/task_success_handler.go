package notify

import (
	"10gen.com/mci"
	"10gen.com/mci/model"
	"10gen.com/mci/web"
	"github.com/10gen-labs/slogger/v1"
)

// Handler for notifying users that a task succeeded. Implements
// NotificationHandler from notification_handler.go.
type TaskSuccessHandler struct {
	TaskNotificationHandler
	Name string
}

func (self *TaskSuccessHandler) GetNotifications(ae *web.App, configName string,
	key *NotificationKey) ([]Email, error) {
	var emails []Email
	preface := mciSuccessPreface
	if key.NotificationRequester == mci.PatchVersionRequester {
		preface = patchSuccessPreface
	}
	triggeredNotifications, err :=
		self.getRecentlyFinishedTasksWithStatus(key, mci.TaskSucceeded, preface, successSubject)

	if err != nil {
		return nil, err
	}

	for _, triggered := range triggeredNotifications {
		email, err := self.TemplateNotification(ae, configName, &triggered)
		if err != nil {
			mci.Logger.Logf(slogger.WARN, "Error templating notification for task `%v`: %v",
				triggered.Current.Id, err)
			continue
		}

		emails = append(emails, email)
	}

	return emails, nil
}

func (self *TaskSuccessHandler) TemplateNotification(ae *web.App, configName string,
	notification *TriggeredTaskNotification) (Email, error) {
	changeInfo, err := self.GetChangeInfo(notification)
	if err != nil {
		return nil, err
	}
	return self.templateNotification(ae, configName, notification, changeInfo)
}

func (self *TaskSuccessHandler) GetChangeInfo(
	notification *TriggeredTaskNotification) ([]ChangeInfo, error) {
	return self.constructChangeInfo([]model.Task{*notification.Current}, &notification.Key)
}
