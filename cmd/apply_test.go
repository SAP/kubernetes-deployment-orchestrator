package cmd

import (
	"bytes"
	"path"
	"path/filepath"
	"runtime"

	semver "github.com/Masterminds/semver/v3"
	"github.com/sap/kubernetes-deployment-orchestrator/pkg/k8s"
	"github.com/sap/kubernetes-deployment-orchestrator/pkg/kdo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fake_k8s_test.go ../pkg/kdo K8s

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
	root       = path.Join(filepath.Dir(b), "..")
	example    = path.Join(root, "charts", "example", "simple")
)

var _ = Describe("Apply Chart", func() {

	It("produces the correct output", func() {
		Skip("unsupported")
		writer := bytes.Buffer{}
		k := &k8s.FakeK8s{
			ApplyStub: func(i k8s.ObjectStream, options *k8s.Options) error {
				return i.Encode()(&writer)
			},
		}
		k.ForSubChartStub = func(s string, app string, version *semver.Version, children int) k8s.K8s {
			return k
		}
		k.GetStub = func(s string, s2 string, options *k8s.Options) (*k8s.Object, error) {
			return &k8s.Object{}, nil
		}

		err := apply(path.Join(example, "cf"), k, kdo.WithNamespace("mynamespace"))
		Expect(err).ToNot(HaveOccurred())
		output := writer.String()
		Expect(output).To(ContainSubstring("CREATE OR REPLACE USER 'uaa'"))
		Expect(k.RolloutStatusCallCount()).To(Equal(1))
		Expect(k.ApplyCallCount()).To(Equal(3))
		Expect(k.ForSubChartCallCount()).To(Equal(3))
		namespace, _, _, _ := k.ForSubChartArgsForCall(0)
		Expect(namespace).To(Equal("mynamespace"))
		namespace, _, _, _ = k.ForSubChartArgsForCall(1)
		Expect(namespace).To(Equal("mynamespace"))
		namespace, _, _, _ = k.ForSubChartArgsForCall(2)
		Expect(namespace).To(Equal("uaa"))
		kind, name, _ := k.RolloutStatusArgsForCall(0)
		Expect(name).To(Equal("uaa-master"))
		Expect(kind).To(Equal("statefulset"))
	})

	It("produces correct objects", func() {
		Skip("unsupported")
		k := k8s.NewK8sInMemory("default")
		err := apply(path.Join(example, "cf"), k, kdo.WithNamespace("mynamespace"))
		Expect(err).ToNot(HaveOccurred())
		uaa := k.ForSubChart("uaa", "uaa", &semver.Version{}, 0).(*k8s.K8sInMemory)
		_, err = uaa.GetObject("secret", "uaa-secret", nil)
		Expect(err).ToNot(HaveOccurred())
		my := k.ForSubChart("mynamespace", "uaa", &semver.Version{}, 0).(*k8s.K8sInMemory)
		_, err = my.GetObject("statefulset", "uaa-master", nil)
		Expect(err).ToNot(HaveOccurred())
	})

})
