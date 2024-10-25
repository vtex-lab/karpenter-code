package main

import (
	"encoding/json"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	kwok "sigs.k8s.io/karpenter/kwok/cloudprovider"
	v1 "sigs.k8s.io/karpenter/pkg/apis/v1"
	"sigs.k8s.io/karpenter/pkg/cloudprovider"
)

var (
	KwokZones        = []string{"test-zone-a", "test-zone-b", "test-zone-c", "test-zone-d"}
	InstanceFamilies = []string{
		"m7i-flex", "r5a", "r6g", "r7g", "vt1",
		"c3", "c5", "c5a", "c5ad", "c6a", "c6g", "c6gn", "c6i", "c7g", "i3en", "m5a", "m6g",
	}
	InstanceSizes = []string{"12xlarge", "16xlarge", "2xlarge", "3xlarge", "4xlarge", "6xlarge", "8xlarge", "large", "xlarge"}
)

func getCPUMemoryForInstance(family, size string) (int, int, int, int) {
	// Define CPU, Memory (GiB), Storage (GiB), and Pod count per instance based on AWS EC2 guidelines
	switch family {
	case "m7i-flex":
		switch size {
		case "12xlarge": return 48, 192, 300, 768
		case "16xlarge": return 64, 256, 400, 1024
		case "2xlarge": return 8, 32, 100, 128
		case "3xlarge": return 12, 48, 150, 192
		case "4xlarge": return 16, 64, 200, 256
		case "6xlarge": return 24, 96, 300, 384
		case "8xlarge": return 32, 128, 400, 512
		case "large": return 2, 8, 50, 32
		case "xlarge": return 4, 16, 100, 64
		}

	case "r5a":
		switch size {
		case "12xlarge": return 48, 384, 300, 768
		case "16xlarge": return 64, 512, 400, 1024
		case "2xlarge": return 8, 64, 100, 128
		case "3xlarge": return 12, 96, 150, 192
		case "4xlarge": return 16, 128, 200, 256
		case "6xlarge": return 24, 192, 300, 384
		case "8xlarge": return 32, 256, 400, 512
		case "large": return 2, 16, 50, 32
		case "xlarge": return 4, 32, 100, 64
		}

	case "r6g":
		switch size {
		case "12xlarge": return 48, 384, 300, 768
		case "16xlarge": return 64, 512, 400, 1024
		case "2xlarge": return 8, 64, 100, 128
		case "3xlarge": return 12, 96, 150, 192
		case "4xlarge": return 16, 128, 200, 256
		case "6xlarge": return 24, 192, 300, 384
		case "8xlarge": return 32, 256, 400, 512
		case "large": return 2, 16, 50, 32
		case "xlarge": return 4, 32, 100, 64
		}

	case "r7g":
		switch size {
		case "12xlarge": return 48, 384, 300, 768
		case "16xlarge": return 64, 512, 400, 1024
		case "2xlarge": return 8, 64, 100, 128
		case "3xlarge": return 12, 96, 150, 192
		case "4xlarge": return 16, 128, 200, 256
		case "6xlarge": return 24, 192, 300, 384
		case "8xlarge": return 32, 256, 400, 512
		case "large": return 2, 16, 50, 32
		case "xlarge": return 4, 32, 100, 64
		}

	case "vt1":
		switch size {
		case "12xlarge": return 48, 192, 300, 768
		case "16xlarge": return 64, 256, 400, 1024
		case "2xlarge": return 8, 32, 100, 128
		case "3xlarge": return 12, 48, 150, 192
		case "4xlarge": return 16, 64, 200, 256
		case "6xlarge": return 24, 96, 300, 384
		case "8xlarge": return 32, 128, 400, 512
		case "large": return 2, 8, 50, 32
		case "xlarge": return 4, 16, 100, 64
		}

	case "c3":
		switch size {
		case "12xlarge": return 48, 96, 300, 768
		case "16xlarge": return 64, 128, 400, 1024
		case "2xlarge": return 8, 16, 100, 128
		case "3xlarge": return 12, 24, 150, 192
		case "4xlarge": return 16, 32, 200, 256
		case "6xlarge": return 24, 48, 300, 384
		case "8xlarge": return 32, 64, 400, 512
		case "large": return 2, 4, 50, 32
		case "xlarge": return 4, 8, 100, 64
		}

	case "c5":
		switch size {
		case "12xlarge": return 48, 96, 300, 768
		case "16xlarge": return 64, 128, 400, 1024
		case "2xlarge": return 8, 16, 100, 128
		case "3xlarge": return 12, 24, 150, 192
		case "4xlarge": return 16, 32, 200, 256
		case "6xlarge": return 24, 48, 300, 384
		case "8xlarge": return 32, 64, 400, 512
		case "large": return 2, 4, 50, 32
		case "xlarge": return 4, 8, 100, 64
		}

	case "c5a":
		switch size {
		case "12xlarge": return 48, 96, 300, 768
		case "16xlarge": return 64, 128, 400, 1024
		case "2xlarge": return 8, 16, 100, 128
		case "3xlarge": return 12, 24, 150, 192
		case "4xlarge": return 16, 32, 200, 256
		case "6xlarge": return 24, 48, 300, 384
		case "8xlarge": return 32, 64, 400, 512
		case "large": return 2, 4, 50, 32
		case "xlarge": return 4, 8, 100, 64
		}

	case "c5ad":
		switch size {
		case "12xlarge": return 48, 96, 300, 768
		case "16xlarge": return 64, 128, 400, 1024
		case "2xlarge": return 8, 16, 100, 128
		case "3xlarge": return 12, 24, 150, 192
		case "4xlarge": return 16, 32, 200, 256
		case "6xlarge": return 24, 48, 300, 384
		case "8xlarge": return 32, 64, 400, 512
		case "large": return 2, 4, 50, 32
		case "xlarge": return 4, 8, 100, 64
		}

	case "c6a":
		switch size {
		case "12xlarge": return 48, 96, 300, 768
		case "16xlarge": return 64, 128, 400, 1024
		case "2xlarge": return 8, 16, 100, 128
		case "3xlarge": return 12, 24, 150, 192
		case "4xlarge": return 16, 32, 200, 256
		case "6xlarge": return 24, 48, 300, 384
		case "8xlarge": return 32, 64, 400, 512
		case "large": return 2, 4, 50, 32
		case "xlarge": return 4, 8, 100, 64
        }
    case "c6g":
		switch size {
		case "12xlarge": return 48, 384, 300, 768
		case "16xlarge": return 64, 512, 400, 1024
		case "2xlarge": return 8, 64, 100, 128
		case "3xlarge": return 12, 96, 150, 192
		case "4xlarge": return 16, 128, 200, 256
		case "6xlarge": return 24, 192, 300, 384
		case "8xlarge": return 32, 256, 400, 512
		case "large": return 2, 16, 50, 32
		case "xlarge": return 4, 32, 100, 64
		}

	case "c6gn":
		switch size {
		case "12xlarge": return 48, 384, 300, 768
		case "16xlarge": return 64, 512, 400, 1024
		case "2xlarge": return 8, 64, 100, 128
		case "3xlarge": return 12, 96, 150, 192
		case "4xlarge": return 16, 128, 200, 256
		case "6xlarge": return 24, 192, 300, 384
		case "8xlarge": return 32, 256, 400, 512
		case "large": return 2, 16, 50, 32
		case "xlarge": return 4, 32, 100, 64
		}

	case "c6i":
		switch size {
		case "12xlarge": return 48, 96, 300, 768
		case "16xlarge": return 64, 128, 400, 1024
		case "2xlarge": return 8, 16, 100, 128
		case "3xlarge": return 12, 24, 150, 192
		case "4xlarge": return 16, 32, 200, 256
		case "6xlarge": return 24, 48, 300, 384
		case "8xlarge": return 32, 64, 400, 512
		case "large": return 2, 4, 50, 32
		case "xlarge": return 4, 8, 100, 64
		}

	case "c7g":
		switch size {
		case "12xlarge": return 48, 384, 300, 768
		case "16xlarge": return 64, 512, 400, 1024
		case "2xlarge": return 8, 64, 100, 128
		case "3xlarge": return 12, 96, 150, 192
		case "4xlarge": return 16, 128, 200, 256
		case "6xlarge": return 24, 192, 300, 384
		case "8xlarge": return 32, 256, 400, 512
		case "large": return 2, 16, 50, 32
		case "xlarge": return 4, 32, 100, 64
		}

	case "i3en":
		switch size {
		case "12xlarge": return 48, 384, 7500, 768 // 7.5 TB storage
		case "16xlarge": return 64, 512, 10000, 1024 // 10 TB storage
		case "2xlarge": return 8, 64, 2500, 128 // 2.5 TB storage
		case "3xlarge": return 12, 96, 3750, 192 // 3.75 TB storage
		case "4xlarge": return 16, 128, 5000, 256 // 5 TB storage
		case "6xlarge": return 24, 192, 7500, 384 // 7.5 TB storage
		case "8xlarge": return 32, 256, 10000, 512 // 10 TB storage
		case "large": return 2, 16, 1250, 32 // 1.25 TB storage
		case "xlarge": return 4, 32, 2500, 64 // 2.5 TB storage
		}

	case "m5a":
		switch size {
		case "12xlarge": return 48, 192, 300, 768
		case "16xlarge": return 64, 256, 400, 1024
		case "2xlarge": return 8, 32, 100, 128
		case "3xlarge": return 12, 48, 150, 192
		case "4xlarge": return 16, 64, 200, 256
		case "6xlarge": return 24, 96, 300, 384
		case "8xlarge": return 32, 128, 400, 512
		case "large": return 2, 8, 50, 32
		case "xlarge": return 4, 16, 100, 64
		}

	case "m6g":
		switch size {
		case "12xlarge": return 48, 192, 300, 768
		case "16xlarge": return 64, 256, 400, 1024
		case "2xlarge": return 8, 32, 100, 128
		case "3xlarge": return 12, 48, 150, 192
		case "4xlarge": return 16, 64, 200, 256
		case "6xlarge": return 24, 96, 300, 384
		case "8xlarge": return 32, 128, 400, 512
		case "large": return 2, 8, 50, 32
		case "xlarge": return 4, 16, 100, 64
		}

    }
	return 0, 0, 0, 0
}

