package plex

const (
	// request auth uri
	auth_uri = "/auth/server"

	// auth failed url
	auth_failed_uri = "/auth/failed"

	// auth success url
	auth_success_uri = "/auth/success"

	// heartbeat uri
	heartbeat_uri = "/heartbeat"
)

const (
	// send msg to client (inner server)
	send_msg_uri = "/logic/send/msg"

	// inner server auth uri
	inner_auth_uri = "/inner/auth/server"
)
