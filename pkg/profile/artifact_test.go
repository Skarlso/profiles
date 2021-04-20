package profile_test

import (
	"context"
	"fmt"
	"reflect"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/profile"
)

const (
	subscriptionName     = "mySub"
	namespace            = "default"
	branch               = "main"
	profileName          = "profileName"
	chartName1           = "chartOneArtifactName"
	chartPath1           = "chart/artifact/path-one"
	chartName2           = "chartTwoArtifactName"
	chartPath2           = "chart/artifact/path-two"
	helmChartName1       = "helmChartArtifactName1"
	helmChartChart1      = "helmChartChartName1"
	helmChartURL1        = "https://org.github.io/chart"
	helmChartVersion1    = "8.8.1"
	kustomizeName1       = "kustomizeOneArtifactName"
	kustomizePath1       = "kustomize/artifact/path-one"
	kustomizeName2       = "kustomizeTwoArtifactName"
	kustomizePath2       = "kustomize/artifact/path-two"
	profileSubKind       = "ProfileSubscription"
	profileSubAPIVersion = "weave.works/v1alpha1"
)

var (
	profileTypeMeta = metav1.TypeMeta{
		Kind:       profileSubKind,
		APIVersion: profileSubAPIVersion,
	}
)

var _ = Describe("Profile", func() {
	var (
		ctx        = context.Background()
		p          *profile.Profile
		fakeClient client.Client
		scheme     *runtime.Scheme
		pSub       profilesv1.ProfileSubscription
		pDef       profilesv1.ProfileDefinition
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		pSub = profilesv1.ProfileSubscription{
			TypeMeta: profileTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name:      subscriptionName,
				Namespace: namespace,
			},
			Spec: profilesv1.ProfileSubscriptionSpec{
				ProfileURL: "https://github.com/org/repo-name",
				Branch:     branch,
				Values: &apiextensionsv1.JSON{
					Raw: []byte(`{"replicaCount": 3,"service":{"port":8081}}`),
				},
				ValuesFrom: []helmv2.ValuesReference{
					{
						Name:     "nginx-values",
						Kind:     "Secret",
						Optional: true,
					},
				},
			},
		}

		pDef = profilesv1.ProfileDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: profileName,
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "Profile",
				APIVersion: "profiles.fluxcd.io/profilesv1",
			},
			Spec: profilesv1.ProfileDefinitionSpec{
				Description: "foo",
				Artifacts: []profilesv1.Artifact{
					{
						Name: chartName1,
						Path: chartPath1,
						Kind: profilesv1.HelmChartKind,
					},
					{
						Name: chartName2,
						Path: chartPath2,
						Kind: profilesv1.HelmChartKind,
					},
					{
						Name: kustomizeName1,
						Path: kustomizePath1,
						Kind: profilesv1.KustomizeKind,
					},
					{
						Name: kustomizeName2,
						Path: kustomizePath2,
						Kind: profilesv1.KustomizeKind,
					},
					{
						Name: helmChartName1,
						Chart: &profilesv1.Chart{
							URL:     helmChartURL1,
							Name:    helmChartChart1,
							Version: helmChartVersion1,
						},
						Kind: profilesv1.HelmChartKind,
					},
				},
			},
		}

		fakeClient = fake.NewClientBuilder().WithScheme(scheme).Build()
	})

	Describe("MakeArtifacts", func() {
		It("creates a slice of runtime.Object", func() {
			Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
			Expect(helmv2.AddToScheme(scheme)).To(Succeed())
			Expect(profilesv1.AddToScheme(scheme)).To(Succeed())

			p = profile.New(pDef, pSub, fakeClient, logr.Discard())
			o, err := p.MakeArtifacts()
			Expect(err).NotTo(HaveOccurred())

			Expect(o).To(HaveLen(7))
			Expect(o[0]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "GitRepository", APIVersion: "source.toolkit.fluxcd.io/v1beta1"}))
			Expect(o[1]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "HelmRelease", APIVersion: "helm.toolkit.fluxcd.io/v2beta1"}))
			Expect(o[2]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "HelmRelease", APIVersion: "helm.toolkit.fluxcd.io/v2beta1"}))
			Expect(o[3]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "Kustomization", APIVersion: "kustomize.toolkit.fluxcd.io/v1beta1"}))
			Expect(o[4]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "Kustomization", APIVersion: "kustomize.toolkit.fluxcd.io/v1beta1"}))
			Expect(o[5]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "HelmRelease", APIVersion: "helm.toolkit.fluxcd.io/v2beta1"}))
			Expect(o[6]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "HelmRepository", APIVersion: "source.toolkit.fluxcd.io/v1beta1"}))
		})
	})

	Describe("MakeOwnerlessArtifacts", func() {
		It("creates a slice of runtime.Object without ownership set up", func() {
			Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
			Expect(helmv2.AddToScheme(scheme)).To(Succeed())
			Expect(profilesv1.AddToScheme(scheme)).To(Succeed())
			Expect(kustomizev1.AddToScheme(scheme)).To(Succeed())

			p = profile.New(pDef, pSub, fakeClient, logr.Discard())
			o, err := p.MakeArtifacts()
			Expect(err).NotTo(HaveOccurred())

			Expect(o).To(HaveLen(7))
			Expect(o[0]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "GitRepository", APIVersion: "source.toolkit.fluxcd.io/v1beta1"}))
			Expect(o[1]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "HelmRelease", APIVersion: "helm.toolkit.fluxcd.io/v2beta1"}))
			Expect(o[2]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "HelmRelease", APIVersion: "helm.toolkit.fluxcd.io/v2beta1"}))
			Expect(o[3]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "Kustomization", APIVersion: "kustomize.toolkit.fluxcd.io/v1beta1"}))
			Expect(o[4]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "Kustomization", APIVersion: "kustomize.toolkit.fluxcd.io/v1beta1"}))
			Expect(o[5]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "HelmRelease", APIVersion: "helm.toolkit.fluxcd.io/v2beta1"}))
			Expect(o[6]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "HelmRepository", APIVersion: "source.toolkit.fluxcd.io/v1beta1"}))

			Expect(o[0].(*sourcev1.GitRepository).OwnerReferences).To(HaveLen(1))
			Expect(o[1].(*helmv2.HelmRelease).OwnerReferences).To(HaveLen(1))
			Expect(o[2].(*helmv2.HelmRelease).OwnerReferences).To(HaveLen(1))
			Expect(o[3].(*kustomizev1.Kustomization).OwnerReferences).To(HaveLen(1))
			Expect(o[4].(*kustomizev1.Kustomization).OwnerReferences).To(HaveLen(1))
			Expect(o[5].(*helmv2.HelmRelease).OwnerReferences).To(HaveLen(1))
			Expect(o[6].(*sourcev1.HelmRepository).OwnerReferences).To(HaveLen(1))
		})
	})

	Describe("ArtifactStatus", func() {
		BeforeEach(func() {
			p = profile.New(pDef, pSub, fakeClient, logr.Discard())
		})

		When("the artifact exists", func() {
			var (
				gitRes                       *sourcev1.GitRepository
				helmRep                      *sourcev1.HelmRepository
				helmRes1, helmRes2           *helmv2.HelmRelease
				kustomizeRes1, kustomizeRes2 *kustomizev1.Kustomization
				condition                    metav1.Condition
			)

			BeforeEach(func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				Expect(profilesv1.AddToScheme(scheme)).To(Succeed())
				Expect(kustomizev1.AddToScheme(scheme)).To(Succeed())

				res, err := p.MakeArtifacts()
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(HaveLen(7))
				gitRes = res[0].(*sourcev1.GitRepository)
				helmRes1 = res[1].(*helmv2.HelmRelease)
				helmRes2 = res[2].(*helmv2.HelmRelease)
				kustomizeRes1 = res[3].(*kustomizev1.Kustomization)
				kustomizeRes2 = res[4].(*kustomizev1.Kustomization)
				helmRep = res[6].(*sourcev1.HelmRepository)
				Expect(fakeClient.Create(ctx, gitRes)).To(Succeed())
				Expect(fakeClient.Create(ctx, helmRes1)).To(Succeed())
				Expect(fakeClient.Create(ctx, helmRes2)).To(Succeed())
				Expect(fakeClient.Create(ctx, kustomizeRes1)).To(Succeed())
				Expect(fakeClient.Create(ctx, kustomizeRes2)).To(Succeed())
				Expect(fakeClient.Create(ctx, helmRep)).To(Succeed())
				condition = metav1.Condition{
					Type:               "Ready",
					Status:             "True",
					Reason:             "foo",
					LastTransitionTime: metav1.Now(),
				}

				conditions := []metav1.Condition{condition}
				gitResNew := gitRes.DeepCopyObject().(*sourcev1.GitRepository)
				gitResNew.Status.Conditions = conditions
				Expect(fakeClient.Status().Patch(ctx, gitResNew, client.MergeFrom(gitRes))).To(Succeed())

				helmRepNew := helmRep.DeepCopyObject().(*sourcev1.HelmRepository)
				helmRepNew.Status.Conditions = conditions
				Expect(fakeClient.Status().Patch(ctx, helmRepNew, client.MergeFrom(helmRep))).To(Succeed())

				helmResNew := helmRes1.DeepCopyObject().(*helmv2.HelmRelease)
				helmResNew.Status.Conditions = conditions
				Expect(fakeClient.Status().Patch(ctx, helmResNew, client.MergeFrom(helmRes1))).To(Succeed())

				helmResNew = helmRes2.DeepCopyObject().(*helmv2.HelmRelease)
				helmResNew.Status.Conditions = conditions
				Expect(fakeClient.Status().Patch(ctx, helmResNew, client.MergeFrom(helmRes2))).To(Succeed())

				helmRes2 = res[5].(*helmv2.HelmRelease)
				Expect(fakeClient.Create(ctx, helmRes2)).To(Succeed())
				helmResNew = helmRes2.DeepCopyObject().(*helmv2.HelmRelease)
				helmResNew.Status.Conditions = conditions
				Expect(fakeClient.Status().Patch(ctx, helmResNew, client.MergeFrom(helmRes2))).To(Succeed())

				kustomizeResNew := kustomizeRes1.DeepCopyObject().(*kustomizev1.Kustomization)
				kustomizeResNew.Status.Conditions = conditions
				Expect(fakeClient.Status().Patch(ctx, kustomizeResNew, client.MergeFrom(kustomizeRes1))).To(Succeed())

				kustomizeResNew = kustomizeRes2.DeepCopyObject().(*kustomizev1.Kustomization)
				kustomizeResNew.Status.Conditions = conditions
				Expect(fakeClient.Status().Patch(ctx, kustomizeResNew, client.MergeFrom(kustomizeRes2))).To(Succeed())
			})

			When("the artifacts are all ready=true", func() {
				It("returns empty condition", func() {
					status, err := p.ArtifactStatus(ctx)
					Expect(err).NotTo(HaveOccurred())
					Expect(status.ResourcesExist).To(BeTrue())
					Expect(status.NotReadyConditions).To(HaveLen(0))
				})
			})

			When("the an artifact is ready=false", func() {
				BeforeEach(func() {
					condition = metav1.Condition{
						Type:               "Ready",
						Status:             "False",
						Reason:             "foo",
						LastTransitionTime: metav1.Now(),
					}
					helmResNew := helmRes1.DeepCopyObject().(*helmv2.HelmRelease)
					helmResNew.Status.Conditions = []metav1.Condition{condition}
					Expect(fakeClient.Status().Patch(ctx, helmResNew, client.MergeFrom(helmRes1))).To(Succeed())
				})

				It("returns the ready=false condition", func() {
					status, err := p.ArtifactStatus(ctx)
					Expect(err).NotTo(HaveOccurred())
					Expect(status.ResourcesExist).To(BeTrue())
					Expect(status.NotReadyConditions).To(HaveLen(1))
					//The GET returned from k8sclient mutatates the time format
					//this hack overrides them so it doesn't cause the equal to fail
					now := metav1.Now()
					status.NotReadyConditions[0].LastTransitionTime = now
					condition.LastTransitionTime = now
					Expect(status.NotReadyConditions[0]).To(Equal(condition))
				})
			})

			When("an artifact is ready=unknown", func() {
				BeforeEach(func() {
					condition = metav1.Condition{
						Type:               "Ready",
						Status:             "Unknown",
						Reason:             "foo",
						LastTransitionTime: metav1.Now(),
					}
					helmResNew := helmRes1.DeepCopyObject().(*helmv2.HelmRelease)
					helmResNew.Status.Conditions = []metav1.Condition{condition}
					Expect(fakeClient.Status().Patch(ctx, helmResNew, client.MergeFrom(helmRes1))).To(Succeed())
				})

				It("returns the ready=unknown condition", func() {
					status, err := p.ArtifactStatus(ctx)
					Expect(err).NotTo(HaveOccurred())
					Expect(status.ResourcesExist).To(BeTrue())
					//The GET returned from k8sclient mutatates the time format
					//this hack overrides them so it doesn't cause the equal to fail
					now := metav1.Now()
					status.NotReadyConditions[0].LastTransitionTime = now
					condition.LastTransitionTime = now
					Expect(status.NotReadyConditions[0]).To(Equal(condition))
				})
			})
		})

		When("the artifact don't exist", func() {
			It("returns false", func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				Expect(profilesv1.AddToScheme(scheme)).To(Succeed())

				status, err := p.ArtifactStatus(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(status.ResourcesExist).To(BeFalse())
			})
		})
	})

	Describe("ReconcileArtifacts", func() {
		assertResources := func() {
			gitRefName := fmt.Sprintf("%s-%s-%s", subscriptionName, "repo-name", branch)
			gitRepo := sourcev1.GitRepository{}
			err := fakeClient.Get(ctx, client.ObjectKey{Name: gitRefName, Namespace: namespace}, &gitRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(gitRepo.Spec.URL).To(Equal("https://github.com/org/repo-name"))
			Expect(gitRepo.Spec.Reference.Branch).To(Equal(branch))
			Expect(gitRepo.OwnerReferences).To(HaveLen(1))
			Expect(gitRepo.OwnerReferences[0].Name).To(Equal(subscriptionName))
			Expect(gitRepo.OwnerReferences[0].Kind).To(Equal(profileSubKind))
			Expect(gitRepo.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
			Expect(*gitRepo.OwnerReferences[0].Controller).To(BeTrue())

			helmReleaseName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, chartName1)
			helmRelease := helmv2.HelmRelease{}
			err = fakeClient.Get(ctx, client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, &helmRelease)
			Expect(err).NotTo(HaveOccurred())
			Expect(helmRelease.Spec.Chart.Spec.Chart).To(Equal(chartPath1))
			Expect(helmRelease.Spec.Chart.Spec.SourceRef).To(Equal(
				helmv2.CrossNamespaceObjectReference{
					Kind:      "GitRepository",
					Name:      gitRefName,
					Namespace: namespace,
				},
			))
			Expect(helmRelease.GetValues()).To(Equal(map[string]interface{}{
				"replicaCount": float64(3),
				"service": map[string]interface{}{
					"port": float64(8081),
				},
			}))
			Expect(helmRelease.Spec.ValuesFrom).To(Equal([]helmv2.ValuesReference{
				{
					Name:     "nginx-values",
					Kind:     "Secret",
					Optional: true,
				},
			}))
			Expect(helmRelease.OwnerReferences).To(HaveLen(1))
			Expect(helmRelease.OwnerReferences[0].Name).To(Equal(subscriptionName))
			Expect(helmRelease.OwnerReferences[0].Kind).To(Equal(profileSubKind))
			Expect(helmRelease.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
			Expect(*helmRelease.OwnerReferences[0].Controller).To(BeTrue())

			helmReleaseName = fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, chartName2)
			helmRelease = helmv2.HelmRelease{}
			err = fakeClient.Get(ctx, client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, &helmRelease)
			Expect(err).NotTo(HaveOccurred())
			Expect(helmRelease.Spec.Chart.Spec.Chart).To(Equal(chartPath2))
			Expect(helmRelease.Spec.Chart.Spec.SourceRef).To(Equal(
				helmv2.CrossNamespaceObjectReference{
					Kind:      "GitRepository",
					Name:      gitRefName,
					Namespace: namespace,
				},
			))
			Expect(helmRelease.GetValues()).To(Equal(map[string]interface{}{
				"replicaCount": float64(3),
				"service": map[string]interface{}{
					"port": float64(8081),
				},
			}))
			Expect(helmRelease.Spec.ValuesFrom).To(Equal([]helmv2.ValuesReference{
				{
					Name:     "nginx-values",
					Kind:     "Secret",
					Optional: true,
				},
			}))
			Expect(helmRelease.OwnerReferences).To(HaveLen(1))
			Expect(helmRelease.OwnerReferences[0].Name).To(Equal(subscriptionName))
			Expect(helmRelease.OwnerReferences[0].Kind).To(Equal(profileSubKind))
			Expect(helmRelease.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
			Expect(*helmRelease.OwnerReferences[0].Controller).To(BeTrue())

			kustomizeName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, kustomizeName1)
			kustomize := kustomizev1.Kustomization{}
			err = fakeClient.Get(ctx, client.ObjectKey{Name: kustomizeName, Namespace: namespace}, &kustomize)
			Expect(err).NotTo(HaveOccurred())
			Expect(kustomize.Spec.Path).To(Equal(kustomizePath1))
			Expect(kustomize.Spec.TargetNamespace).To(Equal(namespace))
			Expect(kustomize.Spec.Prune).To(BeTrue())
			Expect(kustomize.Spec.Interval).To(Equal(metav1.Duration{Duration: time.Minute * 5}))
			Expect(kustomize.Spec.SourceRef).To(Equal(
				kustomizev1.CrossNamespaceSourceReference{
					Kind:      "GitRepository",
					Name:      gitRefName,
					Namespace: namespace,
				},
			))
			Expect(kustomize.OwnerReferences).To(HaveLen(1))
			Expect(kustomize.OwnerReferences[0].Name).To(Equal(subscriptionName))
			Expect(kustomize.OwnerReferences[0].Kind).To(Equal(profileSubKind))
			Expect(kustomize.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
			Expect(*kustomize.OwnerReferences[0].Controller).To(BeTrue())

			kustomizeName = fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, kustomizeName2)
			kustomize = kustomizev1.Kustomization{}
			err = fakeClient.Get(ctx, client.ObjectKey{Name: kustomizeName, Namespace: namespace}, &kustomize)
			Expect(err).NotTo(HaveOccurred())
			Expect(kustomize.Spec.Path).To(Equal(kustomizePath2))
			Expect(kustomize.Spec.TargetNamespace).To(Equal(namespace))
			Expect(kustomize.Spec.Prune).To(BeTrue())
			Expect(kustomize.Spec.Interval).To(Equal(metav1.Duration{Duration: time.Minute * 5}))
			Expect(kustomize.Spec.SourceRef).To(Equal(
				kustomizev1.CrossNamespaceSourceReference{
					Kind:      "GitRepository",
					Name:      gitRefName,
					Namespace: namespace,
				},
			))
			Expect(kustomize.OwnerReferences).To(HaveLen(1))
			Expect(kustomize.OwnerReferences[0].Name).To(Equal(subscriptionName))
			Expect(kustomize.OwnerReferences[0].Kind).To(Equal(profileSubKind))
			Expect(kustomize.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
			Expect(*kustomize.OwnerReferences[0].Controller).To(BeTrue())

			helmRefName := fmt.Sprintf("%s-%s-%s-%s", subscriptionName, "repo-name", branch, helmChartChart1)
			helmRepo := sourcev1.HelmRepository{}
			err = fakeClient.Get(ctx, client.ObjectKey{Name: helmRefName, Namespace: namespace}, &helmRepo)
			Expect(err).NotTo(HaveOccurred())
			Expect(helmRepo.Spec.URL).To(Equal(helmChartURL1))
			Expect(helmRepo.OwnerReferences).To(HaveLen(1))
			Expect(helmRepo.OwnerReferences[0].Name).To(Equal(subscriptionName))
			Expect(helmRepo.OwnerReferences[0].Kind).To(Equal(profileSubKind))
			Expect(helmRepo.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
			Expect(*helmRepo.OwnerReferences[0].Controller).To(BeTrue())

			helmReleaseName = fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, helmChartName1)
			helmRelease = helmv2.HelmRelease{}
			err = fakeClient.Get(ctx, client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, &helmRelease)
			Expect(err).NotTo(HaveOccurred())
			Expect(helmRelease.Spec.Chart.Spec.Chart).To(Equal(helmChartChart1))
			Expect(helmRelease.Spec.Chart.Spec.Version).To(Equal(helmChartVersion1))
			Expect(helmRelease.Spec.Chart.Spec.SourceRef).To(Equal(
				helmv2.CrossNamespaceObjectReference{
					Kind:      "HelmRepository",
					Name:      helmRefName,
					Namespace: namespace,
				},
			))
			Expect(helmRelease.GetValues()).To(Equal(map[string]interface{}{
				"replicaCount": float64(3),
				"service": map[string]interface{}{
					"port": float64(8081),
				},
			}))
			Expect(helmRelease.Spec.ValuesFrom).To(Equal([]helmv2.ValuesReference{
				{
					Name:     "nginx-values",
					Kind:     "Secret",
					Optional: true,
				},
			}))
			Expect(helmRelease.OwnerReferences).To(HaveLen(1))
			Expect(helmRelease.OwnerReferences[0].Name).To(Equal(subscriptionName))
			Expect(helmRelease.OwnerReferences[0].Kind).To(Equal(profileSubKind))
			Expect(helmRelease.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
			Expect(*helmRelease.OwnerReferences[0].Controller).To(BeTrue())

			Expect(*helmRelease.OwnerReferences[0].Controller).To(BeTrue())
		}

		DescribeTable("ReconcileArtifacts",
			func(modifyExistingResources func()) {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				Expect(kustomizev1.AddToScheme(scheme)).To(Succeed())
				Expect(profilesv1.AddToScheme(scheme)).To(Succeed())

				p = profile.New(pDef, pSub, fakeClient, logr.Discard())

				By("creating the resources when none exists")
				err := p.ReconcileArtifacts(ctx)
				Expect(err).NotTo(HaveOccurred())
				assertResources()

				By("updating the resources when they become out of sync or deleted")
				modifyExistingResources()

				err = p.ReconcileArtifacts(ctx)
				Expect(err).NotTo(HaveOccurred())
				assertResources()
			},
			Entry("resource gets deleted", func() {
				gitRefName := fmt.Sprintf("%s-%s-%s", subscriptionName, "repo-name", branch)
				gitRepo := sourcev1.GitRepository{}
				err := fakeClient.Get(ctx, client.ObjectKey{Name: gitRefName, Namespace: namespace}, &gitRepo)
				Expect(err).NotTo(HaveOccurred())
				err = fakeClient.Delete(ctx, &gitRepo)
				Expect(err).NotTo(HaveOccurred())
			}),
			Entry("git resource gets out of sync", func() {
				gitRefName := fmt.Sprintf("%s-%s-%s", subscriptionName, "repo-name", branch)
				gitRepo := sourcev1.GitRepository{}
				err := fakeClient.Get(ctx, client.ObjectKey{Name: gitRefName, Namespace: namespace}, &gitRepo)
				Expect(err).NotTo(HaveOccurred())
				gitRepo.Spec.URL = "example.com"
				Expect(fakeClient.Update(ctx, &gitRepo)).To(Succeed())
			}),
			Entry("kustomize resource gets out of sync", func() {
				kustomizeName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, kustomizeName2)
				kustomize := kustomizev1.Kustomization{}
				err := fakeClient.Get(ctx, client.ObjectKey{Name: kustomizeName, Namespace: namespace}, &kustomize)
				Expect(err).NotTo(HaveOccurred())
				kustomize.Spec.Path = "new/path"
				Expect(fakeClient.Update(ctx, &kustomize)).To(Succeed())
			}),
			Entry("helmrelease resource gets out of sync", func() {
				helmReleaseName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, chartName1)
				helmRelease := helmv2.HelmRelease{}
				err := fakeClient.Get(ctx, client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, &helmRelease)
				Expect(err).NotTo(HaveOccurred())
				helmRelease.Spec.Chart.Spec.Chart = "new chart"
				Expect(fakeClient.Update(ctx, &helmRelease)).To(Succeed())
			}),
			Entry("helmrepo resource gets out of sync", func() {
				helmRefName := fmt.Sprintf("%s-%s-%s-%s", subscriptionName, "repo-name", branch, helmChartChart1)
				helmRepo := sourcev1.HelmRepository{}
				err := fakeClient.Get(ctx, client.ObjectKey{Name: helmRefName, Namespace: namespace}, &helmRepo)
				Expect(err).NotTo(HaveOccurred())
				helmRepo.Spec.URL = "example.com/new"
				Expect(fakeClient.Update(ctx, &helmRepo)).To(Succeed())
			}),
		)

		When("setting the resource owner fails", func() {
			It("errors", func() {
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.ReconcileArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("failed to set resource ownership")))
			})
		})

		When("the GitRepository create fails", func() {
			It("errors", func() {
				// this is a bit of a hack, but by not adding this resource to the scheme
				// we can force the Create call to fail
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				Expect(kustomizev1.AddToScheme(scheme)).To(Succeed())
				Expect(profilesv1.AddToScheme(scheme)).To(Succeed())
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.ReconcileArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("failed to create GitRepository: no kind is registered for the type v1beta1.GitRepository in scheme")))
			})
		})

		When("the HelmRelease create fails", func() {
			It("errors", func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(kustomizev1.AddToScheme(scheme)).To(Succeed())
				Expect(profilesv1.AddToScheme(scheme)).To(Succeed())
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.ReconcileArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("failed to create HelmRelease: no kind is registered for the type v2beta1.HelmRelease in scheme")))
			})
		})

		When("the Kustomization create fails", func() {
			It("errors", func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(profilesv1.AddToScheme(scheme)).To(Succeed())
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.ReconcileArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("failed to create Kustomization: no kind is registered for the type v1beta1.Kustomization in scheme")))
			})
		})

		When("the Kind of artifact is unknown", func() {
			BeforeEach(func() {
				pDef.Spec.Artifacts[0].Kind = "SomeUnknownKind"
			})

			It("errors", func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(profilesv1.AddToScheme(scheme)).To(Succeed())
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.ReconcileArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("artifact kind \"SomeUnknownKind\" not recognized")))
			})
		})

		When("helmRepository has an invalid artifact", func() {
			It("errors", func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				Expect(kustomizev1.AddToScheme(scheme)).To(Succeed())
				Expect(profilesv1.AddToScheme(scheme)).To(Succeed())

				pDef = profilesv1.ProfileDefinition{
					ObjectMeta: metav1.ObjectMeta{
						Name: profileName,
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Profile",
						APIVersion: "profiles.fluxcd.io/profilesv1",
					},
					Spec: profilesv1.ProfileDefinitionSpec{
						Description: "foo",
						Artifacts: []profilesv1.Artifact{
							{
								Name: helmChartName1,
								Chart: &profilesv1.Chart{
									URL:     helmChartURL1,
									Name:    helmChartChart1,
									Version: helmChartVersion1,
								},
								Path: "https://not.empty",
								Kind: profilesv1.HelmChartKind,
							},
						},
					},
				}
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.ReconcileArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("validation failed for artifact helmChartArtifactName1: expected exactly one, got both: chart, path")))
			})
		})
	})
})

