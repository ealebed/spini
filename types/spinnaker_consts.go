package types

const (
	spinnakerRoleGroupDevops = "devops"

	slackPipelineNotifyDeployText = "*Deploy* ${ trigger['artifacts'].?[type == 'docker/image'].![reference] }\n*User:* ${ trigger['user'] }"
)
