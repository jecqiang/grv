package main

import (
	log "github.com/Sirupsen/logrus"
	"sync"
)

const (
	PROMPT_TEXT = ":"
)

type PropertyValue struct {
	Property string
	Value    string
}

type StatusBarView struct {
	rootView RootView
	repoData RepoData
	channels *Channels
	config   ConfigSetter
	active   bool
	lock     sync.Mutex
}

func NewStatusBarView(rootView RootView, repoData RepoData, channels *Channels, config ConfigSetter) *StatusBarView {
	return &StatusBarView{
		rootView: rootView,
		repoData: repoData,
		channels: channels,
		config:   config,
	}
}

func (statusBarView *StatusBarView) Initialise() (err error) {
	return
}

func (statusBarView *StatusBarView) HandleKeyPress(keystring string) (err error) {
	return
}

func (statusBarView *StatusBarView) HandleAction(action Action) (err error) {
	switch action.ActionType {
	case ACTION_PROMPT:
		input := Prompt(PROMPT_TEXT)
		errors := statusBarView.config.Evaluate(input)
		statusBarView.channels.ReportErrors(errors)
	}

	return
}

func (statusBarView *StatusBarView) OnActiveChange(active bool) {
	statusBarView.lock.Lock()
	defer statusBarView.lock.Unlock()

	log.Debugf("StatusBarView active: %v", active)
	statusBarView.active = active

	return
}

func (statusBarView *StatusBarView) ViewId() ViewId {
	return VIEW_STATUS_BAR
}

func (statusBarView *StatusBarView) Render(win RenderWindow) (err error) {
	statusBarView.lock.Lock()
	defer statusBarView.lock.Unlock()

	lineBuilder, err := win.LineBuilder(0, 1)
	if err != nil {
		return
	}

	if statusBarView.active {
		promptText, promptPoint := PromptState()
		lineBuilder.Append("%v", promptText)
		win.SetCursor(0, uint(promptPoint+len(PROMPT_TEXT)))
	} else {
		lineBuilder.Append(" %v", statusBarView.repoData.Path())

		/*lineBuilder.Append(" ")
		viewHierarchy := statusBarView.rootView.ActiveViewHierarchy()

		for _, view := range viewHierarchy {
			if err = view.RenderStatusBar(lineBuilder); err != nil {
				return
			}
		}*/

		win.ApplyStyle(CMP_STATUSBARVIEW_NORMAL)
	}

	return
}

func (statusBarView *StatusBarView) RenderStatusBar(lineBuilder *LineBuilder) (err error) {
	return
}

func (statusBarView *StatusBarView) RenderHelpBar(lineBuilder *LineBuilder) (err error) {
	lineBuilder.AppendWithStyle(CMP_HELPBARVIEW_SPECIAL, "Enter a command")

	return
}

func RenderStatusProperties(lineBuilder *LineBuilder, propertyValues []PropertyValue) {
	for _, propValue := range propertyValues {
		lineBuilder.Append("%v: %v     ", propValue.Property, propValue.Value)
	}
}