func HaveTypeMeta(expected metav1.TypeMeta) types.GomegaMatcher {
	return &typeMetaMatcher{
		expected: expected,
	}
}

type typeMetaMatcher struct {
	expected metav1.TypeMeta
}

func (m *typeMetaMatcher) Match(actual interface{}) (bool, error) {
	ro, ok := actual.(runtime.Object)
	if !ok {
		return false, fmt.Errorf("HaveTypeMeta expects a runtime.Object")
	}

	tm, err := m.typeMetaFromObject(ro)
	if err != nil {
		return false, fmt.Errorf("failed to get the type meta for object %#v: %w", ro, err)
	}

	return reflect.DeepEqual(tm, m.expected), nil
}

func (m *typeMetaMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v\nto have TypeMeta\n\t%#v", actual, m.expected)
}

func (m *typeMetaMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v\nnot to have TypeMeta\n\t%#v", actual, m.expected)
}

func (m *typeMetaMatcher) typeMetaFromObject(ro runtime.Object) (metav1.TypeMeta, error) {
	ta, err := meta.TypeAccessor(ro)
	if err != nil {
		return metav1.TypeMeta{}, fmt.Errorf("failed to get the type meta for object %#v: %w", ro, err)
	}
	return metav1.TypeMeta{APIVersion: ta.GetAPIVersion(), Kind: ta.GetKind()}, nil
}
