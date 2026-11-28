package main

import "context"

// operation 描述一次业务操作；执行时由 Actor 提供自己的状态。
type operation func(activityStore) (ActivityResp, error)

type executor interface {
	ask(context.Context, operation) (ActivityResp, error)
	shutdown()
}

type service struct{ actor executor }

func newService() Service { return &service{actor: newActor()} }

func (s *service) Create(ctx context.Context, id string, quota int) (ActivityResp, error) {
	return s.actor.ask(ctx, func(store activityStore) (ActivityResp, error) { return store.create(id, quota) })
}

func (s *service) Status(ctx context.Context, id string) (ActivityResp, error) {
	return s.actor.ask(ctx, func(store activityStore) (ActivityResp, error) { return store.status(id) })
}

func (s *service) Update(ctx context.Context, id string, quota int) (ActivityResp, error) {
	return s.actor.ask(ctx, func(store activityStore) (ActivityResp, error) { return store.update(id, quota) })
}

func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.actor.ask(ctx, func(store activityStore) (ActivityResp, error) { return ActivityResp{}, store.remove(id) })
	return err
}

func (s *service) Grab(ctx context.Context, id, userID string) (int, error) {
	activity, err := s.actor.ask(ctx, func(store activityStore) (ActivityResp, error) { return store.grab(id, userID) })
	return activity.Remaining, err
}

func (s *service) Shutdown() { s.actor.shutdown() }

// 以下方法已实现业务规则，只能由状态所属的 Actor 串行调用。
type activityStore map[string]*activity

type activity struct {
	id       string
	quota    int
	enrolled map[string]struct{}
}

func (a *activity) snapshot() ActivityResp {
	count := len(a.enrolled)
	return ActivityResp{ID: a.id, Quota: a.quota, Remaining: a.quota - count, EnrolledCount: count}
}

func (store activityStore) create(id string, quota int) (ActivityResp, error) {
	if _, ok := store[id]; ok {
		return ActivityResp{}, ErrExists
	}
	a := &activity{id: id, quota: quota, enrolled: make(map[string]struct{})}
	store[id] = a
	return a.snapshot(), nil
}

func (store activityStore) status(id string) (ActivityResp, error) {
	a := store[id]
	if a == nil {
		return ActivityResp{}, ErrNotFound
	}
	return a.snapshot(), nil
}

func (store activityStore) update(id string, quota int) (ActivityResp, error) {
	a := store[id]
	if a == nil {
		return ActivityResp{}, ErrNotFound
	}
	if quota < len(a.enrolled) {
		return ActivityResp{}, ErrQuotaTooSmall
	}
	a.quota = quota
	return a.snapshot(), nil
}

func (store activityStore) remove(id string) error {
	if store[id] == nil {
		return ErrNotFound
	}
	delete(store, id)
	return nil
}

func (store activityStore) grab(id, userID string) (ActivityResp, error) {
	a := store[id]
	if a == nil {
		return ActivityResp{}, ErrNotFound
	}
	if _, ok := a.enrolled[userID]; ok {
		return ActivityResp{}, ErrRepeat
	}
	if len(a.enrolled) == a.quota {
		return ActivityResp{}, ErrSoldOut
	}
	a.enrolled[userID] = struct{}{}
	return a.snapshot(), nil
}
