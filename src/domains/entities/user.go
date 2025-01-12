package entities

import (
	"github.com/lib/pq"
	"strings"
)

type ScheduleType int
type TagType int

const (
	ScheduleType_ByMessage   ScheduleType = 0
	ScheduleType_Every3Hours ScheduleType = 1
	ScheduleType_Random      ScheduleType = 2
)

const (
	TagType_INSULT    TagType = 0
	TagType_OBSCENITY TagType = 1
	TagType_THREAT    TagType = 2
)

type TargetUser struct {
	CreatorID int64         `db:"creator_id"`
	Target    string        `db:"target"`
	Schedule  ScheduleType  `db:"schedule"`
	Tags      pq.Int32Array `db:"tags"`
}

func NewDefaultUser(creatorId int64, target string) *TargetUser {
	return &TargetUser{
		CreatorID: creatorId,
		Target:    target,
		Schedule:  ScheduleType_ByMessage,
		Tags:      pq.Int32Array{int32(TagType_INSULT)},
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
	var tags string
	for _, tg := range t.Tags {
		switch tg {
		case int32(TagType_INSULT):
			tags += "ОСКОРБЛЕНИЕ, "
		case int32(TagType_OBSCENITY):
			tags += "НЕПРИСТОЙНОСТЬ, "
		case int32(TagType_THREAT):
			tags += "УГРОЗА, "
		}
	}
	tags = strings.TrimSuffix(tags, ", ")
	tags = strings.TrimSuffix(tags, ",")
	flatUser += "\nТеги: " + tags
	return flatUser
}

func (t *TargetUser) ToFlatTags() string {
	var tags string
	for _, tg := range t.Tags {
		switch tg {
		case int32(TagType_INSULT):
			tags += "ОСКОРБЛЕНИЕ, "
		case int32(TagType_OBSCENITY):
			tags += "НЕПРИСТОЙНОСТЬ, "
		case int32(TagType_THREAT):
			tags += "УГРОЗА, "
		}
	}
	tags = strings.TrimSuffix(tags, ", ")
	tags = strings.TrimSuffix(tags, ",")
	return tags
}

func (t *TargetUser) GetTags() []TagType {
	ret := []TagType{}
	for _, tg := range t.Tags {
		ret = append(ret, TagType(tg))
	}
	return ret
}