func priceFromResources(resources corev1.ResourceList) float64 {
	price := 0.0
	for k, v := range resources {
		switch k {
		case corev1.ResourceCPU:
			price += 0.025 * v.AsApproximateFloat64()
		case corev1.ResourceMemory:
			price += 0.001 * v.AsApproximateFloat64() / (1e9) // Convert bytes to GiB
		}
	}
	return price
}

func constructInstanceTypes() []kwok.InstanceTypeOptions {
	var instanceTypesOptions []kwok.InstanceTypeOptions

	for _, family := range InstanceFamilies {
		for _, size := range InstanceSizes {
			cpu, mem, storage, pods := getCPUMemoryForInstance(family, size)
			if cpu == 0 || mem == 0 {
				continue // Skip undefined instances
			}

			// Ensure that only Linux-based instances are generated for simplicity
			for _, arch := range []string{v1.ArchitectureAmd64, v1.ArchitectureArm64} {
				opts := kwok.InstanceTypeOptions{
					Name:             fmt.Sprintf("%s.%s", family, size),
					Architecture:     arch,
					OperatingSystems: []corev1.OSName{corev1.Linux},
					Resources: corev1.ResourceList{
						corev1.ResourceCPU:              resource.MustParse(fmt.Sprintf("%d", cpu)),
						corev1.ResourceMemory:           resource.MustParse(fmt.Sprintf("%dGi", mem)),
						corev1.ResourcePods:             resource.MustParse(fmt.Sprintf("%d", pods)),
						corev1.ResourceEphemeralStorage: resource.MustParse(fmt.Sprintf("%dGi", storage)),
					},
				}
				price := priceFromResources(opts.Resources)

				opts.Offerings = []kwok.KWOKOffering{}
				for _, zone := range KwokZones {
					opts.Offerings = append(opts.Offerings, kwok.KWOKOffering{
						Requirements: []corev1.NodeSelectorRequirement{
							{
								Key:      v1.CapacityTypeLabelKey,
								Operator: corev1.NodeSelectorOpIn,
								Values:   []string{v1.CapacityTypeOnDemand}, // Only on-demand instances
							},
							{
								Key:      corev1.LabelTopologyZone,
								Operator: corev1.NodeSelectorOpIn,
								Values:   []string{zone},
							},
						},
						Offering: cloudprovider.Offering{
							Price:     price,
							Available: true,
						},
					})
				}
				instanceTypesOptions = append(instanceTypesOptions, opts)
			}
		}
	}
	return instanceTypesOptions
}

func main() {
	opts := constructInstanceTypes()
	output, err := json.MarshalIndent(opts, "", "    ")
	if err != nil {
		fmt.Printf("could not marshal generated instance types to JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(string(output))
}

