// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/toddtreece/kubernetes-list-tests/informer/pkg/apis/dashboard/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DashboardLister helps list Dashboards.
// All objects returned here must be treated as read-only.
type DashboardLister interface {
	// List lists all Dashboards in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Dashboard, err error)
	// Dashboards returns an object that can list and get Dashboards.
	Dashboards(namespace string) DashboardNamespaceLister
	DashboardListerExpansion
}

// dashboardLister implements the DashboardLister interface.
type dashboardLister struct {
	indexer cache.Indexer
}

// NewDashboardLister returns a new DashboardLister.
func NewDashboardLister(indexer cache.Indexer) DashboardLister {
	return &dashboardLister{indexer: indexer}
}

// List lists all Dashboards in the indexer.
func (s *dashboardLister) List(selector labels.Selector) (ret []*v1alpha1.Dashboard, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Dashboard))
	})
	return ret, err
}

// Dashboards returns an object that can list and get Dashboards.
func (s *dashboardLister) Dashboards(namespace string) DashboardNamespaceLister {
	return dashboardNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DashboardNamespaceLister helps list and get Dashboards.
// All objects returned here must be treated as read-only.
type DashboardNamespaceLister interface {
	// List lists all Dashboards in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Dashboard, err error)
	// Get retrieves the Dashboard from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Dashboard, error)
	DashboardNamespaceListerExpansion
}

// dashboardNamespaceLister implements the DashboardNamespaceLister
// interface.
type dashboardNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Dashboards in the indexer for a given namespace.
func (s dashboardNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Dashboard, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Dashboard))
	})
	return ret, err
}

// Get retrieves the Dashboard from the indexer for a given namespace and name.
func (s dashboardNamespaceLister) Get(name string) (*v1alpha1.Dashboard, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("dashboard"), name)
	}
	return obj.(*v1alpha1.Dashboard), nil
}
