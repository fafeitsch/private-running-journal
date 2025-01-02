package backup

import (
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os/exec"
)

type Backup struct {
	baseDirectory string
	enabled       bool
	push          bool
}

func Init(baseDirectory string, enabled bool, push bool) *Backup {
	result := &Backup{baseDirectory: baseDirectory, enabled: enabled, push: push}
	shared.Listen(shared.TrackDeletedEvent{}, func(event shared.TrackDeletedEvent) {
		go result.doBackup("delete track")
	})
	shared.Listen(shared.TrackUpsertedEvent{}, func(event shared.TrackUpsertedEvent) {
		go result.doBackup("upsert track")
	})
	shared.Listen(shared.JournalEntryUpsertedEvent{}, func(event shared.JournalEntryUpsertedEvent) {
		go result.doBackup("change journal entry")
	})
	shared.Listen(shared.JournalEntryDeletedEvent{}, func(event shared.JournalEntryDeletedEvent) {
		go result.doBackup("delete journal entry")
	})
	shared.Listen(shared.SettingsChangedEvent{}, func(event shared.SettingsChangedEvent) {
		go result.doBackup("change settings")
	})
	shared.Listen(shared.GitEnablementChangedEvent{}, func(event shared.GitEnablementChangedEvent) {
		result.enabled = event.NewValue
	})
	shared.Listen(shared.GitPushChangedEvent{}, func(event shared.GitPushChangedEvent) {
		result.push = event.NewValue
	})
	return result
}

func (b *Backup) doBackup(message string) {
	if !b.enabled {
		return
	}
	out, err := exec.Command("git", "-C", b.baseDirectory, "add", "--all").CombinedOutput()
	log.Print(string(out))
	if err != nil {
		runtime.EventsEmit(shared.Context, "git-error", string(out))
		log.Printf("Failed to execute git add command: %v", err)
		return
	}
	out, err = exec.Command("git", "-C", b.baseDirectory, "commit", "-m", message).CombinedOutput()
	log.Print(string(out))
	if err != nil {
		runtime.EventsEmit(shared.Context, "git-error", string(out))
		log.Printf("Failed to execute git commit command: %v", err)
		return
	}
	if !b.push {
		return
	}
	out, err = exec.Command("git", "-C", b.baseDirectory, "push").CombinedOutput()
	log.Print(string(out))
	if err != nil {
		runtime.EventsEmit(shared.Context, "git-error", string(out))
		log.Printf("Failed to execute git push command: %v", err)
		return
	}
}

func (b *Backup) Pull() error {
	out, err := exec.Command("git", "-C", b.baseDirectory, "pull").CombinedOutput()
	log.Print(string(out))
	if err != nil {
		runtime.EventsEmit(shared.Context, "git-error", string(out))
		return err
	}
	return nil
}
