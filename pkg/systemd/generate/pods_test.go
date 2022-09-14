package generate

import (
	"testing"

	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestHasPodExitPolicy(t *testing.T) {
	tests := []struct {
		input    []string
		expected bool
	}{
		{
			[]string{"podman", "pod", "create"},
			false,
		},
		{
			[]string{"podman", "pod", "create", "--exit-policy=continue"},
			true,
		},
		{
			[]string{"podman", "pod", "create", "--exit-policy", "continue"},
			true,
		},
	}
	for _, test := range tests {
		assert.Equalf(t, test.expected, hasPodExitPolicy(test.input), "%v", test.input)
	}
}

func TestValidateRestartPolicyPod(t *testing.T) {
	type podInfo struct {
		restart string
	}
	tests := []struct {
		name    string
		podInfo podInfo
		wantErr bool
	}{
		{"good-on", podInfo{restart: "no"}, false},
		{"good-on-success", podInfo{restart: "on-success"}, false},
		{"good-on-failure", podInfo{restart: "on-failure"}, false},
		{"good-on-abnormal", podInfo{restart: "on-abnormal"}, false},
		{"good-on-watchdog", podInfo{restart: "on-watchdog"}, false},
		{"good-on-abort", podInfo{restart: "on-abort"}, false},
		{"good-always", podInfo{restart: "always"}, false},
		{"fail", podInfo{restart: "foobar"}, true},
		{"failblank", podInfo{restart: ""}, true},
	}
	for _, tt := range tests {
		test := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRestartPolicy(test.podInfo.restart); (err != nil) != test.wantErr {
				t.Errorf("ValidateRestartPolicy() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestCreatePodSystemdUnit(t *testing.T) {
	serviceInfo := `# pod-123abc.service
`
	headerInfo := `# autogenerated by Podman CI
`
	podContent := `
[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=102
ExecStart=/usr/bin/podman start jadda-jadda-infra
ExecStop=/usr/bin/podman stop -t 42 jadda-jadda-infra
ExecStopPost=/usr/bin/podman stop -t 42 jadda-jadda-infra
PIDFile=/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid
Type=forking

[Install]
WantedBy=default.target
`
	podGood := serviceInfo + headerInfo + podContent
	podGoodNoHeaderInfo := serviceInfo + podContent

	podGoodWithEmptyPrefix := `# 123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman 123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=102
ExecStart=/usr/bin/podman start jadda-jadda-infra
ExecStop=/usr/bin/podman stop -t 42 jadda-jadda-infra
ExecStopPost=/usr/bin/podman stop -t 42 jadda-jadda-infra
PIDFile=/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid
Type=forking

[Install]
WantedBy=default.target
`

	podGoodCustomWants := `# pod-123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

# User-defined dependencies
Wants=a.service b.service c.target

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=102
ExecStart=/usr/bin/podman start jadda-jadda-infra
ExecStop=/usr/bin/podman stop -t 42 jadda-jadda-infra
ExecStopPost=/usr/bin/podman stop -t 42 jadda-jadda-infra
PIDFile=/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid
Type=forking

[Install]
WantedBy=default.target
`
	podGoodCustomAfter := `# pod-123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

# User-defined dependencies
After=a.service b.service c.target

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=102
ExecStart=/usr/bin/podman start jadda-jadda-infra
ExecStop=/usr/bin/podman stop -t 42 jadda-jadda-infra
ExecStopPost=/usr/bin/podman stop -t 42 jadda-jadda-infra
PIDFile=/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid
Type=forking

[Install]
WantedBy=default.target
`
	podGoodCustomRequires := `# pod-123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

# User-defined dependencies
Requires=a.service b.service c.target

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=102
ExecStart=/usr/bin/podman start jadda-jadda-infra
ExecStop=/usr/bin/podman stop -t 42 jadda-jadda-infra
ExecStopPost=/usr/bin/podman stop -t 42 jadda-jadda-infra
PIDFile=/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid
Type=forking

[Install]
WantedBy=default.target
`
	podGoodCustomDependencies := `# pod-123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

# User-defined dependencies
Wants=a.service b.service c.target
After=a.service b.service c.target
Requires=a.service b.service c.target

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=102
ExecStart=/usr/bin/podman start jadda-jadda-infra
ExecStop=/usr/bin/podman stop -t 42 jadda-jadda-infra
ExecStopPost=/usr/bin/podman stop -t 42 jadda-jadda-infra
PIDFile=/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid
Type=forking

[Install]
WantedBy=default.target
`
	podNoExplicitName := `# pod-123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=
Before=

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStartPre=/bin/rm -f %t/pod-123abc.pid %t/pod-123abc.pod-id
ExecStartPre=/usr/bin/podman pod create --infra-conmon-pidfile %t/pod-123abc.pid --pod-id-file %t/pod-123abc.pod-id --exit-policy=stop foo
ExecStart=/usr/bin/podman pod start --pod-id-file %t/pod-123abc.pod-id
ExecStop=/usr/bin/podman pod stop --ignore --pod-id-file %t/pod-123abc.pod-id -t 10
ExecStopPost=/usr/bin/podman pod rm --ignore -f --pod-id-file %t/pod-123abc.pod-id
PIDFile=%t/pod-123abc.pid
Type=forking

[Install]
WantedBy=default.target
`

	podGoodRestartSec := `# pod-123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
RestartSec=15
TimeoutStopSec=102
ExecStart=/usr/bin/podman start jadda-jadda-infra
ExecStop=/usr/bin/podman stop -t 42 jadda-jadda-infra
ExecStopPost=/usr/bin/podman stop -t 42 jadda-jadda-infra
PIDFile=/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid
Type=forking

[Install]
WantedBy=default.target
`

	podGoodNamedNew := `# pod-123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStartPre=/bin/rm -f %t/pod-123abc.pid %t/pod-123abc.pod-id
ExecStartPre=/usr/bin/podman pod create --infra-conmon-pidfile %t/pod-123abc.pid --pod-id-file %t/pod-123abc.pod-id --exit-policy=stop --name foo "bar=arg with space" --replace
ExecStart=/usr/bin/podman pod start --pod-id-file %t/pod-123abc.pod-id
ExecStop=/usr/bin/podman pod stop --ignore --pod-id-file %t/pod-123abc.pod-id -t 10
ExecStopPost=/usr/bin/podman pod rm --ignore -f --pod-id-file %t/pod-123abc.pod-id
PIDFile=%t/pod-123abc.pid
Type=forking

[Install]
WantedBy=default.target
`

	podGoodNamedNewWithRootArgs := `# pod-123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStartPre=/bin/rm -f %t/pod-123abc.pid %t/pod-123abc.pod-id
ExecStartPre=/usr/bin/podman --events-backend none --runroot /root pod create --infra-conmon-pidfile %t/pod-123abc.pid --pod-id-file %t/pod-123abc.pod-id --exit-policy=stop --name foo "bar=arg with space" --replace
ExecStart=/usr/bin/podman --events-backend none --runroot /root pod start --pod-id-file %t/pod-123abc.pod-id
ExecStop=/usr/bin/podman --events-backend none --runroot /root pod stop --ignore --pod-id-file %t/pod-123abc.pod-id -t 10
ExecStopPost=/usr/bin/podman --events-backend none --runroot /root pod rm --ignore -f --pod-id-file %t/pod-123abc.pod-id
PIDFile=%t/pod-123abc.pid
Type=forking

[Install]
WantedBy=default.target
`

	podGoodNamedNewWithReplaceFalse := `# pod-123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStartPre=/bin/rm -f %t/pod-123abc.pid %t/pod-123abc.pod-id
ExecStartPre=/usr/bin/podman pod create --infra-conmon-pidfile %t/pod-123abc.pid --pod-id-file %t/pod-123abc.pod-id --exit-policy=stop --name foo --replace
ExecStart=/usr/bin/podman pod start --pod-id-file %t/pod-123abc.pod-id
ExecStop=/usr/bin/podman pod stop --ignore --pod-id-file %t/pod-123abc.pod-id -t 10
ExecStopPost=/usr/bin/podman pod rm --ignore -f --pod-id-file %t/pod-123abc.pod-id
PIDFile=%t/pod-123abc.pid
Type=forking

[Install]
WantedBy=default.target
`

	podNewLabelWithCurlyBraces := `# pod-123abc.service
# autogenerated by Podman CI

[Unit]
Description=Podman pod-123abc.service
Documentation=man:podman-generate-systemd(1)
Wants=network-online.target
After=network-online.target
RequiresMountsFor=/var/run/containers/storage
Requires=container-1.service container-2.service
Before=container-1.service container-2.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStartPre=/bin/rm -f %t/pod-123abc.pid %t/pod-123abc.pod-id
ExecStartPre=/usr/bin/podman pod create --infra-conmon-pidfile %t/pod-123abc.pid --pod-id-file %t/pod-123abc.pod-id --name foo --label key={{someval}} --exit-policy=continue --replace
ExecStart=/usr/bin/podman pod start --pod-id-file %t/pod-123abc.pod-id
ExecStop=/usr/bin/podman pod stop --ignore --pod-id-file %t/pod-123abc.pod-id -t 10
ExecStopPost=/usr/bin/podman pod rm --ignore -f --pod-id-file %t/pod-123abc.pod-id
PIDFile=%t/pod-123abc.pid
Type=forking

[Install]
WantedBy=default.target
`

	tests := []struct {
		name     string
		info     podInfo
		want     string
		new      bool
		noHeader bool
		wantErr  bool
	}{
		{"pod",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      42,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				CreateCommand:    []string{"podman", "pod", "create", "--name", "foo", "bar=arg with space"},
			},
			podGood,
			false,
			false,
			false,
		},
		{"pod",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      42,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				Wants:            []string{"a.service", "b.service", "c.target"},
				CreateCommand: []string{
					"podman", "pod", "create", "--name", "foo", "--wants", "a.service",
					"--wants", "b.service", "--wants", "c.target", "bar=arg with space"},
			},
			podGoodCustomWants,
			false,
			false,
			false,
		},
		{"pod",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      42,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				After:            []string{"a.service", "b.service", "c.target"},
				CreateCommand: []string{
					"podman", "pod", "create", "--name", "foo", "--after", "a.service",
					"--after", "b.service", "--after", "c.target", "bar=arg with space"},
			},
			podGoodCustomAfter,
			false,
			false,
			false,
		},
		{"pod",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      42,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				Requires:         []string{"a.service", "b.service", "c.target"},
				CreateCommand: []string{
					"podman", "pod", "create", "--name", "foo", "--requires", "a.service",
					"--requires", "b.service", "--requires", "c.target", "bar=arg with space"},
			},
			podGoodCustomRequires,
			false,
			false,
			false,
		},
		{"pod",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      42,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				Wants:            []string{"a.service", "b.service", "c.target"},
				After:            []string{"a.service", "b.service", "c.target"},
				Requires:         []string{"a.service", "b.service", "c.target"},
				CreateCommand: []string{
					"podman", "pod", "create", "--name", "foo", "--wants", "a.service",
					"--wants", "b.service", "--wants", "c.target", "--after", "a.service",
					"--after", "b.service", "--after", "c.target", "--requires", "a.service",
					"--requires", "b.service", "--requires", "c.target", "bar=arg with space"},
			},
			podGoodCustomDependencies,
			false,
			false,
			false,
		},
		{"pod without --name",
			podInfo{
				Executable:    "/usr/bin/podman",
				ServiceName:   "pod-123abc",
				InfraNameOrID: "jadda-jadda-infra",
				PIDFile:       "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:   10,
				PodmanVersion: "CI",
				GraphRoot:     "/var/lib/containers/storage",
				RunRoot:       "/var/run/containers/storage",
				CreateCommand: []string{"podman", "pod", "create", "foo"},
			},
			podNoExplicitName,
			true,
			false,
			false,
		},
		{"pod restartSec",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      42,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				CreateCommand:    []string{"podman", "pod", "create", "--name", "foo", "bar=arg with space"},
				RestartSec:       15,
			},
			podGoodRestartSec,
			false,
			false,
			false,
		},
		{"pod noHeader",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      42,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				CreateCommand:    []string{"podman", "pod", "create", "--name", "foo", "bar=arg with space"},
			},
			podGoodNoHeaderInfo,
			false,
			true,
			false,
		},
		{"pod with root args",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      42,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				CreateCommand:    []string{"podman", "--events-backend", "none", "--runroot", "/root", "pod", "create", "--name", "foo", "bar=arg with space"},
			},
			podGood,
			false,
			false,
			false,
		},
		{"pod --new",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      10,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				CreateCommand:    []string{"podman", "pod", "create", "--name", "foo", "bar=arg with space"},
			},
			podGoodNamedNew,
			true,
			false,
			false,
		},
		{"pod --new with root args",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      10,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				CreateCommand:    []string{"podman", "--events-backend", "none", "--runroot", "/root", "pod", "create", "--name", "foo", "bar=arg with space"},
			},
			podGoodNamedNewWithRootArgs,
			true,
			false,
			false,
		},
		{"pod --new with --replace=false",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      10,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				CreateCommand:    []string{"podman", "pod", "create", "--name", "foo", "--replace=false"},
			},
			podGoodNamedNewWithReplaceFalse,
			true,
			false,
			false,
		},
		{"pod --new with double curly braces",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      10,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				CreateCommand:    []string{"podman", "pod", "create", "--name", "foo", "--label", "key={{someval}}", "--exit-policy=continue"},
			},
			podNewLabelWithCurlyBraces,
			true,
			false,
			false,
		},
		{"pod --new with ID files",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "pod-123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      10,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				CreateCommand:    []string{"podman", "pod", "create", "--infra-conmon-pidfile", "/tmp/pod-123abc.pid", "--pod-id-file", "/tmp/pod-123abc.pod-id", "--name", "foo", "bar=arg with space"},
			},
			podGoodNamedNew,
			true,
			false,
			false,
		},
		{"pod with empty pod-prefix",
			podInfo{
				Executable:       "/usr/bin/podman",
				ServiceName:      "123abc",
				InfraNameOrID:    "jadda-jadda-infra",
				PIDFile:          "/run/containers/storage/overlay-containers/639c53578af4d84b8800b4635fa4e680ee80fd67e0e6a2d4eea48d1e3230f401/userdata/conmon.pid",
				StopTimeout:      42,
				PodmanVersion:    "CI",
				GraphRoot:        "/var/lib/containers/storage",
				RunRoot:          "/var/run/containers/storage",
				RequiredServices: []string{"container-1", "container-2"},
				CreateCommand:    []string{"podman", "pod", "create", "--name", "foo", "bar=arg with space"},
			},
			podGoodWithEmptyPrefix,
			false,
			false,
			false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(tt.name, func(t *testing.T) {
			opts := entities.GenerateSystemdOptions{
				New:      test.new,
				NoHeader: test.noHeader,
			}
			got, err := executePodTemplate(&test.info, opts)
			if (err != nil) != test.wantErr {
				t.Errorf("CreatePodSystemdUnit() error = \n%v, wantErr \n%v", err, test.wantErr)
				return
			}
			assert.Equal(t, test.want, got)
		})
	}
}
