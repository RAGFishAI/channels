package channels

import (
	"errors"
	"os"
)

var (
	ErrorAuth       = errors.New("auth enabled not initialised")
	ErrorPermission = errors.New("permissions enabled not initialised")
	ErrorChannelId  = errors.New("invalid channelid")
	Empty           string
	TenantId        = os.Getenv("Tenant_ID")
)

func TruncateDescription(description string, limit int) string {
	if len(description) <= limit {
		return description
	}

	truncated := description[:limit] + "..."
	return truncated
}

func AuthandPermission(channel *Channel) error {

	//check auth enable if enabled, use auth pkg otherwise it will return error
	if channel.AuthEnable && !channel.Auth.AuthFlg {

		return ErrorAuth
	}
	//check permission enable if enabled, use team-role pkg otherwise it will return error
	if channel.PermissionEnable && !channel.Auth.PermissionFlg {

		return ErrorPermission

	}

	return nil
}
