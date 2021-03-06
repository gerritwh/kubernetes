/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1_test

import (
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api"
	rbacapi "k8s.io/kubernetes/pkg/apis/rbac"
	_ "k8s.io/kubernetes/pkg/apis/rbac/install"
	"k8s.io/kubernetes/pkg/apis/rbac/v1alpha1"
)

func TestConversion(t *testing.T) {
	testcases := map[string]struct {
		old      *v1alpha1.RoleBinding
		expected *rbacapi.RoleBinding
	}{
		"specific user": {
			old: &v1alpha1.RoleBinding{
				RoleRef:  v1alpha1.RoleRef{Name: "foo", APIGroup: v1alpha1.GroupName},
				Subjects: []v1alpha1.Subject{{Kind: "User", Name: "bob"}},
			},
			expected: &rbacapi.RoleBinding{
				RoleRef:  rbacapi.RoleRef{Name: "foo", APIGroup: v1alpha1.GroupName},
				Subjects: []rbacapi.Subject{{Kind: "User", Name: "bob"}},
			},
		},
		"wildcard user matches authenticated": {
			old: &v1alpha1.RoleBinding{
				RoleRef:  v1alpha1.RoleRef{Name: "foo", APIGroup: v1alpha1.GroupName},
				Subjects: []v1alpha1.Subject{{Kind: "User", Name: "*"}},
			},
			expected: &rbacapi.RoleBinding{
				RoleRef:  rbacapi.RoleRef{Name: "foo", APIGroup: v1alpha1.GroupName},
				Subjects: []rbacapi.Subject{{Kind: "Group", Name: "system:authenticated"}},
			},
		},
	}
	for k, tc := range testcases {
		internal := &rbacapi.RoleBinding{}
		if err := api.Scheme.Convert(tc.old, internal, nil); err != nil {
			t.Errorf("%s: unexpected error: %v", k, err)
		}
		if !reflect.DeepEqual(internal, tc.expected) {
			t.Errorf("%s: expected\n\t%#v, got \n\t%#v", k, tc.expected, internal)
		}
	}
}
