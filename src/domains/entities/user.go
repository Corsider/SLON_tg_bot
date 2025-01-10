package entities

type ScheduleType int

const (
	ScheduleType_ByMessage   ScheduleType = 0
	ScheduleType_Every3Hours ScheduleType = 1
	ScheduleType_Random      ScheduleType = 2
)

type TargetUser struct {
	CreatorID int64        `db:"creator_id"`
	Target    string       `db:"target"`
	Schedule  ScheduleType `db:"schedule"`
	Tags      *string      `db:"tags"` // [курсовая, диплом]
}

func NewDefaultUser(creatorId int64, target string) *TargetUser {
	return &TargetUser{
		CreatorID: creatorId,
		Target:    target,
		Schedule:  ScheduleType_ByMessage,
		Tags:      nil,
	}
}
