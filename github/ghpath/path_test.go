package ghpath_test

import (
	"fmt"
	"testing/quick"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unmango/aferox/github/ghpath"
)

var _ = Describe("Path", func() {
	Describe("Parse", func() {
		DescribeTable("Owner URL",
			Entry(nil, "https://github.com/unmango", "unmango"),
			Entry(nil, "github.com/unmango", "unmango"),
			Entry(nil, "https://api.github.com/unmango", "unmango"),
			Entry(nil, "api.github.com/unmango", "unmango"),
			func(input, name string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Owner()).To(Equal(name))
			},
		)

		DescribeTable("Owner",
			Entry(nil, "unmango", "unmango"),
			func(input, name string) {
				res, err := ghpath.Parse(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Owner()).To(Equal(name))
			},
		)

		DescribeTable("Repository URL",
			Entry(nil, "https://github.com/unmango/go", "go"),
			Entry(nil, "github.com/unmango/go", "go"),
			Entry(nil, "https://api.github.com/unmango/go", "go"),
			Entry(nil, "api.github.com/unmango/go", "go"),
			func(input, name string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Repository()).To(Equal(name))
			},
		)

		DescribeTable("Repository",
			Entry(nil, "unmango/go", "go"),
			func(input, name string) {
				res, err := ghpath.Parse(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Repository()).To(Equal(name))
			},
		)

		DescribeTable("No Repository URL",
			Entry(nil, "https://github.com/unmango"),
			Entry(nil, "github.com/unmango"),
			Entry(nil, "https://api.github.com/unmango"),
			Entry(nil, "api.github.com/unmango"),
			func(input string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				_, err = res.Repository()
				Expect(err).To(MatchError("no repository"))
			},
		)

		DescribeTable("No Repository",
			Entry(nil, "unmango"),
			func(input string) {
				res, err := ghpath.Parse(input)

				Expect(err).NotTo(HaveOccurred())
				_, err = res.Repository()
				Expect(err).To(MatchError("no repository"))
			},
		)

		DescribeTable("Branch URL",
			Entry(nil, "https://github.com/unmango/go/tree/main", "main"),
			Entry(nil, "https://api.github.com/unmango/go/tree/main", "main"),
			Entry(nil, "api.github.com/unmango/go/tree/main", "main"),
			Entry(nil, "github.com/unmango/go/tree/main", "main"),
			Entry(nil, "https://raw.githubusercontent.com/unmango/go/refs/heads/main/fs/fold.go", "main"),
			func(input, name string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Branch()).To(Equal(name))
			},
		)

		DescribeTable("Branch",
			Entry(nil, "unmango/go/tree/main", "main"),
			Entry(nil, "unmango/go/tree/main/fs", "main"),
			Entry(nil, "unmango/go/tree/main/fs/path_test.go", "main"),
			Entry(nil, "unmango/go/tree/feature-name", "feature-name"),
			func(input, name string) {
				res, err := ghpath.Parse(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Branch()).To(Equal(name))
			},
		)

		DescribeTable("Not a Branch URL",
			Entry(nil, "https://github.com/unmango"),
			Entry(nil, "github.com/unmango"),
			Entry(nil, "https://api.github.com/unmango"),
			Entry(nil, "api.github.com/unmango"),
			Entry(nil, "https://github.com/unmango/go"),
			Entry(nil, "github.com/unmango/go"),
			Entry(nil, "https://api.github.com/unmango/go"),
			Entry(nil, "api.github.com/unmango/go"),
			func(input string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				_, err = res.Branch()
				Expect(err).To(MatchError("not a branch"))
			},
		)

		DescribeTable("Not a Branch",
			Entry(nil, "unmango"),
			Entry(nil, "unmango/go"),
			func(input string) {
				res, err := ghpath.Parse(input)

				Expect(err).NotTo(HaveOccurred())
				_, err = res.Branch()
				Expect(err).To(MatchError("not a branch"))
			},
		)

		DescribeTable("No Branch URL",
			Entry(nil, "https://github.com/unmango/go/tree"),
			Entry(nil, "https://api.github.com/unmango/go/tree"),
			Entry(nil, "api.github.com/unmango/go/tree"),
			Entry(nil, "github.com/unmango/go/tree"),
			Entry(nil, "https://raw.githubusercontent.com/unmango/go/refs/heads"),
			Entry(nil, "https://raw.githubusercontent.com/unmango/go/refs"),
			func(input string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				_, err = res.Branch()
				Expect(err).To(MatchError("no branch"))
			},
		)

		DescribeTable("Release URL",
			Entry(nil, "https://github.com/unmango/go/releases/tag/v0.0.69", "v0.0.69"),
			Entry(nil, "https://api.github.com/unmango/go/releases/tag/v0.0.69", "v0.0.69"),
			Entry(nil, "api.github.com/unmango/go/releases/tag/v0.0.69", "v0.0.69"),
			Entry(nil, "github.com/unmango/go/releases/tag/v0.0.69", "v0.0.69"),
			Entry(nil, "https://github.com/unmango/go/releases/download/v0.0.69", "v0.0.69"),
			Entry(nil, "https://api.github.com/unmango/go/releases/download/v0.0.69", "v0.0.69"),
			Entry(nil, "api.github.com/unmango/go/releases/download/v0.0.69", "v0.0.69"),
			Entry(nil, "github.com/unmango/go/releases/download/v0.0.69", "v0.0.69"),
			func(input, name string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Release()).To(Equal(name))
			},
		)

		DescribeTable("Release",
			Entry(nil, "unmango/go/releases/tag/v0.0.69", "v0.0.69"),
			Entry(nil, "unmango/go/releases/download/v0.0.69", "v0.0.69"),
			func(input, name string) {
				res, err := ghpath.Parse(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Release()).To(Equal(name))
			},
		)

		DescribeTable("No Release URL",
			Entry(nil, "https://github.com/unmango/go/releases/bleh/v0.0.69", "v0.0.69"),
			Entry(nil, "https://api.github.com/unmango/go/releases/bleh/v0.0.69", "v0.0.69"),
			Entry(nil, "api.github.com/unmango/go/releases/bleh/v0.0.69", "v0.0.69"),
			Entry(nil, "github.com/unmango/go/releases/bleh/v0.0.69", "v0.0.69"),
			func(input, name string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				_, err = res.Release()
				Expect(err).To(MatchError("no release"))
			},
		)

		DescribeTable("No Release",
			Entry(nil, "unmango/go/releases/bleh/v0.0.69", "v0.0.69"),
			func(input, name string) {
				res, err := ghpath.Parse(input)

				Expect(err).NotTo(HaveOccurred())
				_, err = res.Release()
				Expect(err).To(MatchError("no release"))
			},
		)

		DescribeTable("Asset URL",
			Entry(nil, "https://github.com/unmango/go/releases/tag/v0.0.69/my-asset.tar.gz", "my-asset.tar.gz"),
			Entry(nil, "https://api.github.com/unmango/go/releases/tag/v0.0.69/my-asset.tar.gz", "my-asset.tar.gz"),
			Entry(nil, "api.github.com/unmango/go/releases/tag/v0.0.69/my-asset.tar.gz", "my-asset.tar.gz"),
			Entry(nil, "github.com/unmango/go/releases/tag/v0.0.69/my-asset.tar.gz", "my-asset.tar.gz"),
			Entry(nil, "https://github.com/unmango/go/releases/download/v0.0.69/my-asset.tar.gz", "my-asset.tar.gz"),
			Entry(nil, "https://api.github.com/unmango/go/releases/download/v0.0.69/my-asset.tar.gz", "my-asset.tar.gz"),
			Entry(nil, "api.github.com/unmango/go/releases/download/v0.0.69/my-asset.tar.gz", "my-asset.tar.gz"),
			Entry(nil, "github.com/unmango/go/releases/download/v0.0.69/my-asset.tar.gz", "my-asset.tar.gz"),
			func(input, name string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Asset()).To(Equal(name))
			},
		)

		DescribeTable("Asset",
			Entry(nil, "unmango/go/releases/tag/v0.0.69/my-asset.tar.gz", "my-asset.tar.gz"),
			Entry(nil, "unmango/go/releases/download/v0.0.69/my-asset.tar.gz", "my-asset.tar.gz"),
			func(input, name string) {
				res, err := ghpath.Parse(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Asset()).To(Equal(name))
			},
		)

		DescribeTable("No Asset URL",
			Entry(nil, "https://github.com/unmango/go/releases/tag/v0.0.69"),
			Entry(nil, "https://api.github.com/unmango/go/releases/tag/v0.0.69"),
			Entry(nil, "api.github.com/unmango/go/releases/tag/v0.0.69"),
			Entry(nil, "github.com/unmango/go/releases/tag/v0.0.69"),
			Entry(nil, "https://github.com/unmango/go/releases/download/v0.0.69"),
			Entry(nil, "https://api.github.com/unmango/go/releases/download/v0.0.69"),
			Entry(nil, "api.github.com/unmango/go/releases/download/v0.0.69"),
			Entry(nil, "github.com/unmango/go/releases/download/v0.0.69"),
			func(input string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				_, err = res.Asset()
				Expect(err).To(MatchError("no asset"))
			},
		)

		DescribeTable("No Asset",
			Entry(nil, "unmango/go/releases/tag/v0.0.69"),
			Entry(nil, "unmango/go/releases/download/v0.0.69"),
			func(input string) {
				res, err := ghpath.Parse(input)

				Expect(err).NotTo(HaveOccurred())
				_, err = res.Asset()
				Expect(err).To(MatchError("no asset"))
			},
		)

		DescribeTable("Content URL",
			Entry(nil, "https://github.com/unmango/go/tree/main", []string{}),
			Entry(nil, "https://api.github.com/unmango/go/tree/main", []string{}),
			Entry(nil, "api.github.com/unmango/go/tree/main", []string{}),
			Entry(nil, "github.com/unmango/go/tree/main", []string{}),
			Entry(nil, "https://raw.githubusercontent.com/unmango/go/refs/heads/main/fs/fold.go", []string{"fs", "fold.go"}),
			func(input string, parts []string) {
				res, err := ghpath.ParseUrl(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Content()).To(Equal(parts))
			},
		)

		DescribeTable("Content",
			Entry(nil, "unmango/go/tree/main", []string{}),
			Entry(nil, "unmango/go/tree/main/fs", []string{"fs"}),
			Entry(nil, "unmango/go/tree/main/fs/path_test.go", []string{"fs", "path_test.go"}),
			func(input string, parts []string) {
				res, err := ghpath.Parse(input)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.Content()).To(Equal(parts))
			},
		)

		DescribeTable("Parts",
			Entry(nil,
				[]string{"unmango", "go"},
				"unmango/go",
			),
			Entry(nil,
				[]string{"unmango", "go", "refs", "heads", "main"},
				"unmango/go/refs/heads/main",
			),
			Entry(nil,
				[]string{"unmango", "go", "releases", "tag", "v0.0.69"},
				"unmango/go/releases/tag/v0.0.69",
			),
			func(parts []string, expected string) {
				path, err := ghpath.Parse(parts...)

				Expect(err).NotTo(HaveOccurred())
				Expect(path.String()).To(Equal(expected))
			},
		)
	})

	Describe("OwnerPath", func() {
		It("should Parse repo", func() {
			p := ghpath.NewOwnerPath("testing")

			r, err := p.Parse("repo-name")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Repository()).To(Equal("repo-name"))
		})

		It("should assume release when parsing len 2", func() {
			p := ghpath.NewOwnerPath("testing")

			r, err := p.Parse("repo/release-name")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Repository()).To(Equal("repo"))
			Expect(r.Release()).To(Equal("release-name"))
		})

		It("should Parse release", func() {
			p := ghpath.NewOwnerPath("testing")

			r, err := p.Parse("repo/releases/tag/release-name")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Repository()).To(Equal("repo"))
			Expect(r.Release()).To(Equal("release-name"))
		})

		It("should assume asset when parsing len 5", func() {
			p := ghpath.NewOwnerPath("testing")

			r, err := p.Parse("repo/releases/tag/release-name/asset.tar.gz")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Repository()).To(Equal("repo"))
			Expect(r.Release()).To(Equal("release-name"))
			Expect(r.Asset()).To(Equal("asset.tar.gz"))
		})

		It("should Parse asset", func() {
			p := ghpath.NewOwnerPath("testing")

			r, err := p.Parse("repo/releases/tag/release-name/download/asset.tar.gz")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Repository()).To(Equal("repo"))
			Expect(r.Release()).To(Equal("release-name"))
			Expect(r.Asset()).To(Equal("asset.tar.gz"))
		})

		It("should stringify", func() {
			fn := func(owner string) bool {
				p := ghpath.NewOwnerPath(owner)

				actual := p.String()

				return actual == fmt.Sprintf("https://github.com/%s", owner)
			}

			Expect(quick.Check(fn, nil)).To(Succeed())
		})
	})

	Describe("RepositoryPath", func() {
		It("should Parse release", func() {
			p := ghpath.NewRepositoryPath("owner", "repo")

			r, err := p.Parse("releases/tag/release-name")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Repository()).To(Equal("repo"))
			Expect(r.Release()).To(Equal("release-name"))
		})

		It("should Parse branch from tree", func() {
			p := ghpath.NewRepositoryPath("owner", "repo")

			r, err := p.Parse("tree/branch-name")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Owner()).To(Equal("owner"))
			Expect(r.Repository()).To(Equal("repo"))
			Expect(r.Branch()).To(Equal("branch-name"))
		})

		It("should Parse branch from ref", func() {
			p := ghpath.NewRepositoryPath("owner", "repo")

			r, err := p.Parse("refs/heads/branch-name")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Owner()).To(Equal("owner"))
			Expect(r.Repository()).To(Equal("repo"))
			Expect(r.Branch()).To(Equal("branch-name"))
		})

		It("should stringify", func() {
			fn := func(owner, repo string) bool {
				p := ghpath.NewRepositoryPath(owner, repo)

				actual := p.String()

				return actual == fmt.Sprintf("https://github.com/%s/%s", owner, repo)
			}

			Expect(quick.Check(fn, nil)).To(Succeed())
		})
	})

	Describe("BranchPath", func() {
		It("should parse", func() {
			p := ghpath.NewBranchPath("owner", "repo", "branch")

			r, err := p.Parse("some-content")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Owner()).To(Equal("owner"))
			Expect(r.Repository()).To(Equal("repo"))
			Expect(r.Branch()).To(Equal("branch"))
			Expect(r.Content()).To(ConsistOf("some-content"))
		})

		It("should stringify", func() {
			fn := func(owner, repo, branch string) bool {
				p := ghpath.NewBranchPath(owner, repo, branch)

				actual := p.String()

				return actual == fmt.Sprintf(
					"https://github.com/%s/%s/tree/%s",
					owner, repo, branch,
				)
			}

			Expect(quick.Check(fn, nil)).To(Succeed())
		})
	})

	Describe("ContentPath", func() {
		It("should parse", func() {
			p := ghpath.NewContentPath("owner", "repo", "branch", "content")

			r, err := p.Parse("other-content")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Owner()).To(Equal("owner"))
			Expect(r.Repository()).To(Equal("repo"))
			Expect(r.Branch()).To(Equal("branch"))
			Expect(r.Content()).To(ConsistOf("content", "other-content"))
		})

		It("should stringify", func() {
			fn := func(owner, repo, branch, content string) bool {
				p := ghpath.NewContentPath(owner, repo, branch, content)

				actual := p.String()

				return actual == fmt.Sprintf(
					"https://github.com/%s/%s/tree/%s/%s",
					owner, repo, branch, content,
				)
			}

			Expect(quick.Check(fn, nil)).To(Succeed())
		})
	})

	Describe("ReleasePath", func() {
		It("should Parse", func() {
			p := ghpath.NewReleasePath("owner", "repo", "release")

			r, err := p.Parse("asset-name")

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Repository()).To(Equal("repo"))
			Expect(r.Release()).To(Equal("release"))
			Expect(r.Asset()).To(Equal("asset-name"))
		})

		It("should stringify", func() {
			fn := func(owner, repo, release string) bool {
				p := ghpath.NewReleasePath(owner, repo, release)

				actual := p.String()

				return actual == fmt.Sprintf(
					"https://github.com/%s/%s/releases/tag/%s",
					owner, repo, release,
				)
			}

			Expect(quick.Check(fn, nil)).To(Succeed())
		})
	})

	DescribeTable("ParseOwner",
		Entry(nil, "UnstoppableMango"),
		Entry(nil, "UnstoppableMango/repo"),
		Entry(nil, "UnstoppableMango/repo/releases/tag/tdl"),
		Entry(nil, "UnstoppableMango/repo/releases/this-is-wrong/thing"),
		Entry(nil, "UnstoppableMango/repo/tree/main"),
		func(input string) {
			p, err := ghpath.Parse(input)
			Expect(err).NotTo(HaveOccurred())

			r, err := ghpath.ParseOwner(p)

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Owner).To(Equal("UnstoppableMango"))
		},
	)

	DescribeTable("ParseRepository",
		Entry(nil, "UnstoppableMango/repo"),
		Entry(nil, "UnstoppableMango/repo/releases/tag/tdl"),
		Entry(nil, "UnstoppableMango/repo/releases/this-is-wrong/thing"),
		Entry(nil, "UnstoppableMango/repo/tree/main"),
		func(input string) {
			p, err := ghpath.Parse(input)
			Expect(err).NotTo(HaveOccurred())

			r, err := ghpath.ParseRepository(p)

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Owner).To(Equal("UnstoppableMango"))
			Expect(r.Repository).To(Equal("repo"))
		},
	)

	DescribeTable("ParseRelease",
		Entry(nil, "UnstoppableMango/repo/releases/tag/tdl"),
		func(input string) {
			p, err := ghpath.Parse(input)
			Expect(err).NotTo(HaveOccurred())

			r, err := ghpath.ParseRelease(p)

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Owner).To(Equal("UnstoppableMango"))
			Expect(r.Repository).To(Equal("repo"))
			Expect(r.Release).To(Equal("tdl"))
		},
	)

	DescribeTable("ParseAsset",
		Entry(nil, "UnstoppableMango/repo/releases/tag/tdl/v0.0.69"),
		func(input string) {
			p, err := ghpath.Parse(input)
			Expect(err).NotTo(HaveOccurred())

			r, err := ghpath.ParseAsset(p)

			Expect(err).NotTo(HaveOccurred())
			Expect(r.Owner).To(Equal("UnstoppableMango"))
			Expect(r.Repository).To(Equal("repo"))
			Expect(r.Release).To(Equal("tdl"))
			Expect(r.Asset).To(Equal("v0.0.69"))
		},
	)
})
