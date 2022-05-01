# Security Measures

## Two Factor Authentication
All external services we use, such as Contabo for server hosting, Github and AWS, require two-factor authentication to log in.

## DDoS Protection
Contabo offers a [DDoS protection](https://contabo.com/en/ddos-protection/) for any hosted server.

## Rate Limiting
Security critical routes in our API such as Login, Signup and automatic card generation are limited to 10 requests per minute. Since the state of the rate limiter is saved in memory, it is possible for an attacker to make double the amount of requests in case the state is somehow reset due to a server crash or something similar. The rate limiter also allows an attacker to hit a window where it can make double the amount of requests per minute whenever the rate limiter is reset.

## Fail2Ban
Fail2Ban is used on our VPN server to prevent brute force attacks and ban IPs on unsuccessful ssh login attempts.

## UFW
UFW is used as a firewall to only allow certain connections such as HTTP and ssh. Since we use Docker as our container engine, it is essential to note that Docker overrides the underlying firewall settings and ignores UFW. This has been an ongoing ["issue"](https://www.google.com/search?q=docker+ufw+flaw+github&rlz=1C1ONGR_deDE974DE974&sxsrf=ALiCzsZUw17Eiyo1ybzlTGm6O6gw6oc2Rw%3A1651415445775&ei=lZluYpXxLuqGxc8PzoqEiAI&ved=0ahUKEwjVj5b4wb73AhVqQ_EDHU4FASEQ4dUDCA8&uact=5&oq=docker+ufw+flaw+github&gs_lcp=Cgdnd3Mtd2l6EAM6BwgjELADECc6BwgAEEcQsAM6BAgjECc6BggAEBYQHjoHCCEQChCgAToICCEQFhAdEB46BQghEKABSgQIQRgASgQIRhgAUL4FWNQMYKwNaAFwAXgAgAG2AYgBzQaSAQM4LjGYAQCgAQHIAQnAAQE&sclient=gws-wiz). Since we only want to allow incoming HTTP connections, we can prevent that by making sure to only bind localhost to our docker containers.

## SSL
All incoming HTTP traffic is handled by [traefik](https://traefik.io/), an open-source reverse proxy that is used to proxy all incoming HTTP traffic and automatically generates SSL certificates for all domains. Internal services do not communicate over HTTPS but instead, use HTTP.

## Secrets
Secrets are currently passed via a .env file stored on the server. This should be improved by, for example, using a dedicated secret manager such as [vault](https://www.vaultproject.io/) since an attacker could very easily read those files.

## Dependencies
Since dependencies always pose a certain security risk, we try to minimize using any.

## User Accounts/Authentication
In order to authenticate users, we use JWT tokens to avoid server-side state. HS256 is used for signing the token. The expiration time of a token can be controlled via our config and is currently set to 7 days.
If the token is compromised, there is nothing we can do currently except for disabling the user account.
However, we want to circumvent this in the future by using a blacklist for tokens.

In addition to that, it should also be noted that we do not have any way for users to reset their passwords or validate their email addresses. The latter definitely needs to be fixed since someone can just signup with another person's email address.

## User authorization
Right now, there are no user authorization mechanisms except for disallowing access to specific routes based on whether a user is part of the beta or not. This is handled as a field inside the user collection and saved as a claim in the JWT token and checked on the API level.

## Password Hashing
We use bcrypt to hash passwords which also automatically adds a salt.

## XSS
React itself partially prevents XSS and only works in certain edge cases such as dynamic links or dangerouslySetInnerHTML. Since we do not use the latter and there is no shared content between users, we can safely assume that no XSS is possible.

## CSRF
It is theoretically possible to perform a CSRF attack since the JWT token is stored as an HTTP-only cookie. However, we want to circumvent that by either implementing a mechanism to prevent CSRF or storing the JWT token in local storage.

## CORS
Since our API runs on a different domain, we implement CORS in the backend to allow these requests to be made. This is not a security measure but rather a way to allow the frontend to communicate with the backend.

## Least Privilege
We try to keep the least privilege where possible. This includes developers working on the project and system, such as our CI system connected to AWS.

## SSH authentication
The only way to connect to our VPS is via SSH. That only works via public key authentication.

## VPS
There are still a few things that we need to implement on our VPS for example disabling root login via ssh.