package systempolicy

import (
	"github.com/accuknox/knoxAutoPolicy/src/libs"
	types "github.com/accuknox/knoxAutoPolicy/src/types"
	"github.com/google/go-cmp/cmp"
)

// ==================== //
// == Exact Matching == //
// ==================== //

func IsExistingPolicy(existingPolicies []types.KubeArmorSystemPolicy, newPolicy types.KubeArmorSystemPolicy) bool {
	for _, exist := range existingPolicies {
		if exist.Metadata["clusterName"] == newPolicy.Metadata["clusterName"] &&
			exist.Metadata["namespace"] == newPolicy.Metadata["namespace"] &&
			cmp.Equal(&exist.Spec, &newPolicy.Spec) {
			return true
		}
	}

	return false
}

// ======================= //
// == Policy Name Check == //
// ======================= //

func existPolicyName(policyNamesMap map[string]bool, name string) bool {
	if _, ok := policyNamesMap[name]; ok {
		return true
	}

	return false
}

func GeneratePolicyName(policyNamesMap map[string]bool, policy types.KubeArmorSystemPolicy, clusterName string) types.KubeArmorSystemPolicy {
	procPrefix := "autopol-process-"
	filePrefix := "autopol-file-"
	netPrefix := "autopol-network-"

	polType := policy.Metadata["type"]
	name := "autopol-" + polType + "-" + libs.RandSeq(15)

	for existPolicyName(policyNamesMap, name) {
		if polType == "file" {
			name = filePrefix + libs.RandSeq(15)
		} else if polType == "process" {
			name = procPrefix + libs.RandSeq(15)
		} else { // network
			name = netPrefix + libs.RandSeq(15)
		}
	}

	policyNamesMap[name] = true
	policy.Metadata["name"] = name
	policy.Metadata["clusterName"] = clusterName

	return policy
}

// ====================================== //
// == Update Duplicated Network Policy == //
// ====================================== //

func UpdateDuplicatedPolicy(existingPolicies []types.KubeArmorSystemPolicy, discoveredPolicies []types.KubeArmorSystemPolicy, clusterName string) []types.KubeArmorSystemPolicy {
	newPolicies := []types.KubeArmorSystemPolicy{}

	// update policy name map
	policyNamesMap := map[string]bool{}
	for _, exist := range existingPolicies {
		policyNamesMap[exist.Metadata["name"]] = true
	}

	// enumerate discovered network policy
	for _, policy := range discoveredPolicies {
		// step 1: compare the total network policy spec
		if IsExistingPolicy(existingPolicies, policy) {
			continue
		}

		// step 2: generate policy name
		namedPolicy := GeneratePolicyName(policyNamesMap, policy, clusterName)
		newPolicies = append(newPolicies, namedPolicy)
	}

	return newPolicies
}