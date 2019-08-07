package tailer

import (
	"fmt"

	"github.com/hpcloud/tail"
)

type Tailer struct {
	fileName string
	tailConf tail.Config
	tail     *tail.Tail
	started  bool
}

func New(fileName string) *Tailer {
	return &Tailer{
		fileName: fileName,
		tailConf: tail.Config{
			MustExist: true, // Fail early if the file does not exist
			Follow:    true, // Continue looking for new lines (tail -f)
			ReOpen:    true, // Reopen recreated/truncated files (tail -F)
		},
	}
}

func (t *Tailer) Start() (<-chan *tail.Line, error) {
	if t.started {
		return nil, fmt.Errorf("tailer can be started only once")
	}

	t.started = true
	tf, err := tail.TailFile(t.fileName, t.tailConf)
	if err != nil {
		return nil, err
	}
	t.tail = tf
	return tf.Lines, nil
}

func (t *Tailer) Stop() error {
	if !t.started {
		return fmt.Errorf("tailer can be stopped only after start")
	}
	return t.tail.Stop()
}

func (t *Tailer) Wait() error {
	if t.tail != nil {
		return t.tail.Wait()
	}
	return fmt.Errorf("cannot wait if tailer isn't started")
}
