package slack

import (
	"crypto/hmac"
	"crypto/sha256"
	"cubeflow/pkg/config"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type SlackRequestBody struct {
	Token               string `form:"token"`
	TeamID              string `form:"team_id"`
	TeamDomain          string `form:"team_domain"`
	ChannelID           string `form:"channel_id"`
	ChannelName         string `form:"channel_name"`
	UserID              string `form:"user_id"`
	UserName            string `form:"user_name"`
	Command             string `form:"command"`
	Text                string `form:"text"`
	ApiAppID            string `form:"api_app_id"`
	IsEnterpriseInstall string `form:"is_enterprise_install"`
	ResponseURL         string `form:"response_url"`
	TriggerID           string `form:"trigger_id"`
}

func AuthChecker(c *gin.Context) (bool, string, string, string) {
	var requestBody SlackRequestBody
	if err := c.ShouldBind(&requestBody); err != nil {
		log.Printf("Authchecker bind JSON error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"Authchecker bind JSON error": err.Error()})
		return false, "_", "_", "_"
	}

	slackSignature := c.Request.Header.Get("x-slack-signature")
	slackTimeStamp := c.Request.Header.Get("x-slack-request-timestamp")
	singleLineSlackBody := "v0" + ":" + slackTimeStamp + ":" + CreateSingleLineBody(requestBody)
	key, _ := ComputeHMACSHA256Hex(config.Variable.Slack.SigningSecret, singleLineSlackBody)
	slackPromptInput := url.QueryEscape(requestBody.Text)
	triggerUserName := url.QueryEscape(requestBody.UserName)
	slackChannel := url.QueryEscape(requestBody.ChannelName)
	return slackSignature == "v0="+key, slackPromptInput, triggerUserName, slackChannel
}

func CreateSingleLineBody(body SlackRequestBody) string {
	var keyValuePairs []string

	keyValuePairs = append(keyValuePairs, fmt.Sprintf("token=%s", url.QueryEscape(body.Token)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("team_id=%s", url.QueryEscape(body.TeamID)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("team_domain=%s", url.QueryEscape(body.TeamDomain)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("channel_id=%s", url.QueryEscape(body.ChannelID)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("channel_name=%s", url.QueryEscape(body.ChannelName)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("user_id=%s", url.QueryEscape(body.UserID)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("user_name=%s", url.QueryEscape(body.UserName)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("command=%s", url.QueryEscape(body.Command)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("text=%s", url.QueryEscape(body.Text)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("api_app_id=%s", url.QueryEscape(body.ApiAppID)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("is_enterprise_install=%s", url.QueryEscape(body.IsEnterpriseInstall)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("response_url=%s", url.QueryEscape(body.ResponseURL)))
	keyValuePairs = append(keyValuePairs, fmt.Sprintf("trigger_id=%s", url.QueryEscape(body.TriggerID)))
	return strings.Join(keyValuePairs, "&")
}

func ComputeHMACSHA256Hex(key, data string) (string, error) {
	keyBytes := []byte(key)
	dataBytes := []byte(data)
	h := hmac.New(sha256.New, keyBytes)

	_, err := h.Write(dataBytes)
	if err != nil {
		return "Error when create crypto key", err
	}

	hash := h.Sum(nil)
	hexHash := hex.EncodeToString(hash)
	return hexHash, nil
}
