package slack

import (
	"net/url"

	"cubeflow/pkg/config"

	helperFunc "cubeflow/internal/helpers"
)

func IsChannelMatched(channelToCheck string, appsToCheck string) bool {
	splitParameter := ","
	decodedPrompInput, _ := url.QueryUnescape(appsToCheck)

	serviceNameList := helperFunc.MutipleExec(decodedPrompInput, splitParameter)

	for _, channel := range config.Variable.Slack.Channels {
		if channel.Name == channelToCheck && areAllAppsAuthorized(channel, serviceNameList) {
			return true
		}
	}
	return false
}

func IsChannelExist(channelToCheck string) bool {
	for _, channel := range config.Variable.Slack.Channels {
		if channel.Name == channelToCheck {
			return true
		}
	}
	return false
}

func GetChannelNameByAppName(appName string) (string, bool) {
	for _, channel := range config.Variable.Slack.Channels {
		if isAuthorized(channel.RolloutName, appName) {
			return channel.Name, true
		}
	}
	return "", false
}

func areAllAppsAuthorized(channel config.SlackChannels, appsToCheck []string) bool {
	for _, app := range appsToCheck {
		if !isAuthorized(channel.RolloutName, app) && !isAuthorized(channel.ArgoAppName, app) {
			return false
		}
	}
	return true
}

func isAuthorized(channel []string, appName string) bool {
	for _, app := range channel {
		if app == appName {
			return true
		}
	}
	return false
}
