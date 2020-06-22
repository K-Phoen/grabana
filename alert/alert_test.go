package alert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAlertCanBeCreated(t *testing.T) {
	req := require.New(t)

	a := New("some alert")

	req.Equal("some alert", a.Builder.Name)
	req.Equal(string(LastState), a.Builder.ExecutionErrorState)
	req.Equal(string(KeepLastState), a.Builder.NoDataState)
}

func TestMessageCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", Message("content"))

	req.Equal("content", a.Builder.Message)
}

func TestNotificationCanBeSet(t *testing.T) {
	req := require.New(t)
	channel := &Channel{ID: 1, UID: "channel"}

	a := New("", Notify(channel))

	req.Len(a.Builder.Notifications, 1)
	req.Equal("channel", a.Builder.Notifications[0].UID)
	req.Empty(a.Builder.Notifications[0].ID)
}

func TestNotificationCanBeSetInBulk(t *testing.T) {
	req := require.New(t)
	channel := &Channel{ID: 1, UID: "channel"}
	otherChannel := &Channel{ID: 2, UID: "other-channel"}

	a := New("", NotifyChannels(channel, otherChannel))

	req.Len(a.Builder.Notifications, 2)
	req.ElementsMatch([]string{"channel", "other-channel"}, []string{a.Builder.Notifications[0].UID, a.Builder.Notifications[1].UID})
	req.Empty(a.Builder.Notifications[0].ID)
	req.Empty(a.Builder.Notifications[1].ID)
}

func TestNotificationCanBeSetByChannelID(t *testing.T) {
	req := require.New(t)

	a := New("", NotifyChannel("P-N3fxuZz"))

	req.Len(a.Builder.Notifications, 1)
	req.Equal("P-N3fxuZz", a.Builder.Notifications[0].UID)
	req.Empty(a.Builder.Notifications[0].ID)
}

func TestForIntervalCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", For("1m"))

	req.Equal("1m", a.Builder.For)
}

func TestFrequencyCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", EvaluateEvery("1m"))

	req.Equal("1m", a.Builder.Frequency)
}

func TestErrorModeCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", OnExecutionError(Alerting))

	req.Equal(string(Alerting), a.Builder.ExecutionErrorState)
}

func TestNoDataModeCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", OnNoData(OK))

	req.Equal(string(OK), a.Builder.NoDataState)
}

func TestConditionsCanBeSet(t *testing.T) {
	req := require.New(t)

	a := New("", If(And))

	req.Len(a.Builder.Conditions, 1)
	req.Equal(string(And), a.Builder.Conditions[0].Operator.Type)
}
