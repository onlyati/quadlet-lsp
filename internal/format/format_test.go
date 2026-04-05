package format

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/require"
)

// Test_FormatDocument tests document format method.
func Test_FormatDocument(t *testing.T) {
	source := `# disable-qsr: qsr014
# disable-qsr: qsr015

[Unit]
Description=Nextcloud instance

[Container]
Pod=nc.pod
AutoUpdate=registry

Memory=512M

Volume=nc-app.volume:/var/www/html

Environment=POSTGRES_USER=nextclouduser
# disable-qsr: qsrxxx
Environment=POSTGRES_DB=nextcloud
Environment=POSTGRES_HOST=127.0.0.1
Secret=nc-db-password,type=env,target=POSTGRES_PASSWORD

Environment=NEXTCLOUD_ADMIN_USER=ati
Secret=nc-admin-pw,type=env,target=NEXTCLOUD_ADMIN_PASSWORD

Environment=REDIS_HOST=127.0.0.1
Environment=REDIS_PORT=6379

Environment=SMTP_HOST=smtp.rackhost.hu
Environment=SMTP_SECURE=tls
Environment=SMTP_PORT=587
Environment=SMTP_NAME=noreply@thinkaboutit.tech
Environment=SMTP_DOMAIN=thinkaboutit.tech
Environment=SMTP_FROM_ADDRESS=noreply@thinkaboutit.tech
Secret=tai-noreply,type=env,target=SMTP_PASSWORD

Environment=OBJECTSTORE_S3_BUCKET=dakota-bazooka-metaphor-axes
Environment=OBJECTSTORE_S3_REGION=de
Environment=OBJECTSTORE_S3_HOST=s3.de.io.cloud.ovh.net
Secret=ovh-s3-access-key,type=env,target=OBJECTSTORE_S3_KEY
Secret=ovh-s3-secret-key,type=env,target=OBJECTSTORE_S3_SECRET

PublishPort=8080:8080

[Service]
Restart=on-failure
RestartSec=5
StartLimitBurst=5

[Install]
WantedBy=default.target
`

	expected := `# disable-qsr: qsr014
# disable-qsr: qsr015

[Unit]
Description=Nextcloud instance

[Container]
AutoUpdate=registry
Pod=nc.pod

Environment=NEXTCLOUD_ADMIN_USER=ati
Environment=OBJECTSTORE_S3_BUCKET=dakota-bazooka-metaphor-axes
Environment=OBJECTSTORE_S3_HOST=s3.de.io.cloud.ovh.net
Environment=OBJECTSTORE_S3_REGION=de
# disable-qsr: qsrxxx
Environment=POSTGRES_DB=nextcloud
Environment=POSTGRES_HOST=127.0.0.1
Environment=POSTGRES_USER=nextclouduser
Environment=REDIS_HOST=127.0.0.1
Environment=REDIS_PORT=6379
Environment=SMTP_DOMAIN=thinkaboutit.tech
Environment=SMTP_FROM_ADDRESS=noreply@thinkaboutit.tech
Environment=SMTP_HOST=smtp.rackhost.hu
Environment=SMTP_NAME=noreply@thinkaboutit.tech
Environment=SMTP_PORT=587
Environment=SMTP_SECURE=tls

PublishPort=8080:8080

Volume=nc-app.volume:/var/www/html

Secret=nc-admin-pw,type=env,target=NEXTCLOUD_ADMIN_PASSWORD
Secret=nc-db-password,type=env,target=POSTGRES_PASSWORD
Secret=ovh-s3-access-key,type=env,target=OBJECTSTORE_S3_KEY
Secret=ovh-s3-secret-key,type=env,target=OBJECTSTORE_S3_SECRET
Secret=tai-noreply,type=env,target=SMTP_PASSWORD

Memory=512M

[Service]
Restart=on-failure
RestartSec=5
StartLimitBurst=5

[Install]
WantedBy=default.target
`
	parser := parser.NewParserFromMemory("foo.container", source)
	newText := FormatDocument(parser.Quadlet)

	for i, c := range newText {
		require.Less(t, i, len(expected), "result longer than expected")
		require.Equalf(t, c, rune(expected[i]), "unexpected text at formatting at %d, %c %c", i, c, rune(expected[i]))
	}
}

// Test_FormatDocumentMultiLine tests formatting when continuation sign is used.
func Test_FormatDocumentMultiLine(t *testing.T) {
	source := ` [Unit]
Description=PostgreSQL database for Nextcloud

[Container]
Pod=nc.pod
Image=docker.io/library/postgres:17
Exec=postgres \
  -c listen_addresses=127.0.0.1
AutoUpdate=registry

Label=asd=asd
Label= \
  asd=asd
Memory=512M
`

	expected := `[Unit]
Description=PostgreSQL database for Nextcloud

[Container]
AutoUpdate=registry
Exec=postgres -c listen_addresses=127.0.0.1
Image=docker.io/library/postgres:17
Pod=nc.pod

Label=asd=asd
Label=asd=asd

Memory=512M

`
	parser := parser.NewParserFromMemory("foo.container", source)
	newText := FormatDocument(parser.Quadlet)

	for i, c := range newText {
		require.Less(t, i, len(expected), "result longer than expected")
		require.Equalf(t, c, rune(expected[i]), "unexpected text at formatting at %d, %c %c", i, c, rune(expected[i]))
	}
}

// Test_WrapLine tests that formatting make long lines shorter.
func Test_WrapLine(t *testing.T) {
	source := `[Container]
# Healthcheck options
HealthCmd=/bin/curlcurl -k --fail --connect-timeout 5 -X POST -d '{ "a": "place" }' https://127.0.0.1:3000/api/healthz
HealthRetries=10
HealthStartPeriod=15s
HealthTimeout=15s
`
	expected := `[Container]
# Healthcheck options
HealthCmd=/bin/curlcurl -k --fail --connect-timeout 5 -X POST -d '{ "a": "place" }' \
  https://127.0.0.1:3000/api/healthz
HealthRetries=10
HealthStartPeriod=15s
HealthTimeout=15s

`
	parser := parser.NewParserFromMemory("foo.container", source)
	newText := FormatDocument(parser.Quadlet)

	for i, c := range newText {
		require.Less(t, i, len(expected), "result longer than expected")
		require.Equalf(t, c, rune(expected[i]), "unexpected text at formatting at %d, %c %c", i, c, rune(expected[i]))
	}
}

// Test_Unchanged tests if something format is correct, then do nothing.
func Test_Unchanged(t *testing.T) {
	source := `[Container]
# Healthcheck options
HealthCmd=/bin/curlcurl -k --fail --connect-timeout 5 -X POST -d '{ "a": "place" }' \
  https://127.0.0.1:3000/api/healthz
HealthRetries=10
HealthStartPeriod=15s
HealthTimeout=15s
`
	parser := parser.NewParserFromMemory("foo.container", source)
	newText := FormatDocument(parser.Quadlet)

	for i, c := range newText {
		require.Less(t, i, len(source), "result longer than source")
		require.Equalf(t, c, rune(source[i]), "unexpected text at formatting at %d, %c %c", i, c, rune(source[i]))
	}
}
