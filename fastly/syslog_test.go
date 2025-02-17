package fastly

import (
	"strings"
	"testing"
)

func TestClient_Syslogs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "syslogs/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	caCert := strings.TrimSpace(`
-----BEGIN CERTIFICATE-----
MIICUTCCAfugAwIBAgIBADANBgkqhkiG9w0BAQQFADBXMQswCQYDVQQGEwJDTjEL
MAkGA1UECBMCUE4xCzAJBgNVBAcTAkNOMQswCQYDVQQKEwJPTjELMAkGA1UECxMC
VU4xFDASBgNVBAMTC0hlcm9uZyBZYW5nMB4XDTA1MDcxNTIxMTk0N1oXDTA1MDgx
NDIxMTk0N1owVzELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAlBOMQswCQYDVQQHEwJD
TjELMAkGA1UEChMCT04xCzAJBgNVBAsTAlVOMRQwEgYDVQQDEwtIZXJvbmcgWWFu
ZzBcMA0GCSqGSIb3DQEBAQUAA0sAMEgCQQCp5hnG7ogBhtlynpOS21cBewKE/B7j
V14qeyslnr26xZUsSVko36ZnhiaO/zbMOoRcKK9vEcgMtcLFuQTWDl3RAgMBAAGj
gbEwga4wHQYDVR0OBBYEFFXI70krXeQDxZgbaCQoR4jUDncEMH8GA1UdIwR4MHaA
FFXI70krXeQDxZgbaCQoR4jUDncEoVukWTBXMQswCQYDVQQGEwJDTjELMAkGA1UE
CBMCUE4xCzAJBgNVBAcTAkNOMQswCQYDVQQKEwJPTjELMAkGA1UECxMCVU4xFDAS
BgNVBAMTC0hlcm9uZyBZYW5nggEAMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEE
BQADQQA/ugzBrjjK9jcWnDVfGHlk3icNRq0oV7Ri32z/+HQX67aRfgZu7KWdI+Ju
Wm7DCfrPNGVwFWUQOmsPue9rZBgO
-----END CERTIFICATE-----
`)
	clientCert := strings.TrimSpace(certificate())
	clientKey := strings.TrimSpace(privateKey())

	// Create
	var s *Syslog
	record(t, "syslogs/create", func(c *Client) {
		s, err = c.CreateSyslog(&CreateSyslogInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           String("test-syslog"),
			Address:        String("example.com"),
			Hostname:       String("example.com"),
			Port:           Int(1234),
			UseTLS:         CBool(true),
			TLSCACert:      String(caCert),
			TLSHostname:    String("example.com"),
			TLSClientCert:  String(clientCert),
			TLSClientKey:   String(clientKey),
			Token:          String("abcd1234"),
			Format:         String("format"),
			FormatVersion:  Int(2),
			MessageType:    String("classic"),
			Placement:      String("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "syslogs/cleanup", func(c *Client) {
			c.DeleteSyslog(&DeleteSyslogInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-syslog",
			})

			c.DeleteSyslog(&DeleteSyslogInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-syslog",
			})
		})
	}()

	if s.Name != "test-syslog" {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.Address != "example.com" {
		t.Errorf("bad address: %q", s.Address)
	}
	if s.Hostname != "example.com" {
		t.Errorf("bad hostname: %q", s.Hostname)
	}
	if s.Port != 1234 {
		t.Errorf("bad port: %q", s.Port)
	}
	if !s.UseTLS {
		t.Errorf("bad use_tls: %t", s.UseTLS)
	}
	if s.TLSCACert != caCert {
		t.Errorf("bad tls_ca_cert: %q", s.TLSCACert)
	}
	if s.TLSHostname != "example.com" {
		t.Errorf("bad tls_hostname: %q", s.TLSHostname)
	}
	if s.TLSClientCert != clientCert {
		t.Errorf("bad tls_client_cert: %q", s.TLSClientCert)
	}
	if s.TLSClientKey != clientKey {
		t.Errorf("bad tls_client_key: %q", s.TLSClientKey)
	}
	if s.Token != "abcd1234" {
		t.Errorf("bad token: %q", s.Token)
	}
	if s.Format != "format" {
		t.Errorf("bad format: %q", s.Format)
	}
	if s.FormatVersion != 2 {
		t.Errorf("bad format_version: %d", s.FormatVersion)
	}
	if s.MessageType != "classic" {
		t.Errorf("bad message_type: %s", s.MessageType)
	}
	if s.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", s.Placement)
	}

	// List
	var ss []*Syslog
	record(t, "syslogs/list", func(c *Client) {
		ss, err = c.ListSyslogs(&ListSyslogsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("bad syslogs: %v", ss)
	}

	// Get
	var ns *Syslog
	record(t, "syslogs/get", func(c *Client) {
		ns, err = c.GetSyslog(&GetSyslogInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-syslog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != ns.Name {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.Address != ns.Address {
		t.Errorf("bad address: %q", s.Address)
	}
	if s.Hostname != ns.Hostname {
		t.Errorf("bad hostname: %q", s.Hostname)
	}
	if s.Port != ns.Port {
		t.Errorf("bad port: %q", s.Port)
	}
	if s.UseTLS != ns.UseTLS {
		t.Errorf("bad use_tls: %t", s.UseTLS)
	}
	if s.TLSCACert != ns.TLSCACert {
		t.Errorf("bad tls_ca_cert: %q", s.TLSCACert)
	}
	if s.TLSHostname != ns.TLSHostname {
		t.Errorf("bad tls_hostname: %q", s.TLSHostname)
	}
	if s.TLSClientCert != ns.TLSClientCert {
		t.Errorf("bad tls_client_cert: %q", s.TLSClientCert)
	}
	if s.TLSClientKey != ns.TLSClientKey {
		t.Errorf("bad tls_client_key: %q", s.TLSClientKey)
	}
	if s.Token != ns.Token {
		t.Errorf("bad token: %q", s.Token)
	}
	if s.Format != ns.Format {
		t.Errorf("bad format: %q", s.Format)
	}
	if s.FormatVersion != ns.FormatVersion {
		t.Errorf("bad format_version: %q", s.FormatVersion)
	}
	if s.MessageType != ns.MessageType {
		t.Errorf("bad message_type: %q", s.MessageType)
	}
	if s.Placement != ns.Placement {
		t.Errorf("bad placement: %q", s.Placement)
	}

	// Update
	var us *Syslog
	record(t, "syslogs/update", func(c *Client) {
		us, err = c.UpdateSyslog(&UpdateSyslogInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-syslog",
			NewName:        String("new-test-syslog"),
			FormatVersion:  Int(2),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.Name != "new-test-syslog" {
		t.Errorf("bad name: %q", us.Name)
	}

	if us.FormatVersion != 2 {
		t.Errorf("bad format_version: %d", us.FormatVersion)
	}

	// Delete
	record(t, "syslogs/delete", func(c *Client) {
		err = c.DeleteSyslog(&DeleteSyslogInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-syslog",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListSyslogs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListSyslogs(&ListSyslogsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListSyslogs(&ListSyslogsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateSyslog_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateSyslog(&CreateSyslogInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateSyslog(&CreateSyslogInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetSyslog_validation(t *testing.T) {
	var err error

	_, err = testClient.GetSyslog(&GetSyslogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSyslog(&GetSyslogInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSyslog(&GetSyslogInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSyslog_validation(t *testing.T) {
	var err error

	_, err = testClient.UpdateSyslog(&UpdateSyslogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSyslog(&UpdateSyslogInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSyslog(&UpdateSyslogInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteSyslog_validation(t *testing.T) {
	var err error

	err = testClient.DeleteSyslog(&DeleteSyslogInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSyslog(&DeleteSyslogInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSyslog(&DeleteSyslogInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}
