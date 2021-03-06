package zerolog

import (
	"context"
	"os"

	"omolle.com/fnd/server/model"

	pkgctx "omolle.com/fnd/server/pkg/context"

	"github.com/rs/zerolog"
)

// ZLog represents zerolog logger
type ZLog struct {
	logger *zerolog.Logger
}

// New instantiates new zero logger
func New() *ZLog {
	z := zerolog.New(os.Stdout)
	return &ZLog{
		logger: &z,
	}
}

// Log logs using zerolog
func (z *ZLog) Log(ctx context.Context, source, msg string, err error, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	params["source"] = source

	if user, ok := ctx.Value(pkgctx.KeyString("_authuser")).(*twisk.AuthUser); ok {
		params["id"] = user.ID
		params["username"] = user.Username
		params["tenant_id"] = user.TenantID
	}

	if err != nil {
		params["error"] = err
		z.logger.Error().Fields(params).Msg(msg)
		return
	}

	z.logger.Info().Fields(params).Msg(msg)

}
