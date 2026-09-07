# 🔑 Hasura JWT

## 🗒️ Description

This app allows you to create JWTs for hasura.
And have a minimal signup process via email.
It is small (Image size ~10MB) tool written in golang and minimal dependencies.

## 💡 Features

- 🧑‍🤝‍🧑 Users are stored in Postgres and accessed via GraphQL
- ✨ Integrates with GraphQL and Hasura Permissions
- 🔑 JWT tokens.
- ✉️ Emails sent via SMTP.
- 👨‍💻 Written 100% in Golang.
- 📦 Easy to deploy with Docker or Helm.

## ⚙️ Usage

There are several ways to deploy this project.
There is a ready-made container image on GitHub Packages. 📦

You can use it in your environment.

### Kubernetes

There is also a Helmcahart for Kubernetes available.

For example:
```yaml
apiVersion: helm.cattle.io/v1
kind: HelmChart
spec:
  chart: oci://ghcr.io/53845714nf/charts/hasura-jwt
  version: 0.1.2
  valuesContent: |-
    env:
      APP_URL: "https://your-jwt-app-url"
      HASURA_URL: "https://your-hasura-url"
      HASURA_SECRET:
        valueFrom:
          secretKeyRef:
            name: hasura-secrets
            key: admin-secret
      JWT_KEY:
        valueFrom:
          secretKeyRef:
            name: hasura-secrets
            key: jwt-secret
      EMAIL_VERIFICATION: false
    ingress:
      enabled: false
```

### Docker Compose

There is a Docker Compose File for developers, here the Hasura must be adapted.


### Golang

And last but not least, since it is written in golang,
you can export the project to almost all platforms.
(If a platform is explicitly desired, I can create a Github action for it,
 let me know in an issue)

### Environment Variables for Hasura JWT

- `HASURA_URL` - Must be set to the URL of your Hasura instance
  (e.g. `http://localhost:8080/v1/graphql` is also the default value for development).
- `HASURA_SECRET` - Must be set to the admin secret of your Hasura instance.
- `JWT_KEY` - Must be set to a secret key for signing JWTs.
- `EMAIL_VERIFICATION` - Must be set to `false` if you want to disable email verification.
  Default is `true`. It requires the following SMTP settings.
- `APP_URL` - The URL of the app.
  It is used for creating the email verification link.
  Must be reachable from outside. If you are using a reverse proxy,
  it should be the URL of the reverse proxy.
  In otherwise it must have `:3000` in the end.
- `SMTP_HOST` - The SMTP host to use for sending emails.
- `SMTP_PORT` - The SMTP port to use for sending emails. The default is `587`.
- `SMTP_USER` - The username to use for authenticating with the SMTP server.
  It is used as from email address.
- `SMTP_PASSWORD` - The password to use for authenticating with the SMTP server.

### Environment Variables for Hasura

- `JWT_URL` - Must be set to the URL of your Hasura JWT instance (e.g. `http://localhost:3000`).
- `HASURA_GRAPHQL_UNAUTHORIZED_ROLE` - Set to `anonymous` to get access to the
  public schema without a token, and also for login and signup mutations.

### Volume

You can map a volume with the certificates to `/etc/ssl/certs/` in the container.
This helps by problems with the SMTP Authentication. Certificates from the
Alpine package `ca-certificates` are supported by default.

## 📃 Docs

Please take a look at the GitHub [Wiki](https://github.com/53845714nF/hasura-jwt/wiki)
tab there are sequence diagrams for the process (sign up, login) and a database model.

## 🤖 Similar Projects

There is are similar project like this:

- [Hasura Auth](https://github.com/nhost/hasura-auth/tree/main) - It is now deprecated.
- [Backend-Quickstart](https://github.com/ryaino/Backend-Quickstart) -
It's written in Java, but the last commit was 4 years ago.
- [Hasura Docs](https://hasura.io/docs/latest/actions/codegen/python-flask/) -
  It's a Blog post from official Hasura Documentation,
  there is described how to create JWT with Python and Flask.
