package infra

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestIncusModule_Instance(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)
	RegisterInstanceMetatable(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "create instance builder",
			script: `
				local inst = incus.instance("test-vm")
				assert(inst ~= nil, "instance should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "instance with image",
			script: `
				local inst = incus.instance("test-vm"):image("ubuntu:22.04")
				assert(inst ~= nil, "instance should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "instance with config",
			script: `
				local inst = incus.instance("test-vm")
					:image("ubuntu:22.04")
					:config({
						["limits.cpu"] = "2",
						["limits.memory"] = "2GB"
					})
				assert(inst ~= nil, "instance should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "instance with profile",
			script: `
				local inst = incus.instance("test-vm")
					:image("ubuntu:22.04")
					:profile("default")
					:profile("docker")
				assert(inst ~= nil, "instance should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "instance with device",
			script: `
				local inst = incus.instance("test-vm")
					:image("ubuntu:22.04")
					:device("eth0", {
						type = "nic",
						network = "incusbr0"
					})
				assert(inst ~= nil, "instance should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "instance with delegate_to",
			script: `
				local inst = incus.instance("test-vm")
					:image("ubuntu:22.04")
					:delegate_to("remote-host")
				assert(inst ~= nil, "instance should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "instance ephemeral",
			script: `
				local inst = incus.instance("test-vm")
					:image("ubuntu:22.04")
					:ephemeral(true)
				assert(inst ~= nil, "instance should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "complete instance configuration",
			script: `
				local inst = incus.instance("web-server")
					:image("ubuntu:22.04")
					:profile("default")
					:profile("web")
					:config({
						["limits.cpu"] = "4",
						["limits.memory"] = "4GB",
						["boot.autostart"] = "true"
					})
					:device("root", {
						type = "disk",
						path = "/",
						pool = "default"
					})
					:device("eth0", {
						type = "nic",
						network = "incusbr0"
					})
					:delegate_to("prod-server")
					:ephemeral(false)
				assert(inst ~= nil, "instance should not be nil")
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_Image(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)
	RegisterImageMetatable(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "create image builder",
			script: `
				local img = incus.image("ubuntu-22.04")
				assert(img ~= nil, "image should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "image with source",
			script: `
				local img = incus.image("ubuntu-22.04")
					:source("ubuntu/22.04")
				assert(img ~= nil, "image should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "image with server",
			script: `
				local img = incus.image("ubuntu-22.04")
					:source("ubuntu/22.04")
					:server("https://images.linuxcontainers.org")
				assert(img ~= nil, "image should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "image with delegate_to",
			script: `
				local img = incus.image("ubuntu-22.04")
					:source("ubuntu/22.04")
					:delegate_to("remote-host")
				assert(img ~= nil, "image should not be nil")
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_Network(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)
	RegisterNetworkMetatable(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "create network builder",
			script: `
				local net = incus.network("incusbr0")
				assert(net ~= nil, "network should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "network with type",
			script: `
				local net = incus.network("incusbr0")
					:type("bridge")
				assert(net ~= nil, "network should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "network with config",
			script: `
				local net = incus.network("incusbr0")
					:type("bridge")
					:config({
						["ipv4.address"] = "10.0.0.1/24",
						["ipv4.nat"] = "true",
						["ipv6.address"] = "none"
					})
				assert(net ~= nil, "network should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "network with delegate_to",
			script: `
				local net = incus.network("incusbr0")
					:type("bridge")
					:delegate_to("remote-host")
				assert(net ~= nil, "network should not be nil")
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_Profile(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)
	RegisterProfileMetatable(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "create profile builder",
			script: `
				local prof = incus.profile("docker")
				assert(prof ~= nil, "profile should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "profile with description",
			script: `
				local prof = incus.profile("docker")
					:description("Docker enabled profile")
				assert(prof ~= nil, "profile should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "profile with config",
			script: `
				local prof = incus.profile("docker")
					:config({
						["security.nesting"] = "true",
						["security.syscalls.intercept.mknod"] = "true"
					})
				assert(prof ~= nil, "profile should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "profile with device",
			script: `
				local prof = incus.profile("docker")
					:device("eth0", {
						type = "nic",
						network = "incusbr0"
					})
				assert(prof ~= nil, "profile should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "profile with delegate_to",
			script: `
				local prof = incus.profile("docker")
					:delegate_to("remote-host")
				assert(prof ~= nil, "profile should not be nil")
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_Storage(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)
	RegisterStorageMetatable(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "create storage builder",
			script: `
				local stor = incus.storage("default")
				assert(stor ~= nil, "storage should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "storage with driver",
			script: `
				local stor = incus.storage("default")
					:driver("zfs")
				assert(stor ~= nil, "storage should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "storage with config",
			script: `
				local stor = incus.storage("default")
					:driver("zfs")
					:config({
						source = "/dev/sdb",
						["zfs.pool_name"] = "incus-pool"
					})
				assert(stor ~= nil, "storage should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "storage with delegate_to",
			script: `
				local stor = incus.storage("default")
					:driver("dir")
					:delegate_to("remote-host")
				assert(stor ~= nil, "storage should not be nil")
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_Snapshot(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)
	RegisterSnapshotMetatable(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "create snapshot builder",
			script: `
				local snap = incus.snapshot("web-server", "backup-2024")
				assert(snap ~= nil, "snapshot should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "snapshot stateful",
			script: `
				local snap = incus.snapshot("web-server", "backup-2024")
					:stateful(true)
				assert(snap ~= nil, "snapshot should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "snapshot with delegate_to",
			script: `
				local snap = incus.snapshot("web-server", "backup-2024")
					:stateful(false)
					:delegate_to("remote-host")
				assert(snap ~= nil, "snapshot should not be nil")
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_UtilityFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "exec function exists",
			script: `
				assert(type(incus.exec) == "function", "exec should be a function")
			`,
			wantErr: false,
		},
		{
			name: "list function exists",
			script: `
				assert(type(incus.list) == "function", "list should be a function")
			`,
			wantErr: false,
		},
		{
			name: "info function exists",
			script: `
				assert(type(incus.info) == "function", "info should be a function")
			`,
			wantErr: false,
		},
		{
			name: "delete function exists",
			script: `
				assert(type(incus.delete) == "function", "delete should be a function")
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_Exec(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "exec requires instance",
			script: `
				local result, err = incus.exec({
					command = "ls -la"
				})
				assert(result == nil, "result should be nil when instance is missing")
				assert(err ~= nil, "error should not be nil when instance is missing")
			`,
			wantErr: false,
		},
		{
			name: "exec requires command",
			script: `
				local result, err = incus.exec({
					instance = "test-vm"
				})
				assert(result == nil, "result should be nil when command is missing")
				assert(err ~= nil, "error should not be nil when command is missing")
			`,
			wantErr: false,
		},
		{
			name: "exec with all parameters",
			script: `
				-- Note: This will fail without a real instance, but tests the parameter parsing
				local result, err = incus.exec({
					instance = "test-vm",
					command = "whoami",
					user = "ubuntu",
					group = "ubuntu",
					cwd = "/home/ubuntu",
					env = {
						PATH = "/usr/bin:/bin",
						HOME = "/home/ubuntu"
					},
					target = "remote-host"
				})
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_List(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "list instances",
			script: `
				-- Note: This will fail without incus installed, but tests the parameter parsing
				local result, err = incus.list({
					type = "instances",
					format = "json"
				})
			`,
			wantErr: false,
		},
		{
			name: "list images",
			script: `
				local result, err = incus.list({
					type = "images",
					format = "json"
				})
			`,
			wantErr: false,
		},
		{
			name: "list networks",
			script: `
				local result, err = incus.list({
					type = "networks",
					format = "json"
				})
			`,
			wantErr: false,
		},
		{
			name: "list profiles",
			script: `
				local result, err = incus.list({
					type = "profiles",
					format = "json"
				})
			`,
			wantErr: false,
		},
		{
			name: "list storage",
			script: `
				local result, err = incus.list({
					type = "storage",
					format = "json"
				})
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_Info(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "info requires name",
			script: `
				local result, err = incus.info({
					type = "instance"
				})
				assert(result == nil, "result should be nil when name is missing")
				assert(err ~= nil, "error should not be nil when name is missing")
			`,
			wantErr: false,
		},
		{
			name: "info instance",
			script: `
				local result, err = incus.info({
					type = "instance",
					name = "test-vm"
				})
			`,
			wantErr: false,
		},
		{
			name: "info image",
			script: `
				local result, err = incus.info({
					type = "image",
					name = "ubuntu-22.04"
				})
			`,
			wantErr: false,
		},
		{
			name: "info network",
			script: `
				local result, err = incus.info({
					type = "network",
					name = "incusbr0"
				})
			`,
			wantErr: false,
		},
		{
			name: "info profile",
			script: `
				local result, err = incus.info({
					type = "profile",
					name = "default"
				})
			`,
			wantErr: false,
		},
		{
			name: "info storage",
			script: `
				local result, err = incus.info({
					type = "storage",
					name = "default"
				})
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_Delete(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "delete requires name",
			script: `
				local result, err = incus.delete({
					type = "instance"
				})
				assert(result == nil, "result should be nil when name is missing")
				assert(err ~= nil, "error should not be nil when name is missing")
			`,
			wantErr: false,
		},
		{
			name: "delete instance",
			script: `
				local result, err = incus.delete({
					type = "instance",
					name = "test-vm"
				})
			`,
			wantErr: false,
		},
		{
			name: "delete instance with force",
			script: `
				local result, err = incus.delete({
					type = "instance",
					name = "test-vm",
					force = true
				})
			`,
			wantErr: false,
		},
		{
			name: "delete image",
			script: `
				local result, err = incus.delete({
					type = "image",
					name = "ubuntu-22.04"
				})
			`,
			wantErr: false,
		},
		{
			name: "delete network",
			script: `
				local result, err = incus.delete({
					type = "network",
					name = "incusbr0"
				})
			`,
			wantErr: false,
		},
		{
			name: "delete profile",
			script: `
				local result, err = incus.delete({
					type = "profile",
					name = "docker"
				})
			`,
			wantErr: false,
		},
		{
			name: "delete storage",
			script: `
				local result, err = incus.delete({
					type = "storage",
					name = "default"
				})
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_ChainedMethods(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)
	RegisterInstanceMetatable(L)
	RegisterNetworkMetatable(L)
	RegisterProfileMetatable(L)
	RegisterStorageMetatable(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "complex instance chain",
			script: `
				local inst = incus.instance("web-app")
					:image("ubuntu:22.04")
					:profile("default")
					:profile("docker")
					:config({
						["limits.cpu"] = "4",
						["limits.memory"] = "8GB",
						["boot.autostart"] = "true"
					})
					:device("root", {
						type = "disk",
						path = "/",
						pool = "default",
						size = "50GB"
					})
					:device("eth0", {
						type = "nic",
						network = "incusbr0",
						["ipv4.address"] = "10.0.0.10"
					})
					:device("data", {
						type = "disk",
						source = "/mnt/data",
						path = "/data"
					})
					:delegate_to("prod-server")
					:ephemeral(false)
				assert(inst ~= nil, "instance should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "complex network chain",
			script: `
				local net = incus.network("prod-net")
					:type("bridge")
					:config({
						["bridge.mode"] = "fan",
						["fan.underlay_subnet"] = "auto",
						["ipv4.address"] = "10.10.0.1/24",
						["ipv4.nat"] = "true",
						["ipv4.dhcp"] = "true",
						["ipv6.address"] = "none",
						["dns.mode"] = "managed"
					})
					:delegate_to("network-host")
				assert(net ~= nil, "network should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "complex profile chain",
			script: `
				local prof = incus.profile("gpu-enabled")
					:description("Profile for GPU-enabled containers")
					:config({
						["security.nesting"] = "true",
						["nvidia.runtime"] = "true",
						["nvidia.driver.capabilities"] = "all"
					})
					:device("gpu", {
						type = "gpu",
						gputype = "physical",
						pci = "0000:01:00.0"
					})
					:device("eth0", {
						type = "nic",
						network = "incusbr0"
					})
					:delegate_to("gpu-server")
				assert(prof ~= nil, "profile should not be nil")
			`,
			wantErr: false,
		},
		{
			name: "complex storage chain",
			script: `
				local stor = incus.storage("nvme-pool")
					:driver("zfs")
					:config({
						source = "/dev/nvme0n1",
						["volume.block.filesystem"] = "ext4",
						["volume.size"] = "100GB",
						["zfs.pool_name"] = "nvme-zpool"
					})
					:delegate_to("storage-host")
				assert(stor ~= nil, "storage should not be nil")
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncusModule_ModuleRegistration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewIncusModule(nil)
	module.Register(L)

	// Verificar se o módulo foi registrado corretamente
	err := L.DoString(`
		assert(incus ~= nil, "incus module should be registered")
		assert(type(incus.instance) == "function", "instance should be a function")
		assert(type(incus.image) == "function", "image should be a function")
		assert(type(incus.network) == "function", "network should be a function")
		assert(type(incus.profile) == "function", "profile should be a function")
		assert(type(incus.storage) == "function", "storage should be a function")
		assert(type(incus.snapshot) == "function", "snapshot should be a function")
		assert(type(incus.exec) == "function", "exec should be a function")
		assert(type(incus.list) == "function", "list should be a function")
		assert(type(incus.info) == "function", "info should be a function")
		assert(type(incus.delete) == "function", "delete should be a function")
	`)

	if err != nil {
		t.Errorf("Module registration failed: %v", err)
	}
}
