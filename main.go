/*


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

package main

import (
	"context"
	"flag"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	thatchdv1alpha1 "github.com/sergioifg94/thatchd/api/v1alpha1"
	"github.com/sergioifg94/thatchd/controllers"
	"github.com/sergioifg94/thatchd/pkg/thatchd/strategy"
	"github.com/sergioifg94/thatchd/pkg/thatchd/testsuite"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(thatchdv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "0af988fb.thatchd.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	strategyProviders := []strategy.StrategyProvider{
		&AnnotationsSuiteReconcilerProvider{},
	}

	if err = (&controllers.TestSuiteReconciler{
		Client:            mgr.GetClient(),
		Log:               ctrl.Log.WithName("controllers").WithName("TestSuite"),
		Scheme:            mgr.GetScheme(),
		StrategyProviders: strategyProviders,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TestSuite")
		os.Exit(1)
	}
	if err = (&controllers.TestCaseReconciler{
		Client:            mgr.GetClient(),
		Log:               ctrl.Log.WithName("controllers").WithName("TestCase"),
		Scheme:            mgr.GetScheme(),
		StrategyProviders: strategyProviders,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TestCase")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

type AnnotationsSuiteReconciler struct{}

type AnnotationsSuiteReconcilerProvider struct{}

var _ testsuite.Reconciler = &AnnotationsSuiteReconciler{}
var _ strategy.StrategyProvider = &AnnotationsSuiteReconcilerProvider{}

func (r *AnnotationsSuiteReconciler) Reconcile(c client.Client, namespace, currentState string) (interface{}, error) {
	ns := &v1.Namespace{}
	if err := c.Get(context.TODO(), client.ObjectKey{Name: namespace}, ns); err != nil {
		return nil, err
	}

	result := map[string]string{}
	if ns.Annotations != nil {
		result = ns.Annotations
	}

	return result, nil
}

func (p *AnnotationsSuiteReconcilerProvider) New(_ map[string]string) interface{} {
	return &AnnotationsSuiteReconciler{}
}