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
	shared.RegisterHandler(
		"track moved", func(data ...any) {
			go result.doBackup("move track")
		},
	)
	shared.RegisterHandler(
		"track deleted", func(data ...any) {
			go result.doBackup("delete track")
		},
	)
	shared.RegisterHandler(
		"track upserted", func(data ...any) {
			go result.doBackup("upsert track")
		},
	)
	shared.RegisterHandler(
		"journal entry changed", func(data ...any) {
			go result.doBackup("change journal entry")
		},
	)
	shared.RegisterHandler(
		"settings changed", func(data ...any) {
			go result.doBackup("change settings")
		},
	)
	shared.RegisterHandler(
		"git enablement changed", func(data ...any) {
			result.enabled, _ = data[0].(bool)
		},
	)
	shared.RegisterHandler(
		"git push changed", func(data ...any) {
			result.push, _ = data[0].(bool)
		},
	)
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
