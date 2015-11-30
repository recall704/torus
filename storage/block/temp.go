package block

import (
	"sync"

	"golang.org/x/net/context"

	"github.com/barakmich/agro"
)

var _ agro.BlockStore = &tempBlockStore{}

func init() {
	agro.RegisterBlockStore("temp", openTempBlockStore)
}

type tempBlockStore struct {
	mut     sync.RWMutex
	store   map[agro.BlockID][]byte
	nBlocks uint64
}

func openTempBlockStore(cfg agro.Config, gmd agro.GlobalMetadata) (agro.BlockStore, error) {
	return &tempBlockStore{
		store: make(map[agro.BlockID][]byte),
		// Lie about the number of blocks.
		nBlocks: cfg.StorageSize / 1024,
	}, nil
}

func (t *tempBlockStore) Flush() error { return nil }

func (t *tempBlockStore) Close() error {
	t.mut.Lock()
	t.store = nil
	t.mut.Unlock()
	return nil
}

func (t *tempBlockStore) NumBlocks() uint64 {
	return t.nBlocks
}

func (t *tempBlockStore) GetBlock(_ context.Context, s agro.BlockID) ([]byte, error) {
	t.mut.RLock()
	defer t.mut.RUnlock()

	if t.store == nil {
		return nil, agro.ErrClosed
	}

	x, ok := t.store[s]
	if !ok {
		return nil, agro.ErrBlockNotExist
	}
	return x, nil
}

func (t *tempBlockStore) WriteBlock(_ context.Context, s agro.BlockID, data []byte) error {
	t.mut.Lock()
	defer t.mut.Unlock()

	if t.store == nil {
		return agro.ErrClosed
	}

	t.store[s] = data
	return nil
}

func (t *tempBlockStore) DeleteBlock(_ context.Context, s agro.BlockID) error {
	t.mut.Lock()
	defer t.mut.Unlock()

	if t.store == nil {
		return agro.ErrClosed
	}

	delete(t.store, s)
	return nil
}
