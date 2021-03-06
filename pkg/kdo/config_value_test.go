package kdo

import (
	"bytes"
	"io/ioutil"

	"github.com/k14s/starlark-go/starlark"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sap/kubernetes-deployment-orchestrator/pkg/k8s"
)

var _ = Describe("config value", func() {
	Context("config value", func() {
		Context("string type", func() {
			var cv *jewel

			BeforeEach(func() {
				x, err := makeConfigValue(nil, nil, starlark.Tuple{starlark.String("name")}, []starlark.Tuple{{starlark.String("type"), starlark.String("string")}})
				Expect(err).NotTo(HaveOccurred())
				cv = x.(*jewel)
				configValueStdin = ioutil.NopCloser(bytes.NewBuffer([]byte("test\n")))

			})

			It("behaves like starlark value", func() {
				Expect(cv.String()).To(ContainSubstring("name = name"))
				Expect(func() { cv.Hash() }).Should(Panic())
				Expect(cv.Type()).To(Equal("config_value"))
				Expect(cv.Truth()).To(BeEquivalentTo(true))

				By("attribute name", func() {
					value, err := cv.Attr("name")
					Expect(err).NotTo(HaveOccurred())
					Expect(value).To(Equal(starlark.String("name")))
					Expect(cv.AttrNames()).To(ContainElement("name"))

				})

				By("attribute value", func() {
					value, err := cv.Attr("value")
					Expect(err).NotTo(HaveOccurred())
					Expect(value.(starlark.String).GoString()).To(ContainSubstring("test"))
					Expect(cv.AttrNames()).To(ContainElement("value"))

				})

			})

			It("reads values from k8s", func() {
				k8s := k8s.NewK8sInMemoryEmpty()
				err := cv.read(&vaultK8s{k8s: k8s})
				Expect(err).NotTo(HaveOccurred())

			})
		})
		Context("bool type", func() {
			var cv *jewel

			BeforeEach(func() {
				x, err := makeConfigValue(nil, nil, starlark.Tuple{starlark.String("name")}, []starlark.Tuple{{starlark.String("type"), starlark.String("bool")}})
				Expect(err).NotTo(HaveOccurred())
				cv = x.(*jewel)
				configValueStdin = ioutil.NopCloser(bytes.NewBuffer([]byte("true\n")))

			})

			It("behaves like starlark value", func() {
				value, err := cv.Attr("value")
				Expect(err).NotTo(HaveOccurred())
				Expect(value.(starlark.String).GoString()).To(ContainSubstring("yes"))
			})
		})
		Context("password type", func() {
			var cv *jewel

			BeforeEach(func() {
				x, err := makeConfigValue(nil, nil, starlark.Tuple{starlark.String("name")}, []starlark.Tuple{{starlark.String("type"), starlark.String("password")}})
				Expect(err).NotTo(HaveOccurred())
				cv = x.(*jewel)
				configValueStdin = ioutil.NopCloser(bytes.NewBuffer([]byte("secret\n")))

			})

			It("behaves like starlark value", func() {
				value, err := cv.Attr("value")
				Expect(err).NotTo(HaveOccurred())
				Expect(value.(starlark.String).GoString()).To(ContainSubstring("secret"))
			})
		})
		Context("selection type", func() {
			var cv *jewel

			BeforeEach(func() {
				x, err := makeConfigValue(nil, nil, starlark.Tuple{starlark.String("name")},
					[]starlark.Tuple{
						{starlark.String("type"), starlark.String("selection")},
						{starlark.String("options"), starlark.NewList([]starlark.Value{starlark.String("one"), starlark.String("two")})},
					})
				Expect(err).NotTo(HaveOccurred())
				cv = x.(*jewel)
				configValueStdin = ioutil.NopCloser(bytes.NewBuffer([]byte("one\n")))

			})

			It("behaves like starlark value", func() {
				value, err := cv.Attr("value")
				Expect(err).NotTo(HaveOccurred())
				Expect(value.(starlark.String).GoString()).To(ContainSubstring("one"))
			})
		})
	})

})
