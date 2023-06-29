package global

import "fmt"

const (
	TopicGroupActive = "TopicGroupActive" //群激活的topic
)

func GetGroupActiveTopic(gid uint64) string {
	return fmt.Sprintf("TopicGroupActive:%d", gid)
}
