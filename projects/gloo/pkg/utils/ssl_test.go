package utils_test

import (
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/config/grpc_credential/v2alpha"
	"github.com/gogo/protobuf/types"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

var _ = Describe("Ssl", func() {

	var (
		upstreamCfg      *v1.UpstreamSslConfig
		downstreamCfg    *v1.SslConfig
		tlsSecret        *v1.TlsSecret
		secret           *v1.Secret
		secrets          v1.SecretList
		configTranslator *SslConfigTranslator
	)

	Context("files", func() {
		BeforeEach(func() {
			upstreamCfg = &v1.UpstreamSslConfig{
				Sni: "test.com",
				SslSecrets: &v1.UpstreamSslConfig_SslFiles{
					SslFiles: &v1.SSLFiles{
						RootCa:  "rootca",
						TlsCert: "tlscert",
						TlsKey:  "tlskey",
					},
				},
			}
			downstreamCfg = &v1.SslConfig{
				SniDomains: []string{"test.com", "test1.com"},
				SslSecrets: &v1.SslConfig_SslFiles{
					SslFiles: &v1.SSLFiles{
						RootCa:  "rootca",
						TlsCert: "tlscert",
						TlsKey:  "tlskey",
					},
				},
			}
			configTranslator = NewSslConfigTranslator(nil)
		})

		DescribeTable("should resolve from files",
			func(c func() CertSource) {
				ValidateCommonContextFiles(configTranslator.ResolveCommonSslConfig(c()))
			},
			Entry("upstreamCfg", func() CertSource { return upstreamCfg }),
			Entry("downstreamCfg", func() CertSource { return downstreamCfg }),
		)

		Context("san", func() {
			It("should error with san and not rootca", func() {
				upstreamCfg.SslSecrets.(*v1.UpstreamSslConfig_SslFiles).SslFiles.RootCa = ""
				upstreamCfg.VerifySubjectAltName = []string{"test"}
				_, err := configTranslator.ResolveCommonSslConfig(upstreamCfg)
				Expect(err).To(HaveOccurred())
			})

			It("should add SAN verification when provided", func() {
				upstreamCfg.VerifySubjectAltName = []string{"test"}
				c, err := configTranslator.ResolveCommonSslConfig(upstreamCfg)
				Expect(err).NotTo(HaveOccurred())
				vctx := c.ValidationContextType.(*envoyauth.CommonTlsContext_ValidationContext).ValidationContext
				Expect(vctx.VerifySubjectAltName).To(Equal(upstreamCfg.VerifySubjectAltName))
			})
		})
	})
	Context("secret", func() {
		BeforeEach(func() {
			tlsSecret = &v1.TlsSecret{
				CertChain:  "tlscert",
				PrivateKey: "tlskey",
				RootCa:     "rootca",
			}
			secret = &v1.Secret{
				Kind: &v1.Secret_Tls{
					Tls: tlsSecret,
				},
				Metadata: core.Metadata{
					Name:      "secret",
					Namespace: "secret",
				},
			}
			ref := secret.Metadata.Ref()
			secrets = v1.SecretList{secret}
			upstreamCfg = &v1.UpstreamSslConfig{
				Sni: "test.com",
				SslSecrets: &v1.UpstreamSslConfig_SecretRef{
					SecretRef: &ref,
				},
			}
			downstreamCfg = &v1.SslConfig{
				SniDomains: []string{"test.com", "test1.com"},
				SslSecrets: &v1.SslConfig_SecretRef{
					SecretRef: &ref,
				},
			}
			configTranslator = NewSslConfigTranslator(secrets)
		})

		It("should error with no secret", func() {
			configTranslator = NewSslConfigTranslator(nil)
			_, err := configTranslator.ResolveCommonSslConfig(upstreamCfg)
			Expect(err).To(HaveOccurred())
		})

		It("should error with wrong secret", func() {
			secret.Kind = &v1.Secret_Aws{}
			_, err := configTranslator.ResolveCommonSslConfig(upstreamCfg)
			Expect(err).To(HaveOccurred())
		})

		DescribeTable("should resolve from secret refs",
			func(c func() CertSource) {
				ValidateCommonContextInline(configTranslator.ResolveCommonSslConfig(c()))
			},
			Entry("upstreamCfg", func() CertSource { return upstreamCfg }),
			Entry("downstreamCfg", func() CertSource { return downstreamCfg }),
		)
		DescribeTable("should fail if only cert is not provided",
			func(c func() CertSource) {
				tlsSecret.CertChain = ""
				_, err := configTranslator.ResolveCommonSslConfig(c())
				Expect(err).To(HaveOccurred())

			},
			Entry("upstreamCfg", func() CertSource { return upstreamCfg }),
			Entry("downstreamCfg", func() CertSource { return downstreamCfg }),
		)
		DescribeTable("should fail if only private key is not provided",
			func(c func() CertSource) {
				tlsSecret.PrivateKey = ""
				_, err := configTranslator.ResolveCommonSslConfig(c())
				Expect(err).To(HaveOccurred())

			},
			Entry("upstreamCfg", func() CertSource { return upstreamCfg }),
			Entry("downstreamCfg", func() CertSource { return downstreamCfg }),
		)
		DescribeTable("should not have validation context if no rootca",
			func(c func() CertSource) {
				tlsSecret.RootCa = ""
				cfg, err := configTranslator.ResolveCommonSslConfig(c())
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.ValidationContextType).To(BeNil())

			},
			Entry("upstreamCfg", func() CertSource { return upstreamCfg }),
			Entry("downstreamCfg", func() CertSource { return downstreamCfg }),
		)

		It("should set require client cert for downstream config", func() {
			cfg, err := configTranslator.ResolveDownstreamSslConfig(downstreamCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.RequireClientCertificate.GetValue()).To(BeTrue())
		})

		It("should set alpn for downstream config", func() {
			cfg, err := configTranslator.ResolveDownstreamSslConfig(downstreamCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.CommonTlsContext.AlpnProtocols).To(Equal([]string{"h2", "http/1.1"}))
		})

		It("should NOT set alpn for upstream config", func() {
			cfg, err := configTranslator.ResolveUpstreamSslConfig(upstreamCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.CommonTlsContext.AlpnProtocols).To(BeEmpty())
		})

		It("should not set require client cert for downstream config with no rootca", func() {
			tlsSecret.RootCa = ""
			cfg, err := configTranslator.ResolveDownstreamSslConfig(downstreamCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.RequireClientCertificate.GetValue()).To(BeFalse())
		})

		It("should set sni for upstream config", func() {
			cfg, err := configTranslator.ResolveUpstreamSslConfig(upstreamCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Sni).To(Equal("test.com"))
		})

		Context("san", func() {
			It("should error with san and not rootca", func() {
				tlsSecret.RootCa = ""
				upstreamCfg.VerifySubjectAltName = []string{"test"}
				_, err := configTranslator.ResolveCommonSslConfig(upstreamCfg)
				Expect(err).To(HaveOccurred())
			})

			It("should add SAN verification when provided", func() {
				upstreamCfg.VerifySubjectAltName = []string{"test"}
				c, err := configTranslator.ResolveCommonSslConfig(upstreamCfg)
				Expect(err).NotTo(HaveOccurred())
				vctx := c.ValidationContextType.(*envoyauth.CommonTlsContext_ValidationContext).ValidationContext
				Expect(vctx.VerifySubjectAltName).To(Equal(upstreamCfg.VerifySubjectAltName))
			})
		})
	})

	Context("sds", func() {
		var (
			sdsConfig *v1.SDSConfig
		)
		BeforeEach(func() {
			sdsConfig = &v1.SDSConfig{
				TargetUri:              "TargetUri",
				CertificatesSecretName: "CertificatesSecretName",
				ValidationContextName:  "ValidationContextName",
				CallCredentials: &v1.CallCredentials{
					FileCredentialSource: &v1.CallCredentials_FileCredentialSource{
						TokenFileName: "TokenFileName",
						Header:        "Header",
					},
				},
			}
			upstreamCfg = &v1.UpstreamSslConfig{
				Sni: "test.com",
				SslSecrets: &v1.UpstreamSslConfig_Sds{
					Sds: sdsConfig,
				},
			}
			configTranslator = NewSslConfigTranslator(nil)
		})

		It("should have a sds setup", func() {
			c, err := configTranslator.ResolveCommonSslConfig(upstreamCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(c.TlsCertificateSdsSecretConfigs).To(HaveLen(1))
			Expect(c.ValidationContextType).ToNot(BeNil())

			vctx := c.ValidationContextType.(*envoyauth.CommonTlsContext_ValidationContextSdsSecretConfig).ValidationContextSdsSecretConfig
			cert := c.TlsCertificateSdsSecretConfigs[0]
			Expect(vctx.Name).To(Equal("ValidationContextName"))
			Expect(cert.Name).To(Equal("CertificatesSecretName"))
			// If they are no equivalent, it means that any serialization is different.
			// see here: https://github.com/envoyproxy/go-control-plane/pull/158
			// and here: https://github.com/envoyproxy/envoy/pull/6241
			// this may lead to envoy updates being too frequent
			Expect(vctx.SdsConfig).To(BeEquivalentTo(cert.SdsConfig))

			getGrpcConfig := func(s *envoyauth.SdsSecretConfig) *envoycore.GrpcService_GoogleGrpc {
				return s.SdsConfig.ConfigSourceSpecifier.(*envoycore.ConfigSource_ApiConfigSource).ApiConfigSource.GrpcServices[0].TargetSpecifier.(*envoycore.GrpcService_GoogleGrpc_).GoogleGrpc
			}

			Expect(getGrpcConfig(vctx).ChannelCredentials).To(BeEquivalentTo(&envoycore.GrpcService_GoogleGrpc_ChannelCredentials{
				CredentialSpecifier: &envoycore.GrpcService_GoogleGrpc_ChannelCredentials_LocalCredentials{
					LocalCredentials: &envoycore.GrpcService_GoogleGrpc_GoogleLocalCredentials{},
				},
			}))
			Expect(getGrpcConfig(vctx).CredentialsFactoryName).To(Equal(MetadataPluginName))

			credPlugin := getGrpcConfig(vctx).CallCredentials[0].CredentialSpecifier.(*envoycore.GrpcService_GoogleGrpc_CallCredentials_FromPlugin).FromPlugin
			Expect(credPlugin.Name).To(Equal(MetadataPluginName))
			var credConfig v2alpha.FileBasedMetadataConfig
			types.UnmarshalAny(credPlugin.GetTypedConfig(), &credConfig)

			Expect(credConfig).To(BeEquivalentTo(v2alpha.FileBasedMetadataConfig{
				SecretData: &envoycore.DataSource{
					Specifier: &envoycore.DataSource_Filename{
						Filename: "TokenFileName",
					},
				},
				HeaderKey: "Header",
			}))

		})

		Context("san", func() {
			It("should error with san and not rootca", func() {
				sdsConfig.ValidationContextName = ""
				upstreamCfg.VerifySubjectAltName = []string{"test"}
				_, err := configTranslator.ResolveCommonSslConfig(upstreamCfg)
				Expect(err).To(HaveOccurred())
			})

			It("should add SAN verification when provided", func() {
				upstreamCfg.VerifySubjectAltName = []string{"test"}
				c, err := configTranslator.ResolveCommonSslConfig(upstreamCfg)
				Expect(err).NotTo(HaveOccurred())
				vctx := c.ValidationContextType.(*envoyauth.CommonTlsContext_CombinedValidationContext).CombinedValidationContext
				Expect(vctx.DefaultValidationContext.VerifySubjectAltName).To(Equal(upstreamCfg.VerifySubjectAltName))
			})
		})
	})

})

func ValidateCommonContextFiles(tlsCfg *envoyauth.CommonTlsContext, err error) {

	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	validationCtx := tlsCfg.GetValidationContext()
	ExpectWithOffset(1, validationCtx).ToNot(BeNil())
	ExpectWithOffset(1, validationCtx.TrustedCa.GetFilename()).To(Equal("rootca"))

	ExpectWithOffset(1, tlsCfg.GetTlsCertificates()[0].GetCertificateChain().GetFilename()).To(Equal("tlscert"))
	ExpectWithOffset(1, tlsCfg.GetTlsCertificates()[0].GetPrivateKey().GetFilename()).To(Equal("tlskey"))

}

func ValidateCommonContextInline(tlsCfg *envoyauth.CommonTlsContext, err error) {

	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	validationCtx := tlsCfg.GetValidationContext()
	ExpectWithOffset(1, validationCtx).ToNot(BeNil())
	ExpectWithOffset(1, validationCtx.TrustedCa.GetInlineString()).To(Equal("rootca"))

	ExpectWithOffset(1, tlsCfg.GetTlsCertificates()[0].GetCertificateChain().GetInlineString()).To(Equal("tlscert"))
	ExpectWithOffset(1, tlsCfg.GetTlsCertificates()[0].GetPrivateKey().GetInlineString()).To(Equal("tlskey"))

}
