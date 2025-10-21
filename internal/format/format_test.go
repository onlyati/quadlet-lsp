package format

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FormatDocument(t *testing.T) {
	source := `# disable-qsr: qsr014
# disable-qsr: qsr015

[Unit]
Description=Nextcloud instance

[Container]
Pod=nc.pod
AutoUpdate=registry

Memory=512M

# Volumes
Volume=nc-app.volume:/var/www/html

# Environment variables

# Database variables
Environment=POSTGRES_USER=nextclouduser
Environment=POSTGRES_DB=nextcloud
Environment=POSTGRES_HOST=127.0.0.1
Secret=nc-db-password,type=env,target=POSTGRES_PASSWORD

# Default admin user and password
Environment=NEXTCLOUD_ADMIN_USER=ati
Secret=nc-admin-pw,type=env,target=NEXTCLOUD_ADMIN_PASSWORD

# Redis variables
Environment=REDIS_HOST=127.0.0.1
Environment=REDIS_PORT=6379

# SMTP variables
Environment=SMTP_HOST=smtp.rackhost.hu
Environment=SMTP_SECURE=tls
Environment=SMTP_PORT=587
Environment=SMTP_NAME=noreply@thinkaboutit.tech
Environment=SMTP_DOMAIN=thinkaboutit.tech
Environment=SMTP_FROM_ADDRESS=noreply@thinkaboutit.tech
Secret=tai-noreply,type=env,target=SMTP_PASSWORD

# S3 bucket as primary storage
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
# Base options
AutoUpdate=registry
Pod=nc.pod

# Storage options
Volume=nc-app.volume:/var/www/html

# Network options
PublishPort=8080:8080

# Environment options
Environment=NEXTCLOUD_ADMIN_USER=ati
Environment=OBJECTSTORE_S3_BUCKET=dakota-bazooka-metaphor-axes
Environment=OBJECTSTORE_S3_HOST=s3.de.io.cloud.ovh.net
Environment=OBJECTSTORE_S3_REGION=de
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

# Secret options
Secret=nc-admin-pw,type=env,target=NEXTCLOUD_ADMIN_PASSWORD
Secret=nc-db-password,type=env,target=POSTGRES_PASSWORD
Secret=ovh-s3-access-key,type=env,target=OBJECTSTORE_S3_KEY
Secret=ovh-s3-secret-key,type=env,target=OBJECTSTORE_S3_SECRET
Secret=tai-noreply,type=env,target=SMTP_PASSWORD

# Other options
Memory=512M

[Service]
Restart=on-failure
RestartSec=5
StartLimitBurst=5

[Install]
WantedBy=default.target

`
	newText := FormatDocument(source)

	for i, c := range newText {
		if i > len(expected) {
			t.Fatalf("result longer than expected '%s'", newText[i:])
		}
		if c != rune(expected[i]) {
			t.Fatalf("unpextected text after formatting, at position %d: got: '%v' expected: '%v'", i, newText[i-10:i+1], expected[i-10:i+1])
		}
	}
}

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
# Base options
AutoUpdate=registry
Exec=postgres -c listen_addresses=127.0.0.1
Image=docker.io/library/postgres:17
Pod=nc.pod

# Label options
Label=asd=asd
Label=asd=asd

# Other options
Memory=512M

`
	newText := FormatDocument(source)

	for i, c := range newText {
		if i > len(expected) {
			t.Fatalf("result longer than expected '%s'", newText[i:])
		}
		if c != rune(expected[i]) {
			t.Fatalf("unpextected text after formatting, at position %d: got: '%v' expected: '%v'", i, newText[i-10:i+1], expected[i-10:i+1])
		}
	}
}

func Test_WrapLine(t *testing.T) {
	source := `[Container]
# Healthcheck options
HealthCmd=/bin/curlcurl -k --fail --connect-timeout 5 https://127.0.0.1:3000/api/healthz
HealthRetries=10
HealthStartPeriod=15s
HealthTimeout=15s
`
	expected := `[Container]
# Healthcheck options
HealthCmd=/bin/curlcurl -k --fail --connect-timeout 5 \
  https://127.0.0.1:3000/api/healthz
HealthRetries=10
HealthStartPeriod=15s
HealthTimeout=15s

`
	newText := FormatDocument(source)

	for i, c := range newText {
		if i >= len(expected) {
			t.Fatalf("result longer than expected '%s'", newText[i:])
		}
		if c != rune(expected[i]) {
			t.Fatalf("unpextected text after formatting, at position %d: got:\n'%v'\nexpected:\n'%v'", i, newText[i-10:i+1], expected[i-10:i+1])
		}
	}
}

func Test_Unchanged(t *testing.T) {
	source := `[Container]
# Healthcheck options
HealthCmd=/bin/curlcurl -k --fail --connect-timeout 5 \
  https://127.0.0.1:3000/api/healthz
HealthRetries=10
HealthStartPeriod=15s
HealthTimeout=15s

`
	newText := FormatDocument(source)

	for i, c := range newText {
		if i > len(source) {
			t.Fatalf("result longer than expected '%s'", newText[i:])
		}
		if c != rune(source[i]) {
			t.Fatalf("unpextected text after formatting, at position %d: got:\n'%v'\nexpected:\n'%v'", i, newText[i-10:i+10], source[i-10:i+10])
		}
	}
}

func Test_Wrap(t *testing.T) {
	sources := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas nunc mauris, pharetra quis nisi in, eleifend vulputate nisl. Fusce justo mauris, aliquam sed urna feugiat, accumsan egestas tellus. Maecenas ut felis a leo tincidunt volutpat eget a nibh.",
		"HealthCmd=/bin/curl -k --fail  --connect-timeout 5 https://127.0.0.1:3000/api/healthz",
		"Image=ghcr.io/immich-app/postgres:14-vectorchord0.4.3-pgvectors0.2.0@sha256:41eacbe83eca995561fe43814fd4891e16e39632806253848efaf04d3c8a8b84",
		"HealthCmd=/bin/curlcurl -k --fail --connect-timeout 5 https://127.0.0.1:3000/api/healthzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz asd",
	}
	expected := []string{
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas nunc mauris, \
  pharetra quis nisi in, eleifend vulputate nisl. Fusce justo mauris, aliquam \
  sed urna feugiat, accumsan egestas tellus. Maecenas ut felis a leo \
  tincidunt volutpat eget a nibh.
`,
		`HealthCmd=/bin/curl -k --fail  --connect-timeout 5 \
  https://127.0.0.1:3000/api/healthz
`,
		`Image=ghcr.io/immich-app/postgres:14-vectorchord0.4.3-pgvectors0.2.0@sha256:41eacbe83eca995561fe43814fd4891e16e39632806253848efaf04d3c8a8b84
`,
		`HealthCmd=/bin/curlcurl -k --fail --connect-timeout 5 \
  https://127.0.0.1:3000/api/healthzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz \
  asd
`,
	}

	for i, s := range sources {
		r := wrapLine(s, 80)
		require.Equal(t, expected[i], r, "invalid wrap")
	}
}
