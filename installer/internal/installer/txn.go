package installer

import (
	"os"
	"path/filepath"

	"github.com/imflawlezz/cursor-utils/installer/internal/platform"
)

type txnOp struct {
	path     string
	existed  bool
	previous []byte
	created  bool
	deleted  bool
}

type txn struct {
	ops []txnOp
}

func newTxn() *txn { return &txn{} }

func (t *txn) writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var prev []byte
	existed := false
	if b, err := os.ReadFile(path); err == nil {
		prev = b
		existed = true
	} else if !os.IsNotExist(err) {
		return err
	}
	tmp := path + ".cursor-utils.tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := platform.ReplaceFile(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	t.ops = append(t.ops, txnOp{
		path:     path,
		existed:  existed,
		previous: prev,
		created:  !existed,
	})
	return nil
}

func (t *txn) removeFile(path string) error {
	prev, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	t.ops = append(t.ops, txnOp{
		path:     path,
		existed:  true,
		previous: prev,
		deleted:  true,
	})
	return nil
}

func (t *txn) rollback() {
	for i := len(t.ops) - 1; i >= 0; i-- {
		op := t.ops[i]
		switch {
		case op.deleted:
			_ = os.MkdirAll(filepath.Dir(op.path), 0o755)
			_ = os.WriteFile(op.path, op.previous, 0o644)
		case op.created:
			_ = os.Remove(op.path)
		default:
			if op.existed {
				_ = os.WriteFile(op.path, op.previous, 0o644)
			}
		}
	}
}
