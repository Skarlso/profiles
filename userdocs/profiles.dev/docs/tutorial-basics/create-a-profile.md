---
sidebar_position: 2
---

# Write a profile

:::info

This stage of the tutorial assumes you have prepared your environment correctly.

Please refer back to the [previous section](/docs/tutorial-basics/setup) if not.
:::

In this section we are going to write a very simple profile.

_If you do not wish to learn about writing profiles, please skip ahead to the next section._

---------------------

You will need an account with a git provider (eg GitHub or GitLab) and a text editor.

We are going to write a profile which can be used to install a single component: nginx.
Yes, I know it is not the most exciting thing, but it serves us well for a quick illustration.

Create a new repository:
```bash
mkdir demo-profile
cd demo-profile
git init
# other git setup things, etc
```

At a bare minimum a profile only needs to have one thing in it: a `profile.yaml` file.
In here we define object Kind and the profile components, known as `artifacts`, under the spec.

```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileDefinition
metadata:
  name: nginx
spec:
  artifacts:
    - name: bitnami-nginx
      chart:
        url: https://charts.bitnami.com/bitnami
        name: nginx
        version: "8.9.1"
```

The section of note here is `spec.artifacts`. Each artifact represents a component of a Profile.
Artifacts can be one of:
- Helm Chart
- Raw Kubernetes yaml
- Kustomise patch
- Another profile

In this example we are creating a Profile with just one artifact: a remote Helm Chart.
The artifact type is denoted here by the `chart` key, and we provide further details for the
specific chart with the the `url`, `version` and `name` fields.

Write the above snippet to `profile.yaml` inside your profile repository. Commit the file
and push to the remote.

You have now created a very basic profile which can be consumed by those who have access.
In the next section of this tutorial we will be installing that profile on a cluster.