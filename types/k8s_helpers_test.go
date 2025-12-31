package types

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestNewToleration(t *testing.T) {
	tests := []struct {
		name     string
		nodepool string
		validate func(*testing.T, []apiv1.Toleration)
	}{
		{
			name:     "standard nodepool",
			nodepool: "gke1",
			validate: func(t *testing.T, tolerations []apiv1.Toleration) {
				if len(tolerations) != 3 {
					t.Fatalf("Expected 3 tolerations, got %d", len(tolerations))
				}
				// Check first two default tolerations
				if tolerations[0].Key != "node.kubernetes.io/not-ready" {
					t.Errorf("Expected first toleration key 'node.kubernetes.io/not-ready', got %q", tolerations[0].Key)
				}
				if tolerations[1].Key != "node.kubernetes.io/unreachable" {
					t.Errorf("Expected second toleration key 'node.kubernetes.io/unreachable', got %q", tolerations[1].Key)
				}
				// Check nodepool-specific toleration
				if tolerations[2].Key != "dedicated" {
					t.Errorf("Expected third toleration key 'dedicated', got %q", tolerations[2].Key)
				}
				if tolerations[2].Value != "gke1" {
					t.Errorf("Expected third toleration value 'gke1', got %q", tolerations[2].Value)
				}
				if tolerations[2].Operator != apiv1.TolerationOperator("Equal") {
					t.Errorf("Expected third toleration operator 'Equal', got %q", tolerations[2].Operator)
				}
			},
		},
		{
			name:     "empty nodepool",
			nodepool: "",
			validate: func(t *testing.T, tolerations []apiv1.Toleration) {
				if len(tolerations) != 3 {
					t.Fatalf("Expected 3 tolerations, got %d", len(tolerations))
				}
				if tolerations[2].Value != "" {
					t.Errorf("Expected third toleration value '', got %q", tolerations[2].Value)
				}
			},
		},
		{
			name:     "special characters in nodepool",
			nodepool: "gke-1-prod",
			validate: func(t *testing.T, tolerations []apiv1.Toleration) {
				if tolerations[2].Value != "gke-1-prod" {
					t.Errorf("Expected third toleration value 'gke-1-prod', got %q", tolerations[2].Value)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newToleration(tt.nodepool)
			if result == nil {
				t.Fatal("newToleration returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestNewVolumeMount(t *testing.T) {
	tests := []struct {
		name         string
		application  string
		dependencies []DependsOn
		validate     func(*testing.T, []apiv1.VolumeMount)
	}{
		{
			name:        "config dependency",
			application: "myapp",
			dependencies: []DependsOn{
				{Name: "redis-config"},
			},
			validate: func(t *testing.T, mounts []apiv1.VolumeMount) {
				if len(mounts) != 1 {
					t.Fatalf("Expected 1 volume mount, got %d", len(mounts))
				}
				if mounts[0].Name != "redis" {
					t.Errorf("Expected mount name 'redis', got %q", mounts[0].Name)
				}
				if mounts[0].MountPath != "/app/conf/redis.conf" {
					t.Errorf("Expected mount path '/app/conf/redis.conf', got %q", mounts[0].MountPath)
				}
				if mounts[0].SubPath != "redis.conf" {
					t.Errorf("Expected sub path 'redis.conf', got %q", mounts[0].SubPath)
				}
			},
		},
		{
			name:        "GoogleCloudStorage dependency",
			application: "myapp",
			dependencies: []DependsOn{
				{Name: "GoogleCloudStorage"},
			},
			validate: func(t *testing.T, mounts []apiv1.VolumeMount) {
				if len(mounts) != 1 {
					t.Fatalf("Expected 1 volume mount, got %d", len(mounts))
				}
				if mounts[0].Name != "google-cloud-myapp" {
					t.Errorf("Expected mount name 'google-cloud-myapp', got %q", mounts[0].Name)
				}
				if !mounts[0].ReadOnly {
					t.Error("Expected ReadOnly to be true")
				}
				if mounts[0].MountPath != "/google-cloud-myapp" {
					t.Errorf("Expected mount path '/google-cloud-myapp', got %q", mounts[0].MountPath)
				}
			},
		},
		{
			name:        "maxmind dependency",
			application: "myapp",
			dependencies: []DependsOn{
				{Name: "maxmind"},
			},
			validate: func(t *testing.T, mounts []apiv1.VolumeMount) {
				if len(mounts) != 1 {
					t.Fatalf("Expected 1 volume mount, got %d", len(mounts))
				}
				if mounts[0].Name != "geoip-files" {
					t.Errorf("Expected mount name 'geoip-files', got %q", mounts[0].Name)
				}
				if mounts[0].MountPath != "/usr/share/GeoIP" {
					t.Errorf("Expected mount path '/usr/share/GeoIP', got %q", mounts[0].MountPath)
				}
			},
		},
		{
			name:        "multiple dependencies",
			application: "myapp",
			dependencies: []DependsOn{
				{Name: "redis-config"},
				{Name: "GoogleCloudStorage"},
				{Name: "maxmind"},
			},
			validate: func(t *testing.T, mounts []apiv1.VolumeMount) {
				if len(mounts) != 3 {
					t.Fatalf("Expected 3 volume mounts, got %d", len(mounts))
				}
				// Should be sorted by name
				if mounts[0].Name > mounts[1].Name || mounts[1].Name > mounts[2].Name {
					t.Error("Expected volume mounts to be sorted by name")
				}
			},
		},
		{
			name:        "no matching dependencies",
			application: "myapp",
			dependencies: []DependsOn{
				{Name: "redis"},
				{Name: "postgres"},
			},
			validate: func(t *testing.T, mounts []apiv1.VolumeMount) {
				if len(mounts) != 0 {
					t.Errorf("Expected 0 volume mounts, got %d", len(mounts))
				}
			},
		},
		{
			name:         "empty dependencies",
			application:  "myapp",
			dependencies: []DependsOn{},
			validate: func(t *testing.T, mounts []apiv1.VolumeMount) {
				if len(mounts) != 0 {
					t.Errorf("Expected 0 volume mounts, got %d", len(mounts))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newVolumeMount(tt.application, tt.dependencies)
			if result == nil {
				t.Fatal("newVolumeMount returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestNewVolume(t *testing.T) {
	tests := []struct {
		name         string
		application  string
		dependencies []DependsOn
		validate     func(*testing.T, []apiv1.Volume)
	}{
		{
			name:        "config dependency",
			application: "myapp",
			dependencies: []DependsOn{
				{Name: "redis-config"},
			},
			validate: func(t *testing.T, volumes []apiv1.Volume) {
				if len(volumes) != 1 {
					t.Fatalf("Expected 1 volume, got %d", len(volumes))
				}
				if volumes[0].Name != "redis" {
					t.Errorf("Expected volume name 'redis', got %q", volumes[0].Name)
				}
				if volumes[0].ConfigMap == nil {
					t.Fatal("Expected ConfigMap volume source")
				}
				if volumes[0].ConfigMap.Name != "redis" {
					t.Errorf("Expected ConfigMap name 'redis', got %q", volumes[0].ConfigMap.Name)
				}
			},
		},
		{
			name:        "GoogleCloudStorage dependency",
			application: "myapp",
			dependencies: []DependsOn{
				{Name: "GoogleCloudStorage"},
			},
			validate: func(t *testing.T, volumes []apiv1.Volume) {
				if len(volumes) != 1 {
					t.Fatalf("Expected 1 volume, got %d", len(volumes))
				}
				if volumes[0].Name != "google-cloud-myapp" {
					t.Errorf("Expected volume name 'google-cloud-myapp', got %q", volumes[0].Name)
				}
				if volumes[0].Secret == nil {
					t.Fatal("Expected Secret volume source")
				}
				if volumes[0].Secret.SecretName != "google-cloud-myapp" {
					t.Errorf("Expected Secret name 'google-cloud-myapp', got %q", volumes[0].Secret.SecretName)
				}
			},
		},
		{
			name:        "maxmind dependency",
			application: "myapp",
			dependencies: []DependsOn{
				{Name: "maxmind"},
			},
			validate: func(t *testing.T, volumes []apiv1.Volume) {
				if len(volumes) != 1 {
					t.Fatalf("Expected 1 volume, got %d", len(volumes))
				}
				if volumes[0].Name != "geoip-files" {
					t.Errorf("Expected volume name 'geoip-files', got %q", volumes[0].Name)
				}
				if volumes[0].EmptyDir == nil {
					t.Fatal("Expected EmptyDir volume source")
				}
			},
		},
		{
			name:        "multiple dependencies",
			application: "myapp",
			dependencies: []DependsOn{
				{Name: "redis-config"},
				{Name: "GoogleCloudStorage"},
				{Name: "maxmind"},
			},
			validate: func(t *testing.T, volumes []apiv1.Volume) {
				if len(volumes) != 3 {
					t.Fatalf("Expected 3 volumes, got %d", len(volumes))
				}
				// Should be sorted by name
				if volumes[0].Name > volumes[1].Name || volumes[1].Name > volumes[2].Name {
					t.Error("Expected volumes to be sorted by name")
				}
			},
		},
		{
			name:        "no matching dependencies",
			application: "myapp",
			dependencies: []DependsOn{
				{Name: "redis"},
				{Name: "postgres"},
			},
			validate: func(t *testing.T, volumes []apiv1.Volume) {
				if len(volumes) != 0 {
					t.Errorf("Expected 0 volumes, got %d", len(volumes))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newVolume(tt.application, tt.dependencies)
			if result == nil {
				t.Fatal("newVolume returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestNewResourceRequirements(t *testing.T) {
	tests := []struct {
		name     string
		tier     *Datacenter
		validate func(*testing.T, apiv1.ResourceRequirements)
	}{
		{
			name: "with limits",
			tier: &Datacenter{
				Resources: &ResourceRequirements{
					Limits: &ResourceList{
						CPU:    "1000m",
						Memory: "2Gi",
					},
					Requests: &ResourceList{
						CPU:    "500m",
						Memory: "1Gi",
					},
				},
			},
			validate: func(t *testing.T, req apiv1.ResourceRequirements) {
				if req.Limits == nil {
					t.Fatal("Expected Limits to be set")
				}
				if req.Limits[apiv1.ResourceCPU] != resource.MustParse("1000m") {
					t.Errorf("Expected CPU limit 1000m, got %v", req.Limits[apiv1.ResourceCPU])
				}
				if req.Limits[apiv1.ResourceMemory] != resource.MustParse("2Gi") {
					t.Errorf("Expected Memory limit 2Gi, got %v", req.Limits[apiv1.ResourceMemory])
				}
				if req.Requests[apiv1.ResourceCPU] != resource.MustParse("500m") {
					t.Errorf("Expected CPU request 500m, got %v", req.Requests[apiv1.ResourceCPU])
				}
				if req.Requests[apiv1.ResourceMemory] != resource.MustParse("1Gi") {
					t.Errorf("Expected Memory request 1Gi, got %v", req.Requests[apiv1.ResourceMemory])
				}
			},
		},
		{
			name: "without limits uses requests",
			tier: &Datacenter{
				Resources: &ResourceRequirements{
					Limits: nil,
					Requests: &ResourceList{
						CPU:    "500m",
						Memory: "1Gi",
					},
				},
			},
			validate: func(t *testing.T, req apiv1.ResourceRequirements) {
				if req.Limits == nil {
					t.Fatal("Expected Limits to be set from Requests")
				}
				if req.Limits[apiv1.ResourceCPU] != resource.MustParse("500m") {
					t.Errorf("Expected CPU limit 500m (from requests), got %v", req.Limits[apiv1.ResourceCPU])
				}
				if req.Limits[apiv1.ResourceMemory] != resource.MustParse("1Gi") {
					t.Errorf("Expected Memory limit 1Gi (from requests), got %v", req.Limits[apiv1.ResourceMemory])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newResourceRequirements(tt.tier)
			tt.validate(t, result)
		})
	}
}

func TestNewMetadataLabels(t *testing.T) {
	tests := []struct {
		name        string
		application string
		chaosConfig *ChaosMonkey
		validate    func(*testing.T, map[string]string)
	}{
		{
			name:        "complete chaos config",
			application: "myapp",
			chaosConfig: &ChaosMonkey{
				KillMode:  "fixed",
				KillValue: "1",
				MTBF:      "24h",
			},
			validate: func(t *testing.T, labels map[string]string) {
				if labels["kube-monkey/enabled"] != "enabled" {
					t.Errorf("Expected kube-monkey/enabled 'enabled', got %q", labels["kube-monkey/enabled"])
				}
				if labels["kube-monkey/identifier"] != "myapp" {
					t.Errorf("Expected kube-monkey/identifier 'myapp', got %q", labels["kube-monkey/identifier"])
				}
				if labels["kube-monkey/kill-mode"] != "fixed" {
					t.Errorf("Expected kube-monkey/kill-mode 'fixed', got %q", labels["kube-monkey/kill-mode"])
				}
				if labels["kube-monkey/kill-value"] != "1" {
					t.Errorf("Expected kube-monkey/kill-value '1', got %q", labels["kube-monkey/kill-value"])
				}
				if labels["kube-monkey/mtbf"] != "24h" {
					t.Errorf("Expected kube-monkey/mtbf '24h', got %q", labels["kube-monkey/mtbf"])
				}
			},
		},
		{
			name:        "empty application",
			application: "",
			chaosConfig: &ChaosMonkey{
				KillMode:  "fixed",
				KillValue: "1",
				MTBF:      "24h",
			},
			validate: func(t *testing.T, labels map[string]string) {
				if labels["kube-monkey/identifier"] != "" {
					t.Errorf("Expected kube-monkey/identifier '', got %q", labels["kube-monkey/identifier"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newMetadataLabels(tt.application, tt.chaosConfig)
			if result == nil {
				t.Fatal("newMetadataLabels returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestNewTemplateLabels(t *testing.T) {
	tests := []struct {
		name        string
		application string
		validate    func(*testing.T, map[string]string)
	}{
		{
			name:        "standard application",
			application: "myapp",
			validate: func(t *testing.T, labels map[string]string) {
				if labels["app"] != "myapp" {
					t.Errorf("Expected app 'myapp', got %q", labels["app"])
				}
				if labels["kube-monkey/enabled"] != "enabled" {
					t.Errorf("Expected kube-monkey/enabled 'enabled', got %q", labels["kube-monkey/enabled"])
				}
				if labels["kube-monkey/identifier"] != "myapp" {
					t.Errorf("Expected kube-monkey/identifier 'myapp', got %q", labels["kube-monkey/identifier"])
				}
			},
		},
		{
			name:        "empty application",
			application: "",
			validate: func(t *testing.T, labels map[string]string) {
				if labels["app"] != "" {
					t.Errorf("Expected app '', got %q", labels["app"])
				}
				if labels["kube-monkey/identifier"] != "" {
					t.Errorf("Expected kube-monkey/identifier '', got %q", labels["kube-monkey/identifier"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newTemplateLabels(tt.application)
			if result == nil {
				t.Fatal("newTemplateLabels returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestNewAffinity(t *testing.T) {
	tests := []struct {
		name        string
		application string
		nodepool    string
		validate    func(*testing.T, *apiv1.Affinity)
	}{
		{
			name:        "complete affinity",
			application: "myapp",
			nodepool:    "gke1",
			validate: func(t *testing.T, affinity *apiv1.Affinity) {
				if affinity.NodeAffinity == nil {
					t.Fatal("Expected NodeAffinity to be set")
				}
				if affinity.PodAntiAffinity == nil {
					t.Fatal("Expected PodAntiAffinity to be set")
				}
				if len(affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms) != 1 {
					t.Fatalf("Expected 1 NodeSelectorTerm, got %d", len(affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms))
				}
				nodeSelector := affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0]
				if len(nodeSelector.MatchExpressions) != 1 {
					t.Fatalf("Expected 1 MatchExpression, got %d", len(nodeSelector.MatchExpressions))
				}
				if nodeSelector.MatchExpressions[0].Key != "dedicated" {
					t.Errorf("Expected key 'dedicated', got %q", nodeSelector.MatchExpressions[0].Key)
				}
				if len(nodeSelector.MatchExpressions[0].Values) != 1 || nodeSelector.MatchExpressions[0].Values[0] != "gke1" {
					t.Errorf("Expected values ['gke1'], got %v", nodeSelector.MatchExpressions[0].Values)
				}
				if len(affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution) != 1 {
					t.Fatalf("Expected 1 PodAffinityTerm, got %d", len(affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution))
				}
				podAffinityTerm := affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0]
				if podAffinityTerm.TopologyKey != "kubernetes.io/hostname" {
					t.Errorf("Expected TopologyKey 'kubernetes.io/hostname', got %q", podAffinityTerm.TopologyKey)
				}
				if podAffinityTerm.LabelSelector == nil {
					t.Fatal("Expected LabelSelector to be set")
				}
				if len(podAffinityTerm.LabelSelector.MatchExpressions) != 1 {
					t.Fatalf("Expected 1 MatchExpression, got %d", len(podAffinityTerm.LabelSelector.MatchExpressions))
				}
				if podAffinityTerm.LabelSelector.MatchExpressions[0].Key != "app" {
					t.Errorf("Expected key 'app', got %q", podAffinityTerm.LabelSelector.MatchExpressions[0].Key)
				}
				if len(podAffinityTerm.LabelSelector.MatchExpressions[0].Values) != 1 || podAffinityTerm.LabelSelector.MatchExpressions[0].Values[0] != "myapp" {
					t.Errorf("Expected values ['myapp'], got %v", podAffinityTerm.LabelSelector.MatchExpressions[0].Values)
				}
			},
		},
		{
			name:        "empty nodepool",
			application: "myapp",
			nodepool:    "",
			validate: func(t *testing.T, affinity *apiv1.Affinity) {
				nodeSelector := affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0]
				if len(nodeSelector.MatchExpressions[0].Values) != 1 || nodeSelector.MatchExpressions[0].Values[0] != "" {
					t.Errorf("Expected values [''], got %v", nodeSelector.MatchExpressions[0].Values)
				}
			},
		},
		{
			name:        "empty application",
			application: "",
			nodepool:    "gke1",
			validate: func(t *testing.T, affinity *apiv1.Affinity) {
				podAffinityTerm := affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0]
				if len(podAffinityTerm.LabelSelector.MatchExpressions[0].Values) != 1 || podAffinityTerm.LabelSelector.MatchExpressions[0].Values[0] != "" {
					t.Errorf("Expected values [''], got %v", podAffinityTerm.LabelSelector.MatchExpressions[0].Values)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newAffinity(tt.application, tt.nodepool)
			if result == nil {
				t.Fatal("newAffinity returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestDefaultAnnotations(t *testing.T) {
	tests := []struct {
		name         string
		organization string
		stage        string
		owner        string
		validate     func(*testing.T, map[string]string)
	}{
		{
			name:         "complete annotations",
			organization: "myorg",
			stage:        "production",
			owner:        "john@example.com",
			validate: func(t *testing.T, annotations map[string]string) {
				if annotations["moniker.spinnaker.io/stack"] != "production" {
					t.Errorf("Expected moniker.spinnaker.io/stack 'production', got %q", annotations["moniker.spinnaker.io/stack"])
				}
				if annotations["service.myorg.dev/generated"] != "spini/v1" {
					t.Errorf("Expected service.myorg.dev/generated 'spini/v1', got %q", annotations["service.myorg.dev/generated"])
				}
				if annotations["service.myorg.dev/owners"] != "john@example.com" {
					t.Errorf("Expected service.myorg.dev/owners 'john@example.com', got %q", annotations["service.myorg.dev/owners"])
				}
			},
		},
		{
			name:         "empty organization",
			organization: "",
			stage:        "production",
			owner:        "john@example.com",
			validate: func(t *testing.T, annotations map[string]string) {
				if annotations["service..dev/generated"] != "spini/v1" {
					t.Errorf("Expected service..dev/generated 'spini/v1', got %q", annotations["service..dev/generated"])
				}
			},
		},
		{
			name:         "empty stage",
			organization: "myorg",
			stage:        "",
			owner:        "john@example.com",
			validate: func(t *testing.T, annotations map[string]string) {
				if annotations["moniker.spinnaker.io/stack"] != "" {
					t.Errorf("Expected moniker.spinnaker.io/stack '', got %q", annotations["moniker.spinnaker.io/stack"])
				}
			},
		},
		{
			name:         "empty owner",
			organization: "myorg",
			stage:        "production",
			owner:        "",
			validate: func(t *testing.T, annotations map[string]string) {
				if annotations["service.myorg.dev/owners"] != "" {
					t.Errorf("Expected service.myorg.dev/owners '', got %q", annotations["service.myorg.dev/owners"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultAnnotations(tt.organization, tt.stage, tt.owner)
			if result == nil {
				t.Fatal("defaultAnnotations returned nil")
			}
			tt.validate(t, result)
		})
	}
}
