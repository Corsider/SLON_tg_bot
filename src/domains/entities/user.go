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

func (t *TargetUser) ToFlatUser() string {
	flatUser := "Юзер: " + t.Target + "\nТип расписания: "
	switch t.Schedule {
	case ScheduleType_ByMessage:
		flatUser += "Случайный триггер на сообщение"
	case ScheduleType_Every3Hours:
		flatUser += "Триггер каждые 3 часа"
	case ScheduleType_Random:
		flatUser += "Случайный триггер на случайного юзера с этим типом"
	}
	tags := "*нет тегов*"
	if t.Tags != nil {
		tags = *t.Tags
	}
	flatUser += "\nТеги: " + tags
	return flatUser
}
