package manager

import (
	"testing"

	"github.com/grafana/grafana/pkg/infra/usagestats"
	"github.com/grafana/grafana/pkg/services/encryption/ossencryption"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/kmsproviders/osskmsproviders"
	"github.com/grafana/grafana/pkg/services/secrets"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
)

func SetupTestService(tb testing.TB, store secrets.Store) *SecretsService {
	tb.Helper()
	defaultKey := "SdlklWklckeLS"
	if len(setting.SecretKey) > 0 {
		defaultKey = setting.SecretKey
	}
	raw, err := ini.Load([]byte(`
		[security]
		secret_key = ` + defaultKey))
	require.NoError(tb, err)

	features := featuremgmt.WithFeatures(featuremgmt.FlagEnvelopeEncryption)

	cfg := &setting.Cfg{Raw: raw}
	cfg.IsFeatureToggleEnabled = features.IsEnabled

	settings := &setting.OSSImpl{Cfg: cfg}
	assert.True(tb, settings.IsFeatureToggleEnabled(featuremgmt.FlagEnvelopeEncryption))
	assert.True(tb, features.IsEnabled(featuremgmt.FlagEnvelopeEncryption))

	encryption := ossencryption.ProvideService()
	secretsService, err := ProvideSecretsService(
		store,
		osskmsproviders.ProvideService(encryption, settings),
		encryption,
		settings,
		&usagestats.UsageStatsMock{T: tb},
	)
	require.NoError(tb, err)

	return secretsService
}
