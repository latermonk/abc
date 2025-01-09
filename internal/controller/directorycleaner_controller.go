/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cleanupv1 "example.com/directory-cleaner-operator/api/v1"
)

// DirectoryCleanerReconciler reconciles a DirectoryCleaner object
type DirectoryCleanerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cleanup.example.com,resources=directorycleaners,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cleanup.example.com,resources=directorycleaners/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cleanup.example.com,resources=directorycleaners/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DirectoryCleaner object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *DirectoryCleanerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the DirectoryCleaner instance
	var directoryCleaner cleanupv1.DirectoryCleaner
	if err := r.Get(ctx, req.NamespacedName, &directoryCleaner); err != nil {
		log.Error(err, "unable to fetch DirectoryCleaner")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the node name where this pod is running
	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		return ctrl.Result{}, fmt.Errorf("NODE_NAME environment variable not set")
	}

	// Check if this node should handle the cleaning
	if !isNodeMatch(nodeName, directoryCleaner.Spec.NodeSelector) {
		return ctrl.Result{}, nil
	}

	// Perform the cleaning operation
	cleanedPaths := []string{}
	for _, path := range directoryCleaner.Spec.Paths {
		err := cleanPath(path, directoryCleaner.Spec.Recursive, directoryCleaner.Spec.Force)
		if err != nil {
			log.Error(err, "Failed to clean path", "path", path)
			continue
		}
		cleanedPaths = append(cleanedPaths, path)
	}

	// Update status
	directoryCleaner.Status.LastCleanedTime = metav1.Now()
	directoryCleaner.Status.LastCleanedPaths = cleanedPaths
	if err := r.Status().Update(ctx, &directoryCleaner); err != nil {
		log.Error(err, "Failed to update DirectoryCleaner status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func isNodeMatch(nodeName string, selector map[string]string) bool {
	// TODO: Implement node selector matching logic
	return true
}

func cleanPath(path string, recursive bool, force bool) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.IsDir() {
		if recursive {
			return os.RemoveAll(path)
		}
		return os.Remove(path)
	}

	if force {
		err = os.Chmod(path, 0666)
		if err != nil {
			return err
		}
	}
	return os.Remove(path)
}

func (r *DirectoryCleanerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cleanupv1.DirectoryCleaner{}).
		Complete(r)
}
