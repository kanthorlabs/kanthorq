package tester

import (
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
)

func FakeEvents(subject string, count int) []*entities.Event {
	events := make([]*entities.Event, count)
	for i := 0; i < count; i++ {
		events[i] = entities.NewEvent(subject, xfaker.DataOf16Kb())
	}
	return events
}

func FakeTasks(events []*entities.Event, state entities.TaskState) []*entities.Task {
	tasks := make([]*entities.Task, len(events))
	for i := range events {
		tasks[i] = &entities.Task{
			EventId:        events[i].Id,
			Subject:        events[i].Subject,
			State:          state,
			ScheduleAt:     time.Now().UnixMilli(),
			AttemptCount:   1,
			AttemptedAt:    time.Now().UnixMilli(),
			AttemptedError: []entities.AttemptedError{},
			FinalizedAt:    0,
			Metadata:       events[i].Metadata,
			CreatedAt:      events[i].CreatedAt,
			UpdatedAt:      time.Now().UnixMilli(),
		}
	}
	return tasks
}

func FakeTasksWithSchedule(events []*entities.Event, state entities.TaskState, schedule time.Time) []*entities.Task {
	tasks := make([]*entities.Task, len(events))
	for i := range events {
		tasks[i] = &entities.Task{
			EventId:        events[i].Id,
			Subject:        events[i].Subject,
			State:          state,
			ScheduleAt:     schedule.UnixMilli(),
			AttemptCount:   1,
			AttemptedAt:    time.Now().UnixMilli(),
			AttemptedError: []entities.AttemptedError{},
			FinalizedAt:    0,
			Metadata:       events[i].Metadata,
			CreatedAt:      events[i].CreatedAt,
			UpdatedAt:      time.Now().UnixMilli(),
		}
	}
	return tasks
}
