package session

import (
	"github.com/determination66/earnth"
	"github.com/google/uuid"
)

// SessionManager 友好管理和封装Session
type SessionManager struct {
	Propagator
	Store
	CtxSessKey string
}

func (m *SessionManager) GetSession(ctx *earnth.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	// ctx.Req.earnth.Context().Value(m.CtxSessKey)
	val, ok := ctx.UserValues[m.CtxSessKey]
	if ok {
		return val.(Session), nil
	}
	// 尝试缓存住 session
	sessId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	sess, err := m.Get(ctx.Req.Context(), sessId)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.CtxSessKey] = sess
	// ctx.Req=ctx.Req.Withearnth.Context(earnth.Context.WithValue(ctx.Req.earnth.Context(), m.CtxSessKey, sess))
	return sess, err
}

func (m *SessionManager) InitSession(ctx *earnth.Context, sess Session) (Session, error) {
	id := uuid.New().String()
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	// 注入进去 HTTP 响应里面
	err = m.Inject(id, ctx.Resp)
	return sess, err
}

func (m *SessionManager) RefreshSession(ctx *earnth.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	return m.Refresh(ctx.Req.Context(), sess.ID())
}

func (m *SessionManager) RemoveSession(ctx *earnth.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}
